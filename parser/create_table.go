// nolint: govet
package parser

type CreateTable struct {
	Table                 string                 `"CREATE" "TABLE" @( Ident ( "." Ident )* | QuotedIdent ) "("`
	Entries               []*CreateTableEntry    `@@ ("," @@)* ")"`
	ProvisionedThroughput *ProvisionedThroughput `@@`
}

func (c *CreateTable) children() (children []Node) {
	for _, entry := range c.Entries {
		children = append(children, entry)
	}
	children = append(children, c.ProvisionedThroughput)
	return
}

type CreateTableEntry struct {
	GlobalSecondaryIndex *GlobalSecondaryIndex `  @@`
	LocalSecondaryIndex  *LocalSecondaryIndex  `| @@`
	Attr                 *TableAttr            `| @@` // Must be last.
}

func (c *CreateTableEntry) children() (children []Node) {
	return []Node{c.GlobalSecondaryIndex, c.LocalSecondaryIndex, c.Attr}
}

type ProvisionedThroughput struct {
	ReadCapacityUnits  int64 `"PROVISIONED" "THROUGHPUT" "READ" @Number`
	WriteCapacityUnits int64 `"WRITE" @Number`
}

func (p *ProvisionedThroughput) children() []Node { return nil }

type GlobalSecondaryIndex struct {
	Name                  string                 `"GLOBAL" "SECONDARY"? "INDEX" @(Ident | QuotedIdent)`
	PartitionKey          string                 `"HASH" "(" @(Ident | QuotedIdent) ")"`
	SortKey               string                 `("RANGE" "(" @(Ident | QuotedIdent) ")")?`
	Projection            *Projection            `"PROJECTION" @@`
	ProvisionedThroughput *ProvisionedThroughput `@@`
}

func (c *GlobalSecondaryIndex) children() (children []Node) {
	return []Node{c.Projection, c.ProvisionedThroughput}
}

type Projection struct {
	KeysOnly bool     `  @"KEYS_ONLY"`
	All      bool     `| @"ALL"`
	Include  []string `| "INCLUDE" (@(Ident | QuotedIdent) ("," (@(Ident | QuotedIdent)))*)`
}

func (p *Projection) children() []Node { return nil }

type LocalSecondaryIndex struct {
	Name       string      `"LOCAL" "SECONDARY"? "INDEX" @(Ident | QuotedIdent)`
	SortKey    string      `"RANGE" "(" @( Ident ( "." Ident )* | QuotedIdent ) ")"`
	Projection *Projection `"PROJECTION" @@`
}

func (c *LocalSecondaryIndex) children() (children []Node) {
	return []Node{c.Projection}
}

type TableAttr struct {
	Name string `@(Ident | QuotedIdent)`
	Type string `@Type`
	Key  string `(@("HASH" | "RANGE") "KEY")?`
}

func (c *TableAttr) children() []Node { return nil }
