// nolint: govet
package querybuilder

import (
	"context"
	"database/sql/driver"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/mightyguava/dynamosql/parser"
)

func PrepareDropTable(drop *parser.DropTable) (ExecStmt, error) {
	return execStatementFunc(func(ctx context.Context, dynamo dynamodbiface.DynamoDBAPI, args []driver.NamedValue) (*DriverResult, error) {
		_, err := dynamo.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: &drop.Table,
		})
		return &DriverResult{count: 0}, err
	}), nil
}
