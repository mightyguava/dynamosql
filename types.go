package dynamosql

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Document returns a sql.Scanner that can scan DynamoDB items into a struct or map using dynamodbattribute.UnmarshalMap
// (https://godoc.org/github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute). Refer to the library docs for
// usage.
func Document(v interface{}) sql.Scanner {
	return documentScanner{v: v}
}

type documentScanner struct {
	v interface{}
}

func (d documentScanner) Scan(src interface{}) error {
	srcMap, ok := src.(map[string]*dynamodb.AttributeValue)
	if !ok {
		return fmt.Errorf("dynamosql.Document() can only be used to Scan a document, not %s", reflect.TypeOf(src))
	}
	return dynamodbattribute.UnmarshalMap(srcMap, d.v)
}
