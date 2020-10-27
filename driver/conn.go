package driver

import (
	"context"
	"database/sql/driver"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/mightyguava/dynamosql/querybuilder"
	"github.com/mightyguava/dynamosql/schema"
)

type conn struct {
	dynamo *dynamodb.DynamoDB
	tables *schema.TableLoader
}

var (
	_ driver.Conn               = &conn{}
	_ driver.ExecerContext      = &conn{}
	_ driver.QueryerContext     = &conn{}
	_ driver.ConnPrepareContext = &conn{}
)

func (c conn) Prepare(query string) (driver.Stmt, error) {
	return c.PrepareContext(context.Background(), query)
}

func (c conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	panic("implement me")
}

func (c conn) Close() error {
	return nil
}

func (c conn) Begin() (driver.Tx, error) {
	panic("implement me")
}

func (c conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	q, err := querybuilder.PrepareQuery(ctx, c.tables, query)
	if err != nil {
		return nil, err
	}
	req, err := q.Build(args)
	if err != nil {
		return nil, err
	}
	resp, err := c.dynamo.QueryWithContext(ctx, req)
	if err != nil {
		return nil, err
	}
	return &Rows{
		req:  q.Query,
		resp: resp,
	}, nil
}

func (c conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	panic("implement me")
}
