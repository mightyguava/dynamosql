// nolint: govet
package parser

import (
	"bytes"
	"strconv"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

var (
	Lexer = lexer.Must(lexer.Regexp(`(\s+)` +
		`|\b(?P<Keyword>(?i)SELECT|FROM|WHERE|LIMIT|OFFSET|TRUE|FALSE|NULL|NOT|BETWEEN|AND|OR|USE|INDEX)\b` +
		"|(?P<QuotedIdent>`[^`]+`)" +
		`|(?P<Ident>[a-zA-Z_][a-zA-Z0-9_]*)` +
		`|(?P<Number>[-+]?\d*\.?\d+([eE][-+]?\d+)?)` +
		`|(?P<String>'[^']*'|"[^"]*")` +
		`|(?P<Operators><>|!=|<=|>=|[-+*/%:,.()=<>\[\]])`,
	))
	Parser = participle.MustBuild(
		&Select{},
		participle.Lexer(Lexer),
		participle.Unquote("String"),
		UnquoteIdent(),
		participle.CaseInsensitive("Keyword"),
		participle.UseLookahead(2),
	)
)

// UnquoteIdent removes surrounding backticks (`) from quoted identifiers
func UnquoteIdent() participle.Option {
	return participle.Map(func(t lexer.Token) (lexer.Token, error) {
		t.Value = t.Value[1 : len(t.Value)-1]
		return t, nil
	}, "QuotedIdent")
}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "TRUE"
	return nil
}

// Select based on http://www.h2database.com/html/grammar.html
type Select struct {
	Projection *ProjectionExpression `"SELECT" @@`
	From       string                `"FROM" @Ident`
	Index      *string               `( "USE" "INDEX" "(" @Ident ")" )?`
	Where      *AndExpression        `( "WHERE" @@ )?`
	Limit      *int                  `( "LIMIT" @Number )?`
}

type ProjectionExpression struct {
	All     bool                `  ( @"*" | "document" "(" @"*" ")" )`
	Columns []*ProjectionColumn `| @@ ( "," @@ )*`
}

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
	Function string              `@Ident`
	Args     []*FunctionArgument `"(" @@ ( "," @@ )* ")"`
}

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

func (a FunctionArgument) String() string {
	if a.DocumentPath != nil {
		return a.DocumentPath.String()
	}
	if a.Value != nil {
		return a.Value.Literal()
	}
	return ""
}

type ConditionOperand struct {
	Operand      *DocumentPath `@@`
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
	Value     *Value        `  @@`
	SymbolRef *DocumentPath `| @@`
}

type DocumentPath struct {
	Fragment []PathFragment `@@ ( "." @@ )*`
}

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

type Value struct {
	PlaceHolder *string  `  @":" @Ident`
	Number      *float64 `| @Number`
	String      *string  `| @String`
	Boolean     *Boolean `| @("TRUE" | "FALSE")`
	Null        bool     `| @"NULL"`
}

func (v Value) Literal() string {
	switch {
	case v.PlaceHolder != nil:
		return *v.PlaceHolder
	case v.Number != nil:
		return strconv.FormatFloat(*v.Number, 'g', -1, 64)
	case v.String != nil:
		return *v.String
	case v.Boolean != nil:
		return strconv.FormatBool(bool(*v.Boolean))
	case v.Null:
		return "NULL"
	default:
		panic("unexpected code path")
	}
}
