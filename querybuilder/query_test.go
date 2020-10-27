package querybuilder

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mightyguava/dynamosql/parser"
	"github.com/mightyguava/dynamosql/schema"
	"github.com/mightyguava/dynamosql/testing/fixtures"
)

func TestBuildQuery(t *testing.T) {
	type item struct {
		Query   string
		Request *dynamodb.QueryInput
	}
	table := schema.NewTableFromCreate(fixtures.GameScores.Create)

	queries, err := os.Open("testdata/queries.sql")
	require.NoError(t, err)
	defer queries.Close()
	scanner := bufio.NewScanner(queries)
	var parsed []item
	for scanner.Scan() {
		var ast parser.Select
		queryStr := scanner.Text()
		err := parser.Parser.ParseString(queryStr, &ast)
		assert.NoError(t, err, "Parse: %s", queryStr)
		query, err := buildQuery(table, ast)
		require.NoError(t, err)
		parsed = append(parsed, item{
			Query:   queryStr,
			Request: query.Query,
		})
	}

	g := goldie.New(t,
		goldie.WithDiffEngine(goldie.ColoredDiff),
		goldie.WithFixtureDir("testdata/golden"),
		goldie.WithNameSuffix(".golden.json"))
	for i, q := range parsed {
		g.AssertJson(t, fmt.Sprintf("queries.%02d", i), q)
	}
}
