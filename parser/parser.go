// nolint: govet
package parser

import (
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
	"BETWEEN", "AND", "OR", "USE", "INDEX", "ASC", "DESC", "DROP", "CREATE", "TABLE", "HASH",
	"RANGE", "PROJECTION", "PROVISIONED", "THROUGHPUT", "READ", "WRITE", "GLOBAL", "LOCAL", "INDEX",
	"SECONDARY", "RETURNING", "NONE", "ALL_OLD", "UPDATED_OLD", "ALL_NEW", "UPDATED_NEW", "DELETE", "CHECK",
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
	CreateTable *CreateTable     ` | @@`
	DropTable   *DropTable       ` | @@ ) ";"?`
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
