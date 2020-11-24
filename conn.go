package dynamosql

import (
	"context"
	"database/sql/driver"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/mightyguava/dynamosql/parser"
	"github.com/mightyguava/dynamosql/querybuilder"
	"github.com/mightyguava/dynamosql/schema"
)

type conn struct {
	dynamo      dynamodbiface.DynamoDBAPI
	tables      *schema.TableLoader
	mapToGoType bool
}

func (c conn) CheckNamedValue(value *driver.NamedValue) error {
	return nil
}

var (
	_ driver.Conn               = &conn{}
	_ driver.ConnPrepareContext = &conn{}
	_ driver.NamedValueChecker  = &conn{}
)

func (c conn) Prepare(query string) (driver.Stmt, error) {
	panic("unexpected call to Prepare")
}

func (c conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	ast, err := parser.Parse(query)
	if err != nil {
		return nil, err
	}
	switch {
	case ast.Insert != nil:
		stmt, err := querybuilder.PrepareInsert(ctx, c.tables, ast)
		if err != nil {
			return nil, err
		}
		return &execStmt{
			preparedStmt: stmt,
			dynamo:       c.dynamo,
			mapToGoType:  c.mapToGoType,
		}, nil
	case ast.Select != nil:
		prepared, err := querybuilder.PrepareQuery(ctx, c.tables, query)
		if err != nil {
			return nil, err
		}
		return &queryStmt{
			preparedStmt: prepared,
			dynamo:       c.dynamo,
			mapToGoType:  c.mapToGoType,
		}, err
	case ast.CreateTable != nil:
		prepared, err := querybuilder.PrepareCreateTable(ast)
		if err != nil {
			return nil, err
		}
		return &execStmt{
			preparedStmt: prepared,
			dynamo:       c.dynamo,
			mapToGoType:  c.mapToGoType,
		}, err
	case ast.DropTable != nil:
		prepared, err := querybuilder.PrepareDropTable(ast.DropTable)
		if err != nil {
			return nil, err
		}
		return &execStmt{
			preparedStmt: prepared,
			dynamo:       c.dynamo,
			mapToGoType:  c.mapToGoType,
		}, err
	default:
		panic("unsupported statement")
	}
}

func (c conn) Close() error {
	return nil
}

func (c conn) Begin() (driver.Tx, error) {
	panic("implement me")
}
