// nolint: govet
package parser

type DropTable struct {
	Table string `"DROP" "TABLE" @( Ident ( "." Ident )* | QuotedIdent )`
}

func (d *DropTable) children() []Node { return nil }
