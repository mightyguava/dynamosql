package querybuilder

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/alecthomas/repr"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mightyguava/dynamosql/parser"
	"github.com/mightyguava/dynamosql/schema"
	"github.com/mightyguava/dynamosql/testing/fixtures"
	"github.com/mightyguava/dynamosql/testing/testutil"
)

func TestBuildQuery(t *testing.T) {
	type item struct {
		Query    string
		Prepared *PreparedQuery
	}
	moviesTable := schema.NewTableFromCreate(fixtures.Movies.Create)
	gameScoresTable := schema.NewTableFromCreate(fixtures.GameScores.Create)
	getTable := func(q string) *schema.Table {
		if strings.Contains(q, "FROM movies") {
			return moviesTable
		}
		if strings.Contains(q, "FROM gamescores") {
			return gameScoresTable
		}
		t.Fatalf("invalid table for query %s", q)
		return nil
	}

	queries, err := os.Open("testdata/queries.sql")
	require.NoError(t, err)
	defer queries.Close()
	scanner := bufio.NewScanner(queries)
	var parsed []item
	for scanner.Scan() {
		queryStr := scanner.Text()
		if strings.HasPrefix(queryStr, "--") {
			// skip comments
			continue
		}
		ast, err := parser.Parse(queryStr)
		msg := fmt.Sprintf("Parse: %s\n%s", queryStr, repr.String(ast, repr.Indent("  ")))
		require.NoError(t, err, msg)
		query, err := prepare(getTable(queryStr), ast.Select)
		require.NoError(t, err, msg)
		parsed = append(parsed, item{
			Query:    queryStr,
			Prepared: query,
		})
	}

	g := goldie.New(t,
		goldie.WithFixtureDir("testdata/golden"),
		goldie.WithNameSuffix(".golden.go"))
	for i, q := range parsed {
		t.Log("Query: ", q.Query)
		g.Assert(t, fmt.Sprintf("queries.%02d", i), []byte(testutil.Repr(q)))
	}
}

func TestInvalidQueries(t *testing.T) {
	type item struct {
		Query string
		Error string
	}
	table := schema.NewTableFromCreate(fixtures.GameScores.Create)

	queries, err := os.Open("testdata/queries_invalid.sql")
	require.NoError(t, err)
	defer queries.Close()
	scanner := bufio.NewScanner(queries)
	var parsed []item
	for scanner.Scan() {
		queryStr := scanner.Text()
		if strings.HasPrefix(queryStr, "--") {
			// skip comments
			continue
		}
		ast, err := parser.Parse(queryStr)
		if err == nil {
			var q *PreparedQuery
			q, err = prepare(table, ast.Select)
			assert.Error(t, err, "Query: %s\nPrepared Query: %s", queryStr, testutil.MarshalJSON(q))
		}
		var errStr string
		if err != nil {
			errStr = err.Error()
		}
		parsed = append(parsed, item{
			Query: queryStr,
			Error: errStr,
		})
	}

	g := goldie.New(t,
		goldie.WithFixtureDir("testdata/golden"),
		goldie.WithNameSuffix(".golden.json"))
	g.Assert(t, "queries_invalid", []byte(testutil.MarshalJSON(parsed)))
}

func TestSubstitute(t *testing.T) {
	ctx := NewContext(&schema.Table{}, "")
	require.Equal(t, "hello", ctx.substitute("hello"))
	require.Equal(t, "#select", ctx.substitute("select"))
	require.Equal(t, "#select", ctx.substitute("select"))
	require.Equal(t, "#_gen1", ctx.substitute("foo.bar"))
	require.Equal(t, "#_gen2", ctx.substitute("foo.bar"))
	require.Equal(t, "#_gen3", ctx.substitute("foo.bar2"))
}
