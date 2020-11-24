// nolint: govet
package parser

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/stateful"
)

var (
	Lexer = stateful.MustSimple([]stateful.Rule{
		{"Whitespace", `\s+`, nil},
		{"Bool", `(?i)\b(TRUE|FALSE)\b`, nil},
		{"Type", `(?i)\b(STRING|NUMBER|BINARY)\b`, nil},
		{"Null", `(?i)\bNULL\b`, nil},
		{"Keyword", keywordsRe(), nil},
		{"QuotedIdent", "`[^`]+`", nil},
		{"Ident", `[a-zA-Z_][a-zA-Z0-9_]*`, nil},
		{"Number", `[-+]?\d*\.?\d+([eE][-+]?\d+)?`, nil},
		{"String", `'[^']*'|"[^"]*"`, nil},
		{"Operator", `<>|!=|<=|>=|[-+*/%:?,.()=<>\[\]{};]`, nil},
	},
	)
	parser = participle.MustBuild(
		&AST{},
		participle.Lexer(Lexer),
		participle.Unquote("String"),
		UnquoteIdent(),
		participle.CaseInsensitive("Keyword", "Bool", "Type", "Null"),
		participle.UseLookahead(2),
		participle.Elide("Whitespace"),
	)
)

var keywords = []string{
	"SELECT", "FROM", "WHERE", "LIMIT", "OFFSET", "INSERT", "INTO", "VALUES", "NOT",
	"BETWEEN", "AND", "OR", "USE", "INDEX", "ASC", "DESC", "CREATE", "TABLE", "HASH", "RANGE", "PROJECTION",
	"PROVISIONED", "THROUGHPUT", "READ", "WRITE", "GLOBAL", "LOCAL", "INDEX", "SECONDARY",
	"RETURNING", "NONE", "ALL_OLD", "UPDATED_OLD", "ALL_NEW", "UPDATED_NEW", "DELETE", "CHECK",
}

func keywordsRe() string {
	return `(?i)\b(` + strings.Join(keywords, "|") + `)\b`
}

func Parse(s string) (*AST, error) {
	var ast AST
	err := parser.ParseString("", s, &ast)
	return &ast, err
}

// EBNF grammar for the SQL parser.
func EBNF() string {
	return parser.String()
}

// UnquoteIdent removes surrounding backticks (`) from quoted identifiers
func UnquoteIdent() participle.Option {
	return participle.Map(func(t lexer.Token) (lexer.Token, error) {
		t.Value = t.Value[1 : len(t.Value)-1]
		return t, nil
	}, "QuotedIdent")
}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = strings.ToUpper(values[0]) == "TRUE"
	return nil
}

type ScanDescending bool

func (b *ScanDescending) Capture(values []string) error {
	*b = strings.ToUpper(values[0]) == "DESC"
	return nil
}

// Node is an interface implemented by all AST nodes.
type Node interface {
	node()
}

type AST struct {
	Select      *Select          `(  @@`
	Insert      *InsertOrReplace ` | @@`
	CreateTable *CreateTable     ` | @@ ) ";"?`
}

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

// Select based on http://www.h2database.com/html/grammar.html
type Select struct {
	Projection *ProjectionExpression `"SELECT" @@`
	From       string                `"FROM" @( Ident ( "." Ident )* | QuotedIdent )`
	Index      *string               `( "USE" "INDEX" "(" @Ident ")" )?`
	Where      *AndExpression        `( "WHERE" @@ )?`
	Descending *ScanDescending       `( @"ASC" | @"DESC" )?`
	Limit      *int                  `( "LIMIT" @Number )?`
}

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

func (e *Select) node() {}

type ProjectionExpression struct {
	All     bool                `  ( @"*" | "document" "(" @"*" ")" )`
	Columns []*ProjectionColumn `| @@ ( "," @@ )*`
}

func (e *ProjectionExpression) node() {}

func (e ProjectionExpression) String() string {
	if e.All {
		return ""
	}
	buf := &bytes.Buffer{}
	buf.WriteString(e.Columns[0].String())
	for _, p := range e.Columns[1:] {
		buf.WriteString(", ")
		buf.WriteString(p.String())
	}
	return buf.String()
}

type ProjectionColumn struct {
	Function     *FunctionExpression `  @@`
	DocumentPath *DocumentPath       `| @@`
}

func (c *ProjectionColumn) node() {}

func (c ProjectionColumn) String() string {
	if c.DocumentPath != nil {
		return c.DocumentPath.String()
	}
	if c.Function != nil {
		return c.Function.String()
	}
	return ""
}

type ConditionExpression struct {
	Or []*AndExpression `@@ ( "OR" @@ )*`
}

func (e *ConditionExpression) node() {}

type AndExpression struct {
	And []*Condition `@@ ( "AND" @@ )*`
}

func (e *AndExpression) node() {}

type Condition struct {
	Parenthesized *ConditionExpression `  "(" @@ ")"`
	Not           *NotCondition        `| "NOT" @@`
	Operand       *ConditionOperand    `| @@`
	Function      *FunctionExpression  `| @@`
}

func (e *Condition) node() {}

type NotCondition struct {
	Condition *Condition `@@`
}

func (e *NotCondition) node() {}

type FunctionExpression struct {
	Function string              `@Ident`
	Args     []*FunctionArgument `"(" @@ ( "," @@ )* ")"`
}

func (f *FunctionExpression) node() {}

func (f *FunctionExpression) FirstArgIsRef() bool {
	return len(f.Args) > 0 && f.Args[0].DocumentPath != nil
}

func (f *FunctionExpression) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString(f.Function)
	buf.WriteRune('(')
	for i, arg := range f.Args {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.String())
	}
	buf.WriteRune(')')
	return buf.String()
}

type FunctionArgument struct {
	DocumentPath *DocumentPath `  @@`
	Value        *Value        `| @@`
}

func (a *FunctionArgument) node() {}

func (a FunctionArgument) String() string {
	if a.DocumentPath != nil {
		return a.DocumentPath.String()
	}
	if a.Value != nil {
		return a.Value.String()
	}
	return ""
}

type ConditionOperand struct {
	Operand      *DocumentPath `@@`
	ConditionRHS *ConditionRHS `@@`
}

func (c *ConditionOperand) node() {}

type ConditionRHS struct {
	Compare *Compare `  @@`
	Between *Between `| "BETWEEN" @@`
	In      *In      `| "IN" "(" @@ ")"`
}

func (c *ConditionRHS) node() {}

type In struct {
	Values []*Value `@@ ( "," @@ )*`
}

func (i *In) node() {}

type Compare struct {
	Operator string   `@( "<>" | "<=" | ">=" | "=" | "<" | ">" | "!=" )`
	Operand  *Operand `@@`
}

func (c *Compare) node() {}

type Between struct {
	Start *Operand `@@`
	End   *Operand `"AND" @@`
}

func (b *Between) node() {}

type Operand struct {
	Value     *Value        `  @@`
	SymbolRef *DocumentPath `| @@`
}

func (o *Operand) node() {}

type DocumentPath struct {
	Fragment []*PathFragment `@@ ( "." @@ )*`
}

func (p *DocumentPath) node() {}

// String marshals the DocumentPath into a human readable format. Do not use this function when marshaling
// to expressions, because substitutions need to be applied first for reserved words.
func (p DocumentPath) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString(p.Fragment[0].String())
	for _, f := range p.Fragment[1:] {
		buf.WriteRune('.')
		buf.WriteString(f.String())
	}
	return buf.String()
}

type PathFragment struct {
	Symbol  string `( @Ident | @QuotedIdent )`
	Indexes []int  `( "[" @Number "]" )*`
}

func (p PathFragment) node() {}

func (p PathFragment) String() string {
	if len(p.Indexes) == 0 {
		return p.Symbol
	}
	buf := &bytes.Buffer{}
	buf.WriteString(p.Symbol)
	for _, idx := range p.Indexes {
		buf.WriteRune('[')
		buf.WriteString(strconv.Itoa(idx))
		buf.WriteRune(']')
	}
	return buf.String()
}

type JSONObjectEntry struct {
	Key   string     `@(Ident | String)`
	Value *JSONValue `":" @@`
}

func (j *JSONObjectEntry) node() {}

type JSONObject struct {
	Entries []*JSONObjectEntry `"{" (@@ ("," @@)* ","?)? "}"`
}

func (j *JSONObject) node() {}

func (j *JSONObject) String() string {
	out := make([]string, 0, len(j.Entries))
	for _, entry := range j.Entries {
		out = append(out, strconv.Quote(entry.Key)+":"+entry.Value.String())
	}
	return "{" + strings.Join(out, ",") + "}"
}

type JSONArray struct {
	Entries []*JSONValue `"[" (@@ ("," @@)* ","?)? "]"`
}

func (j *JSONArray) node() {}

func (j *JSONArray) String() string {
	out := make([]string, 0, len(j.Entries))
	for _, v := range j.Entries {
		out = append(out, v.String())
	}
	return "[" + strings.Join(out, ",") + "]"
}

type JSONValue struct {
	Scalar
	Object *JSONObject `| @@`
	Array  *JSONArray  `| @@`
}

type Scalar struct {
	Number  *float64 `  @Number`
	Str     *string  `| @String`
	Boolean *Boolean `| @Bool`
	Null    bool     `| @Null`
}

func (l *Scalar) node() {}
func (l *Scalar) String() string {
	switch {
	case l.Number != nil:
		return strconv.FormatFloat(*l.Number, 'g', -1, 64)
	case l.Str != nil:
		return strconv.Quote(*l.Str)
	case l.Boolean != nil:
		return strconv.FormatBool(bool(*l.Boolean))
	case l.Null:
		return "NULL"
	default:
		panic("unexpected code path")
	}
}

type Value struct {
	Scalar
	PlaceHolder           *string `| @":" @Ident `
	PositionalPlaceholder bool    `| @"?" `
}

func (v *Value) node() {}

func (v Value) String() string {
	switch {
	case v.PlaceHolder != nil:
		return *v.PlaceHolder
	case v.PositionalPlaceholder:
		return "?"
	default:
		return v.Scalar.String()
	}
}
