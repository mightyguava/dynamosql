// nolint: govet
package parser

import (
	"bytes"
	"strconv"
)

// Select based on http://www.h2database.com/html/grammar.html
type Select struct {
	Projection *ProjectionExpression `"SELECT" @@`
	From       string                `"FROM" @( Ident ( "." Ident )* | QuotedIdent )`
	Index      *string               `( "USE" "INDEX" "(" @Ident ")" )?`
	Where      *AndExpression        `( "WHERE" @@ )?`
	Descending *ScanDescending       `( @"ASC" | @"DESC" )?`
	Limit      *int                  `( "LIMIT" @Number )?`
}

func (e *Select) children() (children []Node) {
	return []Node{e.Projection, e.Where}
}

type ProjectionExpression struct {
	All     bool                `  ( @"*" | "document" "(" @"*" ")" )`
	Columns []*ProjectionColumn `| @@ ( "," @@ )*`
}

func (e *ProjectionExpression) children() (children []Node) {
	for _, col := range e.Columns {
		children = append(children, col)
	}
	return
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

func (c *ProjectionColumn) children() (children []Node) {
	return []Node{c.Function, c.DocumentPath}
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

func (e *ConditionExpression) children() (children []Node) {
	for _, or := range e.Or {
		children = append(children, or)
	}
	return
}

type AndExpression struct {
	And []*Condition `@@ ( "AND" @@ )*`
}

func (e *AndExpression) children() (children []Node) {
	for _, and := range e.And {
		children = append(children, and)
	}
	return
}

type Condition struct {
	Parenthesized *ConditionExpression `  "(" @@ ")"`
	Not           *NotCondition        `| "NOT" @@`
	Operand       *ConditionOperand    `| @@`
	Function      *FunctionExpression  `| @@`
}

func (e *Condition) children() (children []Node) {
	return []Node{e.Parenthesized, e.Not, e.Operand, e.Function}
}

type NotCondition struct {
	Condition *Condition `@@`
}

func (e *NotCondition) children() (children []Node) {
	return []Node{e.Condition}
}

type FunctionExpression struct {
	Function string              `@Ident`
	Args     []*FunctionArgument `"(" @@ ( "," @@ )* ")"`
}

func (f *FunctionExpression) children() (children []Node) {
	for _, arg := range f.Args {
		children = append(children, arg)
	}
	return
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

func (a *FunctionArgument) children() (children []Node) {
	return []Node{a.DocumentPath, a.Value}
}

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

func (c *ConditionOperand) children() (children []Node) {
	return []Node{c.Operand, c.ConditionRHS}
}

type ConditionRHS struct {
	Compare *Compare `  @@`
	Between *Between `| "BETWEEN" @@`
	In      *In      `| "IN" "(" @@ ")"`
}

func (c *ConditionRHS) children() (children []Node) {
	return []Node{c.Compare, c.Between, c.In}
}

type In struct {
	Values []*Value `@@ ( "," @@ )*`
}

func (i *In) children() (children []Node) {
	for _, value := range i.Values {
		children = append(children, value)
	}
	return
}

type Compare struct {
	Operator string   `@( "<>" | "<=" | ">=" | "=" | "<" | ">" | "!=" )`
	Operand  *Operand `@@`
}

func (c *Compare) children() (children []Node) {
	return []Node{c.Operand}
}

type Between struct {
	Start *Operand `@@`
	End   *Operand `"AND" @@`
}

func (b *Between) children() (children []Node) {
	return []Node{b.Start, b.End}
}

type Operand struct {
	Value     *Value        `  @@`
	SymbolRef *DocumentPath `| @@`
}

func (o *Operand) children() (children []Node) {
	return []Node{o.Value, o.SymbolRef}
}

type DocumentPath struct {
	Fragment []*PathFragment `@@ ( "." @@ )*`
}

func (p *DocumentPath) children() (children []Node) {
	for _, frag := range p.Fragment {
		children = append(children, frag)
	}
	return
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

func (p PathFragment) children() []Node { return nil }

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
