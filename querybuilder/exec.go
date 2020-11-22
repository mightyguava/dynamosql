package querybuilder

import (
	"context"
	"database/sql/driver"
	"errors"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type ExecStmt interface {
	Do(ctx context.Context, dynamo dynamodbiface.DynamoDBAPI, args []driver.NamedValue) (*DriverResult, error)
}

type DriverResult struct {
	count    int
	returned map[string]*dynamodb.AttributeValue
}

var _ driver.Result = &DriverResult{}

func (i *DriverResult) LastInsertId() (int64, error) {
	return 0, errors.New("LastInsertId is not supported")
}

func (i *DriverResult) RowsAffected() (int64, error) {
	return int64(i.count), nil
}

func (i *DriverResult) Item() map[string]*dynamodb.AttributeValue {
	return i.returned
}
