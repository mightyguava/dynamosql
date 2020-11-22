package schema

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
		expectedTable := &Table{
			Name:    "movies",
			HashKey: "title",
			SortKey: "year",
		}
		require.Equal(t, expectedTable, table)
		table = NewTableFromCreate(fixtures.Movies.Create)
		require.Equal(t, expectedTable, table)
	}

	{
		desc, err := client.DescribeTable(&dynamodb.DescribeTableInput{TableName: &fixtures.GameScores.Table})
		require.NoError(t, err)
		table := NewTable(desc.Table)
		expectedTable := &Table{
			Name:    "gamescores",
			HashKey: "UserId",
			SortKey: "GameTitle",
			Indexes: []Index{
				{
					Name:    "UserWinsIndex",
					HashKey: "UserId",
					SortKey: "Wins",
				},
				{
					Name:    "GameTitleIndex",
					HashKey: "GameTitle",
					SortKey: "TopScore",
					Global:  true,
				},
			},
		}
		require.Equal(t, expectedTable, table)
		table = NewTableFromCreate(fixtures.GameScores.Create)
		require.Equal(t, expectedTable, table)
	}
}

func TestLoadTable(t *testing.T) {
	sess := fixtures.SetUp(t, fixtures.Movies)
	client := dynamodb.New(sess)

	loader := TableLoader{dynamo: client}
	table, err := loader.Get(context.Background(), fixtures.Movies.Table)
	require.NoError(t, err)

	expectedTable := &Table{
		Name:    "movies",
		HashKey: "title",
		SortKey: "year",
	}
	require.Equal(t, expectedTable, table)
}
