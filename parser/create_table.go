// nolint: govet
package parser

type CreateTable struct {
	Table                 string                 `"CREATE" "TABLE" @( Ident ( "." Ident )* | QuotedIdent ) "("`
	Entries               []*CreateTableEntry    `@@ ("," @@)* ")"`
	ProvisionedThroughput *ProvisionedThroughput `@@`
}

func (c *CreateTable) node() {}

type CreateTableEntry struct {
	GlobalSecondaryIndex *GlobalSecondaryIndex `  @@`
	LocalSecondaryIndex  *LocalSecondaryIndex  `| @@`
	Attr                 *TableAttr            `| @@` // Must be last.
}

func (c *CreateTableEntry) node() {}

type ProvisionedThroughput struct {
	ReadCapacityUnits  int64 `"PROVISIONED" "THROUGHPUT" "READ" @Number`
	WriteCapacityUnits int64 `"WRITE" @Number`
}

func (p *ProvisionedThroughput) node() {}

type GlobalSecondaryIndex struct {
	Name                  string                 `"GLOBAL" "SECONDARY"? "INDEX" @(Ident | QuotedIdent)`
	PartitionKey          string                 `"HASH" "(" @(Ident | QuotedIdent) ")"`
	SortKey               string                 `("RANGE" "(" @(Ident | QuotedIdent) ")")?`
	Projection            *Projection            `"PROJECTION" @@`
	ProvisionedThroughput *ProvisionedThroughput `@@`
}

func (c *GlobalSecondaryIndex) node() {}

type Projection struct {
	KeysOnly bool     `  @"KEYS_ONLY"`
	All      bool     `| @"ALL"`
	Include  []string `| "INCLUDE" (@(Ident | QuotedIdent) ("," (@(Ident | QuotedIdent)))*)`
}

func (p *Projection) node() {}

type LocalSecondaryIndex struct {
	Name       string      `"LOCAL" "SECONDARY"? "INDEX" @(Ident | QuotedIdent)`
	SortKey    string      `"RANGE" "(" @( Ident ( "." Ident )* | QuotedIdent ) ")"`
	Projection *Projection `"PROJECTION" @@`
}

func (c *LocalSecondaryIndex) node() {}

type TableAttr struct {
	Name string `@(Ident | QuotedIdent)`
	Type string `@Type`
	Key  string `(@("HASH" | "RANGE") "KEY")?`
}

func (c *TableAttr) node() {}
