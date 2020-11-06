package parser

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mightyguava/dynamosql/testing/testutil"
)

func TestGoldenGoodQueries(t *testing.T) {
	type row struct {
		Query string
		AST   AST
	}

	queries, err := os.Open("testdata/queries.sql")
	require.NoError(t, err)
	defer queries.Close()
	scanner := bufio.NewScanner(queries)
	var parsed []row
	for scanner.Scan() {
		var ast AST
		query := scanner.Text()
		if strings.HasPrefix(query, "--") {
			// skip comments
			continue
		}
		err := Parser.ParseString(query, &ast)
		assert.NoError(t, err, "Parse: %s", query)
		parsed = append(parsed, row{
			Query: query,
			AST:   ast,
		})
	}

	g := goldie.New(t,
		goldie.WithFixtureDir("testdata/golden"),
		goldie.WithNameSuffix(".golden.go"))
	for i, q := range parsed {
		t.Logf("Query: %s", q.Query)
		g.Assert(t, fmt.Sprintf("queries.%02d", i), []byte(testutil.Repr(q)))
	}
}

func TestGoldenBadQueries(t *testing.T) {
	flag.Parse()

	type row struct {
		Query string
		Error string
	}

	queries, err := os.Open("testdata/bad_queries.sql")
	require.NoError(t, err)
	defer queries.Close()
	scanner := bufio.NewScanner(queries)
	var parsed []row
	for scanner.Scan() {
		var ast AST
		query := scanner.Text()
		if strings.HasPrefix(query, "--") {
			// skip comments
			continue
		}
		err := Parser.ParseString(query, &ast)
		require.Errorf(t, err, "Parse %s, expected error but did not", query)
		parsed = append(parsed, row{
			Query: query,
			Error: err.Error(),
		})
	}

	g := goldie.New(t,
		goldie.WithFixtureDir("testdata/golden"),
		goldie.WithNameSuffix(".golden.go"))
	for i, q := range parsed {
		g.Assert(t, fmt.Sprintf("bad_queries.%02d", i), []byte(testutil.MarshalJSON(q)))
	}
}
