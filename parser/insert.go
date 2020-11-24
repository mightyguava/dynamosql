// nolint: govet
package parser

type InsertOrReplace struct {
	Replace   bool              `("INSERT" | @"REPLACE")`
	Into      string            `"INTO" @( Ident ( "." Ident )* | QuotedIdent )`
	Values    []*InsertTerminal `"VALUES" "(" @@ ")" ( "," "(" @@ ")" )* `
	Returning *string           `( "RETURNING" @( "NONE" | "ALL_OLD" ) )?`
}

type InsertTerminal struct {
	Value
	Object *JSONObject `| @@`
}
