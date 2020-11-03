package querybuilder

import (
	"bytes"
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/repr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/mightyguava/dynamosql/parser"
	"github.com/mightyguava/dynamosql/schema"
)

var (
	errPositionalArg = errors.New("unexpected positional arg, use sql.NamedArg to pass named arguments")
	errNamedArg      = errors.New("unexpected named arg, to use named args, provided named placeholders like :param")
)

type PreparedQuery struct {
	Query            *dynamodb.QueryInput
	Columns          []*parser.ProjectionColumn
	NamedParams      NamedParams
	PositionalParams map[int]string
	FixedParams      map[string]interface{}
}

func PrepareQuery(ctx context.Context, tables *schema.TableLoader, query string) (*PreparedQuery, error) {
	var ast parser.Select
	if err := parser.Parser.ParseString(query, &ast); err != nil {
		return nil, err
	}
	table, err := tables.Get(ctx, ast.From)
	if err != nil {
		return nil, err
	}
	return prepare(table, ast)
}

func (pq *PreparedQuery) NewRequest(args []driver.NamedValue) (*dynamodb.QueryInput, error) {
	values, err := bindArgs(pq.FixedParams, pq.NamedParams, pq.PositionalParams, args)
	if err != nil {
		return nil, err
	}
	req := *pq.Query
	req.ExpressionAttributeValues = values
	return &req, nil
}

func bindArgs(fixedParams map[string]interface{}, namedParams NamedParams, positionalParams map[int]string, args []driver.NamedValue) (map[string]*dynamodb.AttributeValue, error) {
	values := make(map[string]*dynamodb.AttributeValue, len(namedParams)+len(fixedParams))

	// Bind fixed params
	for k, v := range fixedParams {
		av, err := toAttributeValue(v)
		if err != nil {
			return nil, err
		}
		values[k] = av
	}

	if len(namedParams) > 0 {
		// Bind named params using the args
		namedParams = namedParams.Clone()
		for _, arg := range args {
			if arg.Name == "" {
				return nil, errPositionalArg
			}
			name := ":" + arg.Name
			_, ok := namedParams[name]
			if !ok {
				return nil, fmt.Errorf("binding %q not found", name)
			}
			av, err := toAttributeValue(arg.Value)
			if err != nil {
				return nil, err
			}
			values[name] = av
			delete(namedParams, name)
		}

		if len(namedParams) > 0 {
			for k := range namedParams {
				return nil, fmt.Errorf("missing argument for binding %q", k)
			}
		}
	} else {
		// Bind positional params using the args
		if len(args) != len(positionalParams) {
			return nil, fmt.Errorf("wrong number of arguments, expected %d, got %d", len(positionalParams), len(args))
		}
		for _, arg := range args {
			if arg.Name != "" {
				return nil, errNamedArg
			}
			name := positionalParams[arg.Ordinal]
			av, err := toAttributeValue(arg.Value)
			if err != nil {
				return nil, err
			}
			values[name] = av
		}
	}

	return values, nil
}

func toAttributeValue(attr interface{}) (*dynamodb.AttributeValue, error) {
	switch v := attr.(type) {
	case string:
		return &dynamodb.AttributeValue{S: &v}, nil
	case float64:
		return &dynamodb.AttributeValue{N: aws.String(strconv.FormatFloat(v, 'g', -1, 64))}, nil
	case int64:
		return &dynamodb.AttributeValue{N: aws.String(strconv.FormatInt(v, 10))}, nil
	case bool:
		return &dynamodb.AttributeValue{BOOL: &v}, nil
	case nil:
		return &dynamodb.AttributeValue{NULL: aws.Bool(true)}, nil
	default:
		return nil, fmt.Errorf("invalid value type %s", reflect.TypeOf(v))
	}
}

type KeyExpression struct{}
type FilterExpression struct{}

type NamedParams map[string]Empty

func (p NamedParams) Clone() NamedParams {
	copy := make(NamedParams)
	for k, v := range p {
		copy[k] = v
	}
	return copy
}

func prepare(table *schema.Table, ast parser.Select) (*PreparedQuery, error) {
	index := ""
	if ast.Index != nil {
		index = *ast.Index
		if !table.HasIndex(index) {
			return nil, fmt.Errorf("unrecognized index %q fro table %q", *ast.Index, ast.From)
		}
	}
	ctx := NewContext(table, index)
	visit := &visitor{Context: ctx}
	kf := extractKeyExpressions(ast.Where, ctx.IsKey)
	keyExpr, err := buildKeyExpression(ctx, kf.Key)
	if err != nil {
		return nil, err
	}
	filterExpr, err := buildFilterExpression(ctx, kf.Filter)
	if err != nil {
		return nil, err
	}
	var projectionExpr *string
	if !ast.Projection.All {
		expr, err := buildProjectionExpression(ctx, ast.Projection)
		if err != nil {
			return nil, err
		}
		projectionExpr = aws.String(expr)
	}
	if len(ctx.PositionalParams) > 0 && len(ctx.FixedParams) > 0 {
		return nil, errors.New("cannot mix positional params (?) with named params (:param)")
	}

	req := &dynamodb.QueryInput{
		TableName:              &ast.From,
		ProjectionExpression:   projectionExpr,
		KeyConditionExpression: aws.String(keyExpr),
	}
	if filterExpr != "" {
		req.FilterExpression = aws.String(filterExpr)
	}
	if len(ctx.Substitutions) > 0 {
		req.ExpressionAttributeNames = aws.StringMap(ctx.Substitutions)
	}
	if index != "" {
		req.IndexName = aws.String(index)
	}
	if ast.Limit != nil {
		req.Limit = aws.Int64(int64(*ast.Limit))
	}
	if ast.Descending != nil {
		req.ScanIndexForward = aws.Bool(!bool(*ast.Descending))
	}
	return &PreparedQuery{
		Query:            req,
		Columns:          ast.Projection.Columns,
		NamedParams:      visit.Context.NamedParams,
		PositionalParams: visit.Context.PositionalParams,
		FixedParams:      visit.Context.FixedParams,
	}, nil
}

// Context tracks expression state as DynamoDB request is built.
type Context struct {
	HashKey          string
	SortKey          string
	NamedParams      map[string]Empty
	PositionalParams map[int]string
	FixedParams      map[string]interface{}
	Substitutions    map[string]string

	positionalParamCount int
	genParamCount        int
	genSubCount          int
}

func NewContext(table *schema.Table, index string) *Context {
	var hashKey, sortKey string
	if index == "" {
		hashKey, sortKey = table.HashKey, table.SortKey
	} else {
		idx := table.GetIndex(index)
		hashKey, sortKey = idx.HashKey, idx.SortKey
	}
	return &Context{
		HashKey:          hashKey,
		SortKey:          sortKey,
		NamedParams:      make(map[string]Empty),
		PositionalParams: make(map[int]string),
		FixedParams:      make(map[string]interface{}),
		Substitutions:    make(map[string]string),
	}
}

// IsKey returns true if the field is a hash key or sort key for the selected table/index.
func (c *Context) IsKey(field string) bool {
	return c.HashKey == field || c.SortKey == field
}

// NextPositionalParam returns the next generated positional placeholder name. DynamoDB does not support positional
// parameters so we simulate it by generating named placeholders.
func (c *Context) NextPositionalParam() (int, string) {
	c.positionalParamCount++
	return c.positionalParamCount, fmt.Sprintf(":_pos%d", c.positionalParamCount)
}

// NextGeneratedParam returns the next generated parameter placeholder name. DynamoDB requires concrete values to
// be passed as expression values, so we simulate expression values by generating placeholders.
func (c *Context) NextGeneratedParam() string {
	c.genParamCount++
	return fmt.Sprintf(":_gen%d", c.genParamCount)
}

// BuildPath applies substitutions for reserved keywords in a parser.DocumentPath and marshals it into string format.
func (c *Context) BuildPath(path *parser.DocumentPath) string {
	buf := &bytes.Buffer{}
	for i, frag := range path.Fragment {
		if i != 0 {
			buf.WriteString(".")
		}
		buf.WriteString(c.substitute(frag.Symbol))
		for _, idx := range frag.Indexes {
			buf.WriteRune('[')
			buf.WriteString(strconv.Itoa(idx))
			buf.WriteRune(']')
		}
	}
	return buf.String()
}

func (c *Context) substitute(symbol string) string {
	if parser.IsReservedWord(symbol) {
		sub := "#" + symbol
		c.Substitutions[sub] = symbol
		return sub
	}
	if !validIdentifierRegexp.MatchString(symbol) {
		c.genSubCount++
		sub := fmt.Sprintf("#_gen%d", c.genSubCount)
		c.Substitutions[sub] = symbol
		return sub
	}
	return symbol
}

// matches an identifier that is valid for use in an expression
var validIdentifierRegexp = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

type keyAndFilter struct {
	Key    *parser.AndExpression
	Filter *parser.AndExpression
}

func extractKeyExpressions(expr *parser.AndExpression, isKey func(string) bool) *keyAndFilter {
	v := &keyAndFilter{
		Key:    &parser.AndExpression{},
		Filter: &parser.AndExpression{},
	}
	if expr == nil {
		return v
	}
	for _, term := range expr.And {
		if term.Operand != nil && isKey(term.Operand.Operand.String()) ||
			term.Function != nil && term.Function.FirstArgIsRef() && isKey(term.Function.Args[0].DocumentPath.String()) {
			v.Key.And = append(v.Key.And, term)
		} else {
			v.Filter.And = append(v.Filter.And, term)
		}
	}
	return v
}

type errHashKey string

func (hashKey errHashKey) Error() string {
	return fmt.Sprintf("partition key must appear exactly once in the WHERE clause, in an equality condition, such as: WHERE %s = :param", string(hashKey))
}

func buildKeyExpression(ctx *Context, key *parser.AndExpression) (string, error) {
	var hashExpr, sortExpr string
	visitor := &visitor{Context: ctx}
	for _, subExpr := range key.And {
		var expr, key string
		if subExpr.Function != nil {
			key = subExpr.Function.Args[0].DocumentPath.String()
			if ctx.HashKey == key {
				return "", errHashKey(ctx.HashKey)
			} else if subExpr.Function.Function != "begins_with" {
				return "", fmt.Errorf("sort key %q may not be used with function %s()", key, subExpr.Function.Function)
			}
			expr = visitor.VisitSimpleExpression(subExpr.Function)
		} else {
			key = subExpr.Operand.Operand.String()
			if key == ctx.HashKey {
				if subExpr.Operand.ConditionRHS.Compare == nil || subExpr.Operand.ConditionRHS.Compare.Operator != "=" {
					return "", errHashKey(ctx.HashKey)
				}
			}
			expr = visitor.VisitSimpleExpression(subExpr.Operand)
		}
		if ctx.HashKey == key {
			if hashExpr != "" {
				return "", fmt.Errorf("partition key %q can only appear once in WHERE clause", key)
			}
			hashExpr = expr
		} else if ctx.SortKey == key {
			if sortExpr != "" {
				return "", fmt.Errorf("sort key %q can only appear once in WHERE clause", key)
			}
			sortExpr = expr
		}
	}
	if hashExpr == "" {
		return "", errHashKey(ctx.HashKey)
	}
	if sortExpr != "" {
		return hashExpr + " AND " + sortExpr, nil
	}
	return hashExpr, nil
}

func buildFilterExpression(ctx *Context, filter *parser.AndExpression) (string, error) {
	v := &visitor{Context: ctx}
	return v.VisitFilterExpression(filter)
}

type Empty struct{}

type Acc struct {
	Key    string
	Filter string
}

type visitor struct {
	*Context
}

// VisitFilterExpression visits all nodes in the filter expression tree to build a filter expression.
func (v *visitor) VisitFilterExpression(n interface{}) (string, error) {
	if reflect.ValueOf(n).IsNil() {
		return "", nil
	}
	switch node := n.(type) {
	case *parser.ConditionExpression:
		filter := make([]string, 0, len(node.Or))
		for _, expr := range node.Or {
			acc, err := v.VisitFilterExpression(expr)
			if err != nil {
				return "", err
			}
			if acc != "" {
				filter = append(filter, acc)
			}
		}
		return strings.Join(filter, " OR "), nil
	case *parser.AndExpression:
		filter := make([]string, 0, len(node.And))
		for _, expr := range node.And {
			acc, err := v.VisitFilterExpression(expr)
			if err != nil {
				return "", err
			}
			if acc != "" {
				filter = append(filter, acc)
			}
		}
		return strings.Join(filter, " AND "), nil
	case *parser.Condition:
		switch {
		case node.Operand != nil:
			if v.Context.IsKey(node.Operand.Operand.String()) {
				return "", fmt.Errorf("partition key %q may not appear in nested expression", node.Operand.Operand.String())
			}
			return v.VisitSimpleExpression(node.Operand), nil
		case node.Function != nil:
			if node.Function.FirstArgIsRef() && v.Context.IsKey(node.Function.Args[0].DocumentPath.String()) {
				return "", fmt.Errorf("partition key %q may not appear in nested expression", node.Function.Args[0].DocumentPath)
			}
			return v.VisitSimpleExpression(node.Function), nil
		case node.Not != nil:
			return v.VisitFilterExpression(node.Not)
		case node.Parenthesized != nil:
			return v.VisitFilterExpression(node.Parenthesized)
		default:
			panic("invalid condition subtype")
		}
	case *parser.ParenthesizedExpression:
		acc, err := v.VisitFilterExpression(node.ConditionExpression)
		if err != nil {
			return "", err
		}
		return "(" + acc + ")", nil
	case *parser.NotCondition:
		acc, err := v.VisitFilterExpression(node.Condition)
		if err != nil {
			return "", err
		}
		return "NOT " + acc, nil
	default:
		panic("invalid type " + repr.String(node))
	}
}

// VisitSimpleExpression visits unary and binary expressions down to the leaf nodes.
func (v *visitor) VisitSimpleExpression(n interface{}) string {
	switch node := n.(type) {
	case *parser.ConditionOperand:
		return v.BuildPath(node.Operand) + " " + v.VisitSimpleExpression(node.ConditionRHS)
	case *parser.FunctionExpression:
		argStr := make([]string, len(node.Args))
		for i, arg := range node.Args {
			argStr[i] = v.VisitSimpleExpression(arg)
		}
		return fmt.Sprintf("%s(%s)", node.Function, strings.Join(argStr, ", "))
	case *parser.FunctionArgument:
		if node.DocumentPath != nil {
			return v.BuildPath(node.DocumentPath)
		}
		return v.VisitSimpleExpression(node.Value)
	case *parser.ConditionRHS:
		switch {
		case node.Compare != nil:
			return v.VisitSimpleExpression(node.Compare)
		case node.Between != nil:
			return v.VisitSimpleExpression(node.Between)
		case node.In != nil:
			return v.VisitSimpleExpression(node.In)
		default:
			panic("invalid rhs")
		}
	case *parser.Compare:
		return node.Operator + " " + v.VisitSimpleExpression(node.Operand)
	case *parser.Between:
		return fmt.Sprintf("BETWEEN %s AND %s",
			v.VisitSimpleExpression(node.Start), v.VisitSimpleExpression(node.End))
	case *parser.In:
		return fmt.Sprintf(" IN (%s)", v.VisitSimpleExpression(node.Values))
	case *parser.Operand:
		if node.SymbolRef != nil {
			return v.BuildPath(node.SymbolRef)
		}
		return v.VisitSimpleExpression(node.Value)
	case *parser.Value:
		switch {
		case node.PlaceHolder != nil:
			v.Context.NamedParams[*node.PlaceHolder] = Empty{}
			return *node.PlaceHolder
		case node.PositionalPlaceholder != nil:
			num, str := v.Context.NextPositionalParam()
			v.Context.PositionalParams[num] = str
			return str
		case node.Number != nil:
			name := v.Context.NextGeneratedParam()
			v.Context.FixedParams[name] = *node.Number
			return name
		case node.String != nil:
			name := v.Context.NextGeneratedParam()
			v.Context.FixedParams[name] = *node.String
			return name
		case node.Boolean != nil:
			name := v.Context.NextGeneratedParam()
			v.Context.FixedParams[name] = bool(*node.Boolean)
			return name
		case node.Null:
			name := v.Context.NextGeneratedParam()
			v.Context.FixedParams[name] = nil
			return name
		default:
			panic("invalid value" + repr.String(node))
		}
	default:
		panic("visitor does not recognize type " + reflect.TypeOf(node).String())
	}
}

func buildProjectionExpression(ctx *Context, expr *parser.ProjectionExpression) (string, error) {
	cols := make([]*parser.DocumentPath, 0, len(expr.Columns))
	for _, col := range expr.Columns {
		if col.DocumentPath != nil {
			cols = append(cols, col.DocumentPath)
		} else if col.Function != nil {
			fc, err := extractProjectionsFromFunction(col.Function)
			if err != nil {
				return "", err
			}
			cols = append(cols, fc...)
		} else {
			return "", fmt.Errorf("unexpected ProjectionColumn %v", *col)
		}
	}
	buf := &bytes.Buffer{}
	for i, col := range cols {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(ctx.BuildPath(col))
	}
	return buf.String(), nil
}

func extractProjectionsFromFunction(expr *parser.FunctionExpression) ([]*parser.DocumentPath, error) {
	if expr.Function != "document" {
		return nil, fmt.Errorf("function %q not allowed in projection", expr.Function)
	}
	cols := make([]*parser.DocumentPath, 0, len(expr.Args))
	for _, arg := range expr.Args {
		if arg.DocumentPath == nil {
			return nil, fmt.Errorf("args to document() must be document paths, got %s", arg.String())
		}
		cols = append(cols, arg.DocumentPath)
	}
	return cols, nil
}
