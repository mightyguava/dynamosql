package dynamosql

import (
	"database/sql/driver"
	"testing"

	"github.com/alecthomas/participle"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/require"

	"github.com/mightyguava/dynamosql/parser"
)

func TestPluck(t *testing.T) {
	projectionParser := participle.MustBuild(
		&parser.ProjectionColumn{},
		participle.Lexer(parser.Lexer),
	)
	item := map[string]*dynamodb.AttributeValue{
		"field": {S: aws.String("foo")},
		"numberSet": {
			NS: aws.StringSlice([]string{"100", "101"}),
		},
		"binarySet": {
			BS: [][]byte{
				[]byte("bytes1"),
				[]byte("bytes2"),
			},
		},
		"null": {
			NULL: aws.Bool(true),
		},
		"nestedDocument": {M: map[string]*dynamodb.AttributeValue{
			"nestedValue": {S: aws.String("nested")},
			"nestedList": {
				L: []*dynamodb.AttributeValue{
					{N: aws.String("15")},
					{L: []*dynamodb.AttributeValue{
						{SS: aws.StringSlice([]string{"n1", "n2"})},
					}},
				},
			},
		}},
		"list": {
			L: []*dynamodb.AttributeValue{
				{N: aws.String("3")},
				{
					M: map[string]*dynamodb.AttributeValue{
						"deepField": {
							SS: aws.StringSlice([]string{"a", "b", "c"}),
						},
					},
				},
			},
		},
	}
	tests := []struct {
		name   string
		path   string
		result driver.Value
	}{
		{
			name:   "top-level field",
			path:   "field",
			result: "foo",
		},
		{
			name:   "nested document",
			path:   "nestedDocument",
			result: item["nestedDocument"].M,
		},
		{
			name:   "nested value",
			path:   "nestedDocument.nestedValue",
			result: "nested",
		},
		{
			name:   "index nested list",
			path:   "nestedDocument.nestedList[0]",
			result: "15",
		},
		{
			name:   "index deep nested list",
			path:   "nestedDocument.nestedList[1][0][1]",
			result: "n2",
		},
		{
			name:   "number set",
			path:   "numberSet",
			result: aws.StringValueSlice(item["numberSet"].NS),
		},
		{
			name:   "index number set",
			path:   "numberSet[1]",
			result: "101",
		},
		{
			name:   "index out of bounds",
			path:   "numberSet[2]",
			result: nil,
		},
		{
			name:   "list nested field",
			path:   "list[1].deepField[2]",
			result: "c",
		},
		{
			name:   "missing nested field",
			path:   "list[1].deepField[2].missingField",
			result: nil,
		},
		{
			name:   "missing nested index",
			path:   "nestedDocument.nestedValue[2]",
			result: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var projection parser.ProjectionColumn
			require.NoError(t, projectionParser.ParseString(test.path, &projection))
			v := pluck(&dynamodb.AttributeValue{M: item}, projection.DocumentPath)
			require.Equal(t, test.result, v)
		})
	}
}
