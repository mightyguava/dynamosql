package querybuilder

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/repr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/mightyguava/dynamosql/parser"
	"github.com/mightyguava/dynamosql/schema"
)

type PreparedQuery struct {
	Query *dynamodb.QueryInput
}

type KeyExpression struct{}
type FilterExpression struct{}

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

func buildQuery(table *schema.Table, ast parser.Select) (*PreparedQuery, error) {
	visit := &visitor{Context: &Context{Table: table}}
	result, err := visit.VisitNodes(ast.Where)
	if err != nil {
		return nil, err
	}
	expressionValues := make(map[string]*dynamodb.AttributeValue)
	req := &dynamodb.QueryInput{
		TableName:                 &ast.From,
		KeyConditionExpression:    aws.String(result.Key),
		ExpressionAttributeValues: expressionValues,
	}
	if result.Filter != "" {
		req.FilterExpression = aws.String(result.Filter)
	}
	return &PreparedQuery{Query: req}, nil
}

type Context struct {
	Table *schema.Table
}

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
		return fmt.Sprintf(" BETWEEN %s AND %s",
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
			return ":" + *node.PlaceHolder
		case node.Number != nil:
			return strconv.FormatFloat(*node.Number, 'g', -1, 64)
		case node.String != nil:
			return *node.String
		case node.Boolean != nil:
			return strconv.FormatBool(bool(*node.Boolean))
		case node.Null:
			return "NULL"
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
