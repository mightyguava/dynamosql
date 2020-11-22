package parser

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mightyguava/dynamosql/testing/testutil"
)

func TestGoldenGoodQueries(t *testing.T) {
	type row struct {
		Query string
		AST   *AST
	}

	queries, err := os.Open("testdata/queries.sql")
	require.NoError(t, err)
	defer queries.Close()
	scanner := bufio.NewScanner(queries)
	var parsed []row
	for scanner.Scan() {
		query := scanner.Text()
		if strings.HasPrefix(query, "--") {
			// skip comments
			continue
		}
		ast, err := Parse(query)
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
		query := scanner.Text()
		if strings.HasPrefix(query, "--") {
			// skip comments
			continue
		}
		_, err := Parse(query)
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

func TestParseInsert(t *testing.T) {
	tests := []struct {
		name  string
		query string
		ast   *AST
	}{
		{
			name:  "basic",
			query: "INSERT INTO `movies` VALUES (?)",
			ast: &AST{Insert: &Insert{
				Into:   "movies",
				Values: []*InsertTerminal{{Value: Value{PositionalPlaceholder: true}}},
			}},
		},
		{
			name:  "basic replace",
			query: "REPLACE INTO `movies` VALUES (?)",
			ast: &AST{Replace: &Insert{
				Into:   "movies",
				Values: []*InsertTerminal{{Value: Value{PositionalPlaceholder: true}}},
			}},
		},
		{
			name:  "replace returning",
			query: "REPLACE INTO `movies` VALUES (?) RETURNING ALL_OLD",
			ast: &AST{Replace: &Insert{
				Into:      "movies",
				Values:    []*InsertTerminal{{Value: Value{PositionalPlaceholder: true}}},
				Returning: aws.String("ALL_OLD"),
			}},
		},
		{
			name: "literal",
			query: `
INSERT INTO movies
VALUES ('{"title":"hello","year":2009}'),
       ('{"title":"foo","year":2938}');
`,
			ast: &AST{
				Insert: &Insert{
					Into: "movies",
					Values: []*InsertTerminal{
						{
							Value: Value{Scalar: Scalar{Str: aws.String(`{"title":"hello","year":2009}`)}},
						},
						{
							Value: Value{Scalar: Scalar{Str: aws.String(`{"title":"foo","year":2938}`)}},
						},
					},
				},
			},
		},
		{
			name: "object",
			query: `
INSERT INTO movies
VALUES ({title:"hello",year:2009}),
       ({title:"foo",year:2938});
`,
			ast: &AST{
				Insert: &Insert{
					Into: "movies",
					Values: []*InsertTerminal{
						{
							Object: &JSONObject{[]*JSONObjectEntry{
								{"title", &JSONValue{Scalar: Scalar{Str: aws.String("hello")}}},
								{"year", &JSONValue{Scalar: Scalar{Number: aws.Float64(2009)}}},
							}},
						},
						{
							Object: &JSONObject{[]*JSONObjectEntry{
								{"title", &JSONValue{Scalar: Scalar{Str: aws.String("foo")}}},
								{"year", &JSONValue{Scalar: Scalar{Number: aws.Float64(2938)}}},
							}},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Log("Query: ", test.query)
			ast, err := Parse(test.query)
			require.NoError(t, err)
			assert.Equal(t, test.ast, ast)
		})
	}
}
