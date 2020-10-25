package query

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"golang.org/x/sync/singleflight"
)

// Table contains the schema for a DynamoDB table
type Table struct {
	HashKey                string
	SortKey                string
	LocalSecondaryIndexes  []LocalSecondaryIndex
	GlobalSecondaryIndexes []GlobalSecondaryIndex

	Desc *dynamodb.TableDescription
}

// NewTable parses a dynamodb.TableDescription into a simplified Table schema
func NewTable(desc *dynamodb.TableDescription) *Table {
	var lsi []LocalSecondaryIndex
	if len(desc.LocalSecondaryIndexes) > 0 {
		lsi = make([]LocalSecondaryIndex, 0, len(desc.LocalSecondaryIndexes))
		for _, indexDesc := range desc.LocalSecondaryIndexes {
			index := LocalSecondaryIndex{
				Name: *indexDesc.IndexName,
			}
			_, index.SortKey = parseKeySchema(indexDesc.KeySchema)
			lsi = append(lsi, index)
		}
	}
	var gsi []GlobalSecondaryIndex
	if len(desc.GlobalSecondaryIndexes) > 0 {
		gsi = make([]GlobalSecondaryIndex, 0, len(desc.GlobalSecondaryIndexes))
		for _, indexDesc := range desc.GlobalSecondaryIndexes {
			index := GlobalSecondaryIndex{
				Name: *indexDesc.IndexName,
			}
			index.HashKey, index.SortKey = parseKeySchema(indexDesc.KeySchema)
			gsi = append(gsi, index)
		}
	}
	hash, sort := parseKeySchema(desc.KeySchema)
	return &Table{
		HashKey:                hash,
		SortKey:                sort,
		LocalSecondaryIndexes:  lsi,
		GlobalSecondaryIndexes: gsi,
		Desc:                   desc,
	}
}

func parseKeySchema(schema []*dynamodb.KeySchemaElement) (hash, sort string) {
	for _, key := range schema {
		if *key.KeyType == dynamodb.KeyTypeHash {
			hash = *key.AttributeName
		} else if *key.KeyType == dynamodb.KeyTypeRange {
			sort = *key.AttributeName
		}
	}
	return
}

// LocalSecondaryIndex is the schema for a local secondary index.
type LocalSecondaryIndex struct {
	Name    string
	SortKey string
}

// GlobalSecondaryIndex is the schema for a global secondary index
type GlobalSecondaryIndex struct {
	Name    string
	HashKey string
	SortKey string
}

// TableLoader is a loading cache of DynamoDB table schemas.
type TableLoader struct {
	dynamo *dynamodb.DynamoDB
	tables sync.Map
	load   singleflight.Group
}

// Get retrieves a cached table schema, loading it from DynamoDB if not found.
// If multiple Get are issued against the same table concurrently, only a single request will be made to load the table.
func (l *TableLoader) Get(ctx context.Context, name string) (*Table, error) {
	table, ok := l.tables.Load(name)
	if ok {
		return table.(*Table), nil
	}

	resultChan := l.load.DoChan(name, func() (interface{}, error) {
		desc, err := l.dynamo.DescribeTableWithContext(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(name)})
		if err != nil {
			return nil, err
		}
		table := NewTable(desc.Table)
		l.tables.Store(name, table)
		return table, nil
	})
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-resultChan:
		if result.Err != nil {
			return nil, result.Err
		}
		return result.Val.(*Table), nil
	}
}
