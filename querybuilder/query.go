package querybuilder

import (
	"bytes"
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

	"github.com/mightyguava/dynamosql/parser"
	"github.com/mightyguava/dynamosql/schema"
)

var (
	errPositionalArg = errors.New("positional args not supported, use sql.NamedArg to pass named arguments")
)

type PreparedQuery struct {
	Query       *dynamodb.QueryInput
	Columns     []*parser.ProjectionColumn
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
	values, err := bindArgs(pq.FixedParams, pq.FreeParams, args)
	if err != nil {
		return nil, err
	}
	req := *pq.Query
	req.ExpressionAttributeValues = values
	return &req, nil
}

func bindArgs(fixedParams map[string]interface{}, freeParams FreeParams, args []driver.NamedValue) (map[string]*dynamodb.AttributeValue, error) {
	values := make(map[string]*dynamodb.AttributeValue, len(freeParams)+len(fixedParams))

	// Bind fixed params
	for k, v := range fixedParams {
		av, err := toAttributeValue(v)
		if err != nil {
			return nil, err
		}
		values[k] = av
	}

	// Bind free params using the args
	freeParams = freeParams.Clone()
	for _, arg := range args {
		if arg.Name == "" {
			return nil, errPositionalArg
		}
		name := ":" + arg.Name
		_, ok := freeParams[name]
		if !ok {
			return nil, fmt.Errorf("binding %q not found", name)
		}
		av, err := toAttributeValue(arg.Value)
		if err != nil {
			return nil, err
		}
		values[name] = av
		delete(freeParams, name)
	}
	if len(freeParams) > 0 {
		for k := range freeParams {
			return nil, fmt.Errorf("missing argument for binding %q", k)
		}
	}
	return values, nil
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
	ctx := NewContext(table)
	visit := &visitor{Context: ctx}
	kf := extractKeyExpressions(ast.Where, table.IsKey)
	keyExpr, err := buildKeyExpression(ctx, kf.Key)
	if err != nil {
		return nil, err
	}
	filterExpr, err := buildFilterExpression(ctx, kf.Filter)
	if err != nil {
		return nil, err
	}
	var projectionExpr *string
	if !ast.Projection.All {
		expr, err := buildProjectionExpression(ast.Projection)
		if err != nil {
			return nil, err
		}
		projectionExpr = aws.String(expr)
	}
	req := &dynamodb.QueryInput{
		TableName:              &ast.From,
		ProjectionExpression:   projectionExpr,
		KeyConditionExpression: aws.String(keyExpr),
	}
	if filterExpr != "" {
		req.FilterExpression = aws.String(filterExpr)
	}
	return &PreparedQuery{
		Query:       req,
		Columns:     ast.Projection.Columns,
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

type keyAndFilter struct {
	Key    *parser.AndExpression
	Filter *parser.AndExpression
}

func extractKeyExpressions(expr *parser.AndExpression, isKey func(string) bool) *keyAndFilter {
	v := &keyAndFilter{
		Key:    &parser.AndExpression{},
		Filter: &parser.AndExpression{},
	}
	if expr == nil {
		return v
	}
	for _, term := range expr.And {
		if term.Operand != nil && isKey(term.Operand.Operand.Symbol) ||
			term.Function != nil && term.Function.FirstArgIsRef() && isKey(term.Function.Args[0].DocumentPath.String()) {
			v.Key.And = append(v.Key.And, term)
		} else {
			v.Filter.And = append(v.Filter.And, term)
		}
	}
	return v
}

type errHashKey string

func (hashKey errHashKey) Error() string {
	return fmt.Sprintf("partition key must appear exactly once in the WHERE clause, in an equality condition, such as: WHERE %s = :param", string(hashKey))
}

func buildKeyExpression(ctx *Context, key *parser.AndExpression) (string, error) {
	var hashExpr, sortExpr string
	visitor := &visitor{Context: ctx}
	for _, subExpr := range key.And {
		var expr, key string
		if subExpr.Function != nil {
			key = subExpr.Function.Args[0].DocumentPath.String()
			if ctx.Table.HashKey == key {
				return "", errHashKey(ctx.Table.HashKey)
			} else if subExpr.Function.Function != "begins_with" {
				return "", fmt.Errorf("sort key %q may not be used with function %s()", key, subExpr.Function.Function)
			}
			expr = visitor.VisitSimpleExpression(subExpr.Function)
		} else {
			key = subExpr.Operand.Operand.Symbol
			if key == ctx.Table.HashKey {
				if subExpr.Operand.ConditionRHS.Compare == nil || subExpr.Operand.ConditionRHS.Compare.Operator != "=" {
					return "", errHashKey(ctx.Table.HashKey)
				}
			}
			expr = visitor.VisitSimpleExpression(subExpr.Operand)
		}
		if ctx.Table.HashKey == key {
			if hashExpr != "" {
				return "", fmt.Errorf("partition key %q can only appear once in WHERE clause", key)
			}
			hashExpr = expr
		} else if ctx.Table.SortKey == key {
			if sortExpr != "" {
				return "", fmt.Errorf("sort key %q can only appear once in WHERE clause", key)
			}
			sortExpr = expr
		}
	}
	if hashExpr == "" {
		return "", errHashKey(ctx.Table.HashKey)
	}
	if sortExpr != "" {
		return hashExpr + " AND " + sortExpr, nil
	}
	return hashExpr, nil
}

func buildFilterExpression(ctx *Context, filter *parser.AndExpression) (string, error) {
	v := &visitor{Context: ctx}
	return v.VisitFilterExpression(filter)
}

type Empty struct{}

type Acc struct {
	Key    string
	Filter string
}

type visitor struct {
	*Context
}

// VisitFilterExpression visits all nodes in the filter expression tree to build a filter expression.
func (v *visitor) VisitFilterExpression(n interface{}) (string, error) {
	if reflect.ValueOf(n).IsNil() {
		return "", nil
	}
	switch node := n.(type) {
	case *parser.ConditionExpression:
		filter := make([]string, 0, len(node.Or))
		for _, expr := range node.Or {
			acc, err := v.VisitFilterExpression(expr)
			if err != nil {
				return "", err
			}
			if acc != "" {
				filter = append(filter, acc)
			}
		}
		return strings.Join(filter, " OR "), nil
	case *parser.AndExpression:
		filter := make([]string, 0, len(node.And))
		for _, expr := range node.And {
			acc, err := v.VisitFilterExpression(expr)
			if err != nil {
				return "", err
			}
			if acc != "" {
				filter = append(filter, acc)
			}
		}
		return strings.Join(filter, " AND "), nil
	case *parser.Condition:
		switch {
		case node.Operand != nil:
			if v.Context.Table.IsKey(node.Operand.Operand.Symbol) {
				return "", fmt.Errorf("partition key %q may not appear in nested expression", node.Operand.Operand.Symbol)
			}
			return v.VisitSimpleExpression(node.Operand), nil
		case node.Function != nil:
			if node.Function.FirstArgIsRef() && v.Context.Table.IsKey(node.Function.Args[0].DocumentPath.String()) {
				return "", fmt.Errorf("partition key %q may not appear in nested expression", node.Function.Args[0].DocumentPath)
			}
			return v.VisitSimpleExpression(node.Function), nil
		case node.Not != nil:
			return v.VisitFilterExpression(node.Not)
		case node.Parenthesized != nil:
			return v.VisitFilterExpression(node.Parenthesized)
		default:
			panic("invalid condition subtype")
		}
	case *parser.ParenthesizedExpression:
		acc, err := v.VisitFilterExpression(node.ConditionExpression)
		if err != nil {
			return "", err
		}
		return "(" + acc + ")", nil
	case *parser.NotCondition:
		acc, err := v.VisitFilterExpression(node.Condition)
		if err != nil {
			return "", err
		}
		return "NOT " + acc, nil
	default:
		panic("invalid type " + repr.String(node))
	}
}

// VisitSimpleExpression visits unary and binary expressions down to the leaf nodes.
func (v *visitor) VisitSimpleExpression(n interface{}) string {
	switch node := n.(type) {
	case *parser.ConditionOperand:
		return node.Operand.Symbol + " " + v.VisitSimpleExpression(node.ConditionRHS)
	case *parser.FunctionExpression:
		argStr := make([]string, len(node.Args))
		for i, arg := range node.Args {
			argStr[i] = v.VisitSimpleExpression(arg)
		}
		return fmt.Sprintf("%s(%s)", node.Function, strings.Join(argStr, ", "))
	case *parser.FunctionArgument:
		if node.DocumentPath != nil {
			return node.DocumentPath.String()
		}
		return v.VisitSimpleExpression(node.Value)
	case *parser.ConditionRHS:
		switch {
		case node.Compare != nil:
			return v.VisitSimpleExpression(node.Compare)
		case node.Between != nil:
			return v.VisitSimpleExpression(node.Between)
		case node.In != nil:
			return v.VisitSimpleExpression(node.In)
		default:
			panic("invalid rhs")
		}
	case *parser.Compare:
		return node.Operator + " " + v.VisitSimpleExpression(node.Operand)
	case *parser.Between:
		return fmt.Sprintf("BETWEEN %s AND %s",
			v.VisitSimpleExpression(node.Start), v.VisitSimpleExpression(node.End))
	case *parser.In:
		return fmt.Sprintf(" IN (%s)", v.VisitSimpleExpression(node.Values))
	case *parser.Operand:
		if node.SymbolRef != nil {
			return node.SymbolRef.Symbol
		}
		return v.VisitSimpleExpression(node.Value)
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
		panic("visitor does not recognize type " + reflect.TypeOf(node).String())
	}
}

func buildProjectionExpression(expr *parser.ProjectionExpression) (string, error) {
	cols := make([]*parser.DocumentPath, 0, len(expr.Columns))
	for _, col := range expr.Columns {
		if col.DocumentPath != nil {
			cols = append(cols, col.DocumentPath)
		} else if col.Function != nil {
			fc, err := extractProjectionsFromFunction(col.Function)
			if err != nil {
				return "", err
			}
			cols = append(cols, fc...)
		} else {
			return "", fmt.Errorf("unexpected ProjectionColumn %v", *col)
		}
	}
	buf := &bytes.Buffer{}
	for i, col := range cols {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(col.String())
	}
	return buf.String(), nil
}

func extractProjectionsFromFunction(expr *parser.FunctionExpression) ([]*parser.DocumentPath, error) {
	if expr.Function != "document" {
		return nil, fmt.Errorf("function %q not allowed in projection", expr.Function)
	}
	cols := make([]*parser.DocumentPath, 0, len(expr.Args))
	for _, arg := range expr.Args {
		if arg.DocumentPath == nil {
			return nil, fmt.Errorf("args to document() must be document paths, got %s", arg.String())
		}
		cols = append(cols, arg.DocumentPath)
	}
	return cols, nil
}
