// nolint: govet
package parser

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

var (
	Lexer = lexer.Must(lexer.Regexp(`(\s+)` +
		`|(?P<Keyword>(?i)SELECT|FROM|WHERE|MINUS|EXCEPT|INTERSECT|ORDER|LIMIT|OFFSET|TRUE|FALSE|NULL|IS|NOT|ANY|SOME|BETWEEN|AND|OR|AS)` +
		`|(?P<Ident>[a-zA-Z_][a-zA-Z0-9_]*)` +
		`|(?P<Number>[-+]?\d*\.?\d+([eE][-+]?\d+)?)` +
		`|(?P<String>'[^']*'|"[^"]*")` +
		`|(?P<Operators><>|!=|<=|>=|[-+*/%,.()=<>])`,
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
	Expression *SelectExpression    `"SELECT" @@`
	From       *From                `"FROM" @@`
	Limit      *ConditionExpression `( "LIMIT" @@ )?`
	Offset     *ConditionExpression `( "OFFSET" @@ )?`
}

type From struct {
	Table string               `@Ident ( @"." @Ident )*`
	Where *ConditionExpression `( "WHERE" @@ )?`
}

type SelectExpression struct {
	All         bool     `  @"*"`
	Projections []string `| @Ident ( "," @Ident )*`
}

type ConditionExpression struct {
	Parenthesized *ParenthesizedExpression `@@`
	Condition     *Condition               `| @@`
	More          []MoreConditions         `  (@@)*`
}

type ParenthesizedExpression struct {
	ConditionExpression *ConditionExpression `"(" @@ ")"`
}

type Condition struct {
	Operand *ConditionOperand `  @@`
	Not     *Condition        `| "NOT" @@`
}

type MoreConditions struct {
	LogicalOperator string     `@("AND" | "OR")`
	Condition       *Condition `@@`
}

type ConditionOperand struct {
	Operand      *Operand      `@@`
	ConditionRHS *ConditionRHS `[ @@ ]`
}

type ConditionRHS struct {
	Compare *Compare `  @@`
	Between *Between `| "BETWEEN" @@`
	In      []Value  `| "(" @@ ( "," @@ )* ")"`
}

type Compare struct {
	Operator string  `@( "<>" | "<=" | ">=" | "=" | "<" | ">" | "!=" )`
	Operand  Operand `@@`
}

type Between struct {
	Start *Operand `@@`
	End   *Operand `"AND" @@`
}

type Operand struct {
	Value     Value      `  @@`
	SymbolRef *SymbolRef `| @@`
}

type SymbolRef struct {
	Symbol string `@Ident @{ "." Ident }`
}

type Value struct {
	Number  *float64 `   @Number`
	String  *string  ` | @String`
	Boolean *Boolean ` | @("TRUE" | "FALSE")`
	Null    bool     ` | @"NULL"`
}
