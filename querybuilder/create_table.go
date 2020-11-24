// nolint: govet
package querybuilder

import (
	"context"
	"database/sql/driver"
	"strings"

	"github.com/alecthomas/repr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/mightyguava/dynamosql/parser"
)

func PrepareCreateTable(ast *parser.AST) (ExecStmt, error) {
	return execStatementFunc(func(ctx context.Context, dynamo dynamodbiface.DynamoDBAPI, args []driver.NamedValue) (*DriverResult, error) {
		stmt := ast.CreateTable
		req := &dynamodb.CreateTableInput{
			TableName:             &stmt.Table,
			ProvisionedThroughput: mapProvisionedThroughput(stmt.ProvisionedThroughput),
		}
		for _, entry := range stmt.Entries {
			switch {
			case entry.Attr != nil:
				attr := entry.Attr
				typ := strings.ToUpper(attr.Type[0:1])
				req.AttributeDefinitions = append(req.AttributeDefinitions, &dynamodb.AttributeDefinition{
					AttributeName: &attr.Name,
					AttributeType: &typ,
				})
				if attr.Key != "" {
					req.KeySchema = append(req.KeySchema, &dynamodb.KeySchemaElement{
						AttributeName: &attr.Name,
						KeyType:       aws.String(strings.ToUpper(attr.Key)),
					})
				}

			case entry.GlobalSecondaryIndex != nil:
				gsi := entry.GlobalSecondaryIndex
				reqgsi := &dynamodb.GlobalSecondaryIndex{
					IndexName: &gsi.Name,
					KeySchema: []*dynamodb.KeySchemaElement{
						{AttributeName: &gsi.PartitionKey, KeyType: aws.String("HASH")},
					},
					Projection:            mapProjection(gsi.Projection),
					ProvisionedThroughput: mapProvisionedThroughput(gsi.ProvisionedThroughput),
				}
				if gsi.SortKey != "" {
					reqgsi.KeySchema = append(reqgsi.KeySchema,
						&dynamodb.KeySchemaElement{AttributeName: &gsi.SortKey, KeyType: aws.String("RANGE")},
					)
				}
				req.GlobalSecondaryIndexes = append(req.GlobalSecondaryIndexes, reqgsi)

			case entry.LocalSecondaryIndex != nil:
				lsi := entry.LocalSecondaryIndex
				req.LocalSecondaryIndexes = append(req.LocalSecondaryIndexes, &dynamodb.LocalSecondaryIndex{
					IndexName: &lsi.Name,
					KeySchema: []*dynamodb.KeySchemaElement{
						{AttributeName: &lsi.SortKey, KeyType: aws.String("RANGE")},
					},
					Projection: mapProjection(lsi.Projection),
				})

			default:
				panic(repr.String(entry))
			}
		}
		_, err := dynamo.CreateTableWithContext(ctx, req)
		return &DriverResult{count: 0}, err
	}), nil
}

func mapProjection(projection *parser.Projection) *dynamodb.Projection {
	out := &dynamodb.Projection{}
	switch {
	case projection.All:
		out.ProjectionType = aws.String("ALL")

	case projection.KeysOnly:
		out.ProjectionType = aws.String("KEYS_ONLY")

	default:
		out.ProjectionType = aws.String("INCLUDE")
		for _, attr := range projection.Include {
			attr := attr
			out.NonKeyAttributes = append(out.NonKeyAttributes, &attr)
		}
	}
	return out
}

func mapProvisionedThroughput(throughput *parser.ProvisionedThroughput) *dynamodb.ProvisionedThroughput {
	return &dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  &throughput.ReadCapacityUnits,
		WriteCapacityUnits: &throughput.WriteCapacityUnits,
	}
}
