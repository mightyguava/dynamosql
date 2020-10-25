package querybuilder

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/mightyguava/dynamosql/parser"
	"github.com/mightyguava/dynamosql/schema"
)

type PreparedQuery struct {
	Query *dynamodb.QueryInput
}

func PrepareQuery(ctx context.Context, tables *schema.TableLoader, query string) (*PreparedQuery, error) {
	var ast parser.Select
	if err := parser.Parser.ParseString(query, &ast); err != nil {
		return nil, err
	}
	//table, err := tables.Get(ctx, ast.From)
	//if err != nil {
	//	return nil, err
	//}
	var keyConditionExpr []string
	expressionValues := make(map[string]*dynamodb.AttributeValue)
	if ast.Where != nil {
		for _, condition := range ast.Where.And {
			binaryExpr := condition.Operand
			expr := fmt.Sprintf("%s %s :%s", binaryExpr.Operand.Symbol, binaryExpr.ConditionRHS.Compare.Operator, binaryExpr.Operand.Symbol)
			keyConditionExpr = append(keyConditionExpr, expr)
			expressionValues[":"+binaryExpr.Operand.Symbol] = &dynamodb.AttributeValue{S: binaryExpr.ConditionRHS.Compare.Operand.Value.String}
		}
	}
	req := &dynamodb.QueryInput{
		TableName:                 &ast.From,
		KeyConditionExpression:    aws.String(strings.Join(keyConditionExpr, " AND ")),
		ExpressionAttributeValues: expressionValues,
	}
	return &PreparedQuery{Query: req}, nil
}
