package parser

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	require.NoError(t, enc.Encode(parsed))

	g := goldie.New(t,
		goldie.WithDiffEngine(goldie.ColoredDiff),
		goldie.WithFixtureDir("testdata/golden"),
		goldie.WithNameSuffix(".golden.json"))
	for i, q := range parsed {
		g.AssertJson(t, fmt.Sprintf("queries.%02d", i), q)
	}
}
