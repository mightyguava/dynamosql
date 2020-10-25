package query

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/require"

	"github.com/mightyguava/dynamosql/testing/fixtures"
)

func TestNewTable(t *testing.T) {
	sess := fixtures.SetUp(t, fixtures.GameScores, fixtures.Movies)
	client := dynamodb.New(sess)

	{
		desc, err := client.DescribeTable(&dynamodb.DescribeTableInput{TableName: &fixtures.Movies.Table})
		require.NoError(t, err)
		table := NewTable(desc.Table)
		require.Same(t, table.Desc, desc.Table)
		table.Desc = nil
		expectedTable := &Table{
			HashKey: "title",
			SortKey: "year",
		}
		require.Equal(t, expectedTable, table)
	}

	{
		desc, err := client.DescribeTable(&dynamodb.DescribeTableInput{TableName: &fixtures.GameScores.Table})
		require.NoError(t, err)
		table := NewTable(desc.Table)
		require.Same(t, table.Desc, desc.Table)
		table.Desc = nil
		expectedTable := &Table{
			HashKey: "UserId",
			SortKey: "GameTitle",
			LocalSecondaryIndexes: []LocalSecondaryIndex{
				{
					Name:    "UserWinsIndex",
					SortKey: "Wins",
				},
			},
			GlobalSecondaryIndexes: []GlobalSecondaryIndex{
				{
					Name:    "GameTitleIndex",
					HashKey: "GameTitle",
					SortKey: "TopScore",
				},
			},
		}
		require.Equal(t, expectedTable, table)
	}
}

func TestLoadTable(t *testing.T) {
	sess := fixtures.SetUp(t, fixtures.Movies)
	client := dynamodb.New(sess)

	loader := TableLoader{dynamo: client}
	table, err := loader.Get(context.Background(), fixtures.Movies.Table)
	require.NoError(t, err)

	table.Desc = nil
	expectedTable := &Table{
		HashKey: "title",
		SortKey: "year",
	}
	require.Equal(t, expectedTable, table)
}
