// nolint: govet
package parser

import (
	"bytes"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

var (
	Lexer = lexer.Must(lexer.Regexp(`(\s+)` +
		`|(?P<Keyword>(?i)SELECT|FROM|WHERE|MINUS|EXCEPT|INTERSECT|ORDER|LIMIT|OFFSET|TRUE|FALSE|NULL|IS|NOT|ANY|SOME|BETWEEN|AND|OR|AS)` +
		`|(?P<Ident>[a-zA-Z_][a-zA-Z0-9_]*)` +
		`|(?P<IndexArray>(\[\d\])+)` +
		`|(?P<Number>[-+]?\d*\.?\d+([eE][-+]?\d+)?)` +
		`|(?P<String>'[^']*'|"[^"]*")` +
		`|(?P<Operators><>|!=|<=|>=|[-+*/%:,.()=<>])`,
	))
	Parser = participle.MustBuild(
		&Select{},
		participle.Lexer(Lexer),
		participle.Unquote("String"),
		participle.CaseInsensitive("Keyword"),
	)
)

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "TRUE"
	return nil
}

// Select based on http://www.h2database.com/html/grammar.html
type Select struct {
	Projection *ProjectionExpression `"SELECT" @@`
	From       string                `"FROM" @Ident ( @"." @Ident )*`
	Where      *AndExpression        `( "WHERE" @@ )?`
	Limit      *int                  `( "LIMIT" @Number )?`
}

type ProjectionExpression struct {
	All         bool           `  @"*"`
	Projections []DocumentPath `| @@ ( "," @@ )*`
}

func (e ProjectionExpression) String() string {
	if e.All {
		return ""
	}
	buf := &bytes.Buffer{}
	buf.WriteString(e.Projections[0].String())
	for _, p := range e.Projections[1:] {
		buf.WriteString(", ")
		buf.WriteString(p.String())
	}
	return buf.String()
}

type ConditionExpression struct {
	Or []*AndExpression `@@ ( "OR" @@ )*`
}

type AndExpression struct {
	And []*Condition `@@ ( "AND" @@ )*`
}

type ParenthesizedExpression struct {
	ConditionExpression *ConditionExpression `@@`
}

type Condition struct {
	Parenthesized *ParenthesizedExpression `  "(" @@ ")"`
	Not           *NotCondition            `| "NOT" @@`
	Operand       *ConditionOperand        `| @@`
	Function      *FunctionExpression      `| @@`
}

type NotCondition struct {
	Condition *Condition `@@`
}

type FunctionExpression struct {
	Function      string   `@Ident`
	PathArgument  string   `"(" @Ident`
	MoreArguments []*Value `    ( "," @@ )* ")"`
}

type ConditionOperand struct {
	Operand      *SymbolRef    `@@`
	ConditionRHS *ConditionRHS `@@`
}

type ConditionRHS struct {
	Compare *Compare `  @@`
	Between *Between `| "BETWEEN" @@`
	In      *In      `| "IN" "(" @@ ")"`
}

type In struct {
	Values []Value `@@ ( "," @@ )*`
}

type Compare struct {
	Operator string   `@( "<>" | "<=" | ">=" | "=" | "<" | ">" | "!=" )`
	Operand  *Operand `@@`
}

type Between struct {
	Start *Operand `@@`
	End   *Operand `"AND" @@`
}

type Operand struct {
	Value     *Value     `  @@`
	SymbolRef *SymbolRef `| @@`
}

type DocumentPath struct {
	Fragment []PathFragment `@@ ( "." @@)*`
}

func (p DocumentPath) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString(p.Fragment[0].Symbol)
	for _, f := range p.Fragment[1:] {
		buf.WriteRune('.')
		buf.WriteString(f.Symbol)
	}
	return buf.String()
}

type SymbolRef struct {
	Symbol string `@Ident @{ "." Ident }`
}

type PathFragment struct {
	Symbol string `@Ident @IndexArray*`
}

type Value struct {
	PlaceHolder *string  `  @":" @Ident`
	Number      *float64 `| @Number`
	String      *string  `| @String`
	Boolean     *Boolean `| @("TRUE" | "FALSE")`
	Null        bool     `| @"NULL"`
}
