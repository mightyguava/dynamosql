// nolint: govet
package main

import (
	"github.com/alecthomas/kong"
	"github.com/alecthomas/repr"

	"github.com/mightyguava/dynamosql/parser"
)

var (
	cli struct {
		SQL string `arg:"" required:"" help:"SQL to parse."`
	}
)

func main() {
	ctx := kong.Parse(&cli)
	sql := &parser.Select{}
	err := parser.Parser.ParseString(cli.SQL, sql)
	repr.Println(sql, repr.Indent("  "), repr.OmitEmpty(true))
	ctx.FatalIfErrorf(err)
}
