// nolint: govet
package parser

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

var (
	Lexer = lexer.Must(lexer.Regexp(`(\s+)` +
		`|(?P<Keyword>(?i)SELECT|FROM|WHERE|MINUS|EXCEPT|INTERSECT|ORDER|LIMIT|OFFSET|TRUE|FALSE|NULL|IS|NOT|ANY|SOME|BETWEEN|AND|OR|AS)` +
		`|(?P<Function>(?i)attribute_exists|attribute_not_exists|attribute_type|begins_with|contains|size)` +
		`|(?P<Ident>[a-zA-Z_][a-zA-Z0-9_]*)` +
		`|(?P<NamedParameter>:[a-zA-Z_][a-zA-Z0-9_]*)` +
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
	Or []*OrCondition `@@ ( "OR" @@ )*`
}

type OrCondition struct {
	And []*Condition `@@ ( "AND" @@ )*`
}

type ParenthesizedExpression struct {
	ConditionExpression *ConditionExpression
}

type Condition struct {
	Parenthesized *ConditionExpression `  "(" @@ ")"`
	Not           *Condition           `| "NOT" @@`
	Operand       *ConditionOperand    `| @@`
	Function      *FunctionExpression  `| @@`
}

type FunctionExpression struct {
	Function      string  `@Function`
	PathArgument  string  `"(" @Ident`
	MoreArguments []Value `    ( "," @@ )* ")"`
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
	PlaceHolder *string  `  @NamedParameter`
	Number      *float64 `| @Number`
	String      *string  `| @String`
	Boolean     *Boolean `| @("TRUE" | "FALSE")`
	Null        bool     `| @"NULL"`
}
