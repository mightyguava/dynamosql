package query

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type PreparedQuery struct {
	Query *dynamodb.QueryInput
}

func PrepareQuery(query string) (*PreparedQuery, error) {
	return nil, nil
}
