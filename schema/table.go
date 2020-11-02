package schema

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"golang.org/x/sync/singleflight"
)

// Table contains the schema for a DynamoDB table
type Table struct {
	HashKey string
	SortKey string
	Indexes []Index
}

// NewTable parses a dynamodb.TableDescription into a simplified Table schema
// nolint: dupl
func NewTable(desc *dynamodb.TableDescription) *Table {
	var indexes []Index
	indexCount := len(desc.LocalSecondaryIndexes) + len(desc.GlobalSecondaryIndexes)
	if indexCount > 0 {
		indexes = make([]Index, 0, indexCount)
	}
	for _, indexDesc := range desc.LocalSecondaryIndexes {
		index := Index{
			Name: *indexDesc.IndexName,
		}
		index.HashKey, index.SortKey = parseKeySchema(indexDesc.KeySchema)
		indexes = append(indexes, index)
	}
	for _, indexDesc := range desc.GlobalSecondaryIndexes {
		index := Index{
			Name:   *indexDesc.IndexName,
			Global: true,
		}
		index.HashKey, index.SortKey = parseKeySchema(indexDesc.KeySchema)
		indexes = append(indexes, index)
	}
	hash, sort := parseKeySchema(desc.KeySchema)
	return &Table{
		HashKey: hash,
		SortKey: sort,
		Indexes: indexes,
	}
}

// NewTableFromCreate parses a dynamodb.CreateTableInput into a simplified Table schema
// nolint: dupl
func NewTableFromCreate(desc *dynamodb.CreateTableInput) *Table {
	var indexes []Index
	indexCount := len(desc.LocalSecondaryIndexes) + len(desc.GlobalSecondaryIndexes)
	if indexCount > 0 {
		indexes = make([]Index, 0, indexCount)
	}
	for _, indexDesc := range desc.LocalSecondaryIndexes {
		index := Index{
			Name:   *indexDesc.IndexName,
			Global: false,
		}
		index.HashKey, index.SortKey = parseKeySchema(indexDesc.KeySchema)
		indexes = append(indexes, index)
	}
	for _, indexDesc := range desc.GlobalSecondaryIndexes {
		index := Index{
			Name:   *indexDesc.IndexName,
			Global: true,
		}
		index.HashKey, index.SortKey = parseKeySchema(indexDesc.KeySchema)
		indexes = append(indexes, index)
	}
	hash, sort := parseKeySchema(desc.KeySchema)
	return &Table{
		HashKey: hash,
		SortKey: sort,
		Indexes: indexes,
	}
}

// IsKey returns whether the attribute is a hash or sort key.
func (t *Table) IsKey(name string) bool {
	return t.HashKey == name || t.SortKey == name
}

// HasIndex returns true if the table contains an index with a matching name.
func (t *Table) HasIndex(name string) bool {
	for _, idx := range t.Indexes {
		if idx.Name == name {
			return true
		}
	}
	return false
}

// GetIndex returns an index with a matching name, or nil if not found.
func (t *Table) GetIndex(name string) *Index {
	for _, idx := range t.Indexes {
		if idx.Name == name {
			return &idx
		}
	}
	return nil
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

// Index is the schema for a secondary index.
type Index struct {
	Name    string
	HashKey string
	SortKey string
	Global  bool
}

// TableLoader is a loading cache of DynamoDB table schemas.
type TableLoader struct {
	dynamo dynamodbiface.DynamoDBAPI
	tables sync.Map
	load   singleflight.Group
}

func NewTableLoader(dynamo dynamodbiface.DynamoDBAPI) *TableLoader {
	return &TableLoader{dynamo: dynamo}
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
