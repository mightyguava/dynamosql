package driver

import (
	"context"
	"database/sql/driver"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/mightyguava/dynamosql/parser"
)

type conn struct {
	dynamo *dynamodb.DynamoDB
}

var (
	_ driver.Conn               = &conn{}
	_ driver.ExecerContext      = &conn{}
	_ driver.QueryerContext     = &conn{}
	_ driver.ConnPrepareContext = &conn{}
)

func (c conn) Prepare(query string) (driver.Stmt, error) {
	panic("implement me")
}

func (c conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	panic("implement me")
}

func (c conn) Close() error {
	panic("implement me")
}

func (c conn) Begin() (driver.Tx, error) {
	panic("implement me")
}

func (c conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	var ast parser.Select
	if err := parser.Parser.ParseString(query, &ast); err != nil {
		return nil, err
	}
	req := &dynamodb.QueryInput{
		TableName:              aws.String(ast.From.Table),
		KeyConditionExpression: aws.String("title = :title"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":title": {S: aws.String("The Godfather")},
		},
	}
	resp, err := c.dynamo.QueryWithContext(ctx, req)
	if err != nil {
		return nil, err
	}
	return &Rows{
		ast:  ast,
		req:  req,
		resp: resp,
	}, nil
}

func (c conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	panic("implement me")
}
