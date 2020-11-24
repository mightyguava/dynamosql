// nolint: govet
package parser

type DropTable struct {
	Table string `"DROP" "TABLE" @( Ident ( "." Ident )* | QuotedIdent )`
}

func (d *DropTable) node() {}
