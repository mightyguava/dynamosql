package querybuilder

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/alecthomas/repr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	dberr "github.com/mightyguava/dynamosql/errors"
	"github.com/mightyguava/dynamosql/parser"
	"github.com/mightyguava/dynamosql/schema"
)

type PreparedQuery struct {
	Query       *dynamodb.QueryInput
	FreeParams  FreeParams
	FixedParams map[string]interface{}
}

func PrepareQuery(ctx context.Context, tables *schema.TableLoader, query string) (*PreparedQuery, error) {
	var ast parser.Select
	if err := parser.Parser.ParseString(query, &ast); err != nil {
		return nil, err
	}
	table, err := tables.Get(ctx, ast.From)
	if err != nil {
		return nil, err
	}
	return buildQuery(table, ast)
}

func (pq *PreparedQuery) Build(args []driver.NamedValue) (*dynamodb.QueryInput, error) {
	values := make(map[string]*dynamodb.AttributeValue, len(pq.FreeParams)+len(pq.FixedParams))
	for k, v := range pq.FixedParams {
		av, err := toAttributeValue(v)
		if err != nil {
			return nil, err
		}
		values[k] = av
	}
	if err := bindArgs(pq.FreeParams, args, values); err != nil {
		return nil, err
	}
	req := *pq.Query
	req.ExpressionAttributeValues = values
	return &req, nil
}

func bindArgs(freeParams FreeParams, args []driver.NamedValue, values map[string]*dynamodb.AttributeValue) error {
	freeParams = freeParams.Clone()
	for _, arg := range args {
		if arg.Name == "" {
			return dberr.ErrPositionalArg
		}
		name := ":" + arg.Name
		_, ok := freeParams[name]
		if !ok {
			return fmt.Errorf("binding %q not found", name)
		}
		av, err := toAttributeValue(arg.Value)
		if err != nil {
			return err
		}
		values[name] = av
		delete(freeParams, name)
	}
	if len(freeParams) > 0 {
		for k := range freeParams {
			return fmt.Errorf("missing argument for binding %q", k)
		}
	}
	return nil
}

func toAttributeValue(attr interface{}) (*dynamodb.AttributeValue, error) {
	switch v := attr.(type) {
	case string:
		return &dynamodb.AttributeValue{S: &v}, nil
	case float64:
		return &dynamodb.AttributeValue{N: aws.String(strconv.FormatFloat(v, 'g', -1, 64))}, nil
	case int64:
		return &dynamodb.AttributeValue{N: aws.String(strconv.FormatInt(v, 10))}, nil
	case bool:
		return &dynamodb.AttributeValue{BOOL: &v}, nil
	case nil:
		return &dynamodb.AttributeValue{NULL: aws.Bool(true)}, nil
	default:
		return nil, fmt.Errorf("invalid value type %s", reflect.TypeOf(v))
	}
}

type KeyExpression struct{}
type FilterExpression struct{}

type FreeParams map[string]Empty

func (p FreeParams) Clone() FreeParams {
	copy := make(FreeParams)
	for k, v := range p {
		copy[k] = v
	}
	return copy
}

func buildQuery(table *schema.Table, ast parser.Select) (*PreparedQuery, error) {
	visit := &visitor{Context: NewContext(table)}
	result, err := visit.VisitNodes(ast.Where)
	if err != nil {
		return nil, err
	}
	var projectionExpr *string
	if !ast.Projection.All {
		projectionExpr = aws.String(strings.Join(ast.Projection.Projections, ","))
	}
	req := &dynamodb.QueryInput{
		TableName:              &ast.From,
		ProjectionExpression:   projectionExpr,
		KeyConditionExpression: aws.String(result.Key),
	}
	if result.Filter != "" {
		req.FilterExpression = aws.String(result.Filter)
	}
	return &PreparedQuery{
		Query:       req,
		FreeParams:  visit.Context.FreeParams,
		FixedParams: visit.Context.FixedParams,
	}, nil
}

type Context struct {
	Table       *schema.Table
	FreeParams  map[string]Empty
	FixedParams map[string]interface{}

	genParamCount int
}

func NewContext(table *schema.Table) *Context {
	return &Context{
		Table:       table,
		FreeParams:  make(map[string]Empty),
		FixedParams: make(map[string]interface{}),
	}
}

// NextGeneratedParam returns the next generated parameter placeholder name.
func (c *Context) NextGeneratedParam() string {
	c.genParamCount++
	return fmt.Sprintf(":_gen%d", c.genParamCount)
}

type Empty struct{}

type Acc struct {
	Key    string
	Filter string
}

type visitor struct {
	*Context
}

func (v *visitor) VisitNodes(n interface{}) (Acc, error) {
	switch node := n.(type) {
	case *parser.ConditionExpression:
		var filter, key []string
		for _, expr := range node.Or {
			acc, err := v.VisitNodes(expr)
			if err != nil {
				return Acc{}, err
			}
			if acc.Filter != "" {
				filter = append(filter, acc.Filter)
			}
			if acc.Key != "" {
				return Acc{}, errors.New("primary key cannot appear in OR expression")
			}
		}
		return Acc{
			Filter: strings.Join(key, " OR "),
		}, nil
	case *parser.AndExpression:
		var filter, key []string
		for _, expr := range node.And {
			acc, err := v.VisitNodes(expr)
			if err != nil {
				return Acc{}, err
			}
			if acc.Filter != "" {
				filter = append(filter, acc.Filter)
			}
			if acc.Key != "" {
				key = append(key, acc.Key)
			}
		}
		return Acc{
			Key:    strings.Join(key, " AND "),
			Filter: strings.Join(filter, " AND "),
		}, nil
	case *parser.Condition:
		switch {
		case node.Operand != nil:
			expr := v.VisitTerm(node.Operand)
			if v.Table.IsKey(node.Operand.Operand.Symbol) {
				return Acc{
					Key: expr,
				}, nil
			} else {
				return Acc{
					Filter: expr,
				}, nil
			}
		case node.Function != nil:
			return v.VisitNodes(node.Function)
		case node.Not != nil:
			return v.VisitNodes(node.Not)
		case node.Parenthesized != nil:
			return v.VisitNodes(node.Parenthesized)
		default:
			panic("invalid condition subtype")
		}
	case *parser.ParenthesizedExpression:
		acc, err := v.VisitNodes(node.ConditionExpression)
		if err != nil {
			return Acc{}, err
		}
		return Acc{
			Key:    "(" + acc.Key + ")",
			Filter: "(" + acc.Filter + ")",
		}, nil
	case *parser.NotCondition:
		acc, err := v.VisitNodes(node.Condition)
		if err != nil {
			return Acc{}, err
		}
		if acc.Key != "" {
			return Acc{}, errors.New("primary key cannot appear in NOT expression")
		}
		return Acc{
			Filter: "NOT " + acc.Filter,
		}, nil
	case *parser.FunctionExpression:
		args := node.PathArgument
		more := v.VisitTerm(node.MoreArguments)
		if more != "" {
			args = args + "," + more
		}
		expr := fmt.Sprintf("%s(%s)", node.Function, args)
		if v.Table.IsKey(node.PathArgument) {
			return Acc{Key: expr}, nil
		} else {
			return Acc{Filter: expr}, nil
		}
	default:
		panic("invalid type " + repr.String(node))
	}
}

func (v *visitor) VisitTerm(n interface{}) string {
	switch node := n.(type) {
	case *parser.ConditionOperand:
		return node.Operand.Symbol + " " + v.VisitTerm(node.ConditionRHS)
	case *parser.ConditionRHS:
		switch {
		case node.Compare != nil:
			return v.VisitTerm(node.Compare)
		case node.Between != nil:
			return v.VisitTerm(node.Between)
		case node.In != nil:
			return v.VisitTerm(node.In)
		default:
			panic("invalid rhs")
		}
	case *parser.Compare:
		return node.Operator + " " + v.VisitTerm(node.Operand)
	case *parser.Between:
		return fmt.Sprintf("BETWEEN %s AND %s",
			v.VisitTerm(node.Start), v.VisitTerm(node.End))
	case *parser.In:
		return fmt.Sprintf(" IN (%s)", v.VisitTerm(node.Values))
	case []parser.Value:
		var values []string
		for _, value := range node {
			values = append(values, v.VisitTerm(value))
		}
		return strings.Join(values, ",")
	case *parser.Operand:
		if node.SymbolRef != nil {
			return node.SymbolRef.Symbol
		} else {
			return v.VisitTerm(node.Value)
		}
	case *parser.Value:
		switch {
		case node.PlaceHolder != nil:
			v.Context.FreeParams[*node.PlaceHolder] = Empty{}
			return *node.PlaceHolder
		case node.Number != nil:
			name := v.Context.NextGeneratedParam()
			v.Context.FixedParams[name] = *node.Number
			return name
		case node.String != nil:
			name := v.Context.NextGeneratedParam()
			v.Context.FixedParams[name] = *node.String
			return name
		case node.Boolean != nil:
			name := v.Context.NextGeneratedParam()
			v.Context.FixedParams[name] = bool(*node.Boolean)
			return name
		case node.Null:
			name := v.Context.NextGeneratedParam()
			v.Context.FixedParams[name] = nil
			return name
		default:
			panic("invalid value" + repr.String(node))
		}
	default:
		panic("invalid term: " + repr.String(node))
	}
}

//type visitor func(ctx Context, v interface{})
//
//func visitAndExpression(and *parser.AndExpression) {
//	for _, expr := range and.And {
//		switch {
//		case expr.Operand != nil:
//			visitOperand(ctx, expr.Operand)
//		case expr.Function != nil:
//		case expr.Not != nil:
//		case expr.Parenthesized != nil:
//		}
//	}
//}
//
//func visitOperand(ctx Context, operand *parser.ConditionOperand) {
//	operand.ConditionRHS
//	if ctx.Table.IsKey(operand.Operand.Symbol) {
//		ctx.KeyCondition =
//	}
//}
