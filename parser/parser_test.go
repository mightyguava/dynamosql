package parser

import (
	"bufio"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden test update")

func TestGolden(t *testing.T) {
	flag.Parse()

	type row struct {
		Query string
		AST   Select
	}

	queries, err := os.Open("testdata/queries.sql")
	require.NoError(t, err)
	defer queries.Close()
	scanner := bufio.NewScanner(queries)
	var parsed []row
	for scanner.Scan() {
		var ast Select
		query := scanner.Text()
		err := Parser.ParseString(query, &ast)
		assert.NoError(t, err, "Parse: %s", query)
		parsed = append(parsed, row{
			Query: query,
			AST:   ast,
		})
	}

	if *update {
		f, err := os.Create("testdata/queries.ast.go")
		require.NoError(t, err)
		repr.New(f).Print(parsed)
	} else {
		data, err := ioutil.ReadFile("testdata/queries.ast.json")
		require.NoError(t, err)
		var expected []Select
		require.NoError(t, json.Unmarshal(data, &expected))
		require.Equal(t, expected, parsed)
	}
}
