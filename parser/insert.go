// nolint: govet
package parser

type InsertOrReplace struct {
	Replace   bool              `("INSERT" | @"REPLACE")`
	Into      string            `"INTO" @( Ident ( "." Ident )* | QuotedIdent )`
	Values    []*InsertTerminal `"VALUES" "(" @@ ")" ( "," "(" @@ ")" )* `
	Returning *string           `( "RETURNING" @( "NONE" | "ALL_OLD" ) )?`
}

func (i *InsertOrReplace) children() (children []Node) {
	for _, value := range i.Values {
		children = append(children, value)
	}
	return
}

type InsertTerminal struct {
	Value
	Object *JSONObject `| @@`
}

func (i *InsertTerminal) children() (children []Node) {
	return append(i.Value.children(), i.Object)
}
