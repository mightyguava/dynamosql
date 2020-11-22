package querybuilder

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/alecthomas/repr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/mightyguava/dynamosql/parser"
	"github.com/mightyguava/dynamosql/schema"
)

type PreparedInsert struct {
	Table       *schema.Table
	Placeholder string
	Values      []map[string]*dynamodb.AttributeValue
	Returning   *string
	Replace     bool
}

func PrepareInsert(ctx context.Context, tables *schema.TableLoader, ast *parser.AST) (*PreparedInsert, error) {
	var ins *parser.Insert
	replace := false
	switch {
	case ast.Insert != nil:
		ins = ast.Insert
		if ins.Returning != nil && *ins.Returning != "NONE" {
			return nil, errors.New("RETURNING is not allowed on INSERT")
		}
	case ast.Replace != nil:
		ins = ast.Replace
		replace = true
	default:
		return nil, fmt.Errorf("expected INSERT but got %s", repr.String(ast))
	}
	table, err := tables.Get(ctx, ins.Into)
	if err != nil {
		return nil, err
	}
	var literals []string
	var usePlaceholder bool
	var placeholder string
	for _, v := range ins.Values {
		switch {
		case v.PositionalPlaceholder:
			usePlaceholder = true
		case v.PlaceHolder != nil:
			usePlaceholder = true
			placeholder = *v.PlaceHolder
		case v.Object != nil:
			literals = append(literals, v.Object.String())
		case v.Str != nil:
			literals = append(literals, *v.Str)
		default:
			return nil, fmt.Errorf("VALUES expression may must be a placeholder or string, but was %s", repr.String(v))
		}
	}
	if usePlaceholder && len(ins.Values) > 1 {
		return nil, errors.New("when using placeholder parameters, INSERT may contain exactly one placeholder")
	}

	var values []map[string]*dynamodb.AttributeValue
	if len(literals) > 0 {
		values = make([]map[string]*dynamodb.AttributeValue, 0, len(literals))
		for _, l := range literals {
			av, err := jsonStringToDynamodbMap(l)
			if err != nil {
				return nil, err
			}
			values = append(values, av)
		}
	}

	return &PreparedInsert{
		Table:       table,
		Placeholder: placeholder,
		Values:      values,
		Returning:   ins.Returning,
		Replace:     replace,
	}, nil
}

func (p *PreparedInsert) Do(ctx context.Context, dynamo dynamodbiface.DynamoDBAPI, args []driver.NamedValue) (*DriverResult, error) {
	if len(args) > 0 && len(p.Values) > 0 {
		return nil, errors.New("no arguments expected")
	}
	var values []map[string]*dynamodb.AttributeValue
	if len(p.Values) > 0 {
		values = p.Values
	} else {
		if len(args) > 1 && len(p.Values) == 0 {
			return nil, errors.New("too many arguments")
		}
		arg := args[0]
		if arg.Name != p.Placeholder {
			return nil, fmt.Errorf("unexpected named argument %q", arg.Name)
		}
		var err error
		values, err = argToListOfMaps(arg.Value)
		if err != nil {
			return nil, err
		}
	}

	if len(values) == 0 {
		return nil, errors.New("no values to insert")
	}
	if len(values) == 1 {
		resp, err := dynamo.PutItem(p.toPutItem(values[0]))
		if err != nil {
			return nil, err
		}
		return &DriverResult{
			count:    1,
			returned: resp.Attributes,
		}, nil
	}
	if p.Returning != nil && *p.Returning != "NONE" {
		return nil, errors.New("cannot use RETURNING with more than 1 item")
	}
	_, err := dynamo.TransactWriteItems(p.toTransactWrite(values))
	if err != nil {
		return nil, err
	}
	return &DriverResult{count: len(values)}, nil
}

func (p *PreparedInsert) toTransactWrite(items []map[string]*dynamodb.AttributeValue) *dynamodb.TransactWriteItemsInput {
	conditionExpr, exprAttrNames := p.conditionExpr()

	puts := make([]*dynamodb.TransactWriteItem, 0, len(items))
	for _, item := range items {
		puts = append(puts, &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				ConditionExpression:      conditionExpr,
				ExpressionAttributeNames: exprAttrNames,
				Item:                     item,
				TableName:                &p.Table.Name,
			},
		})
	}
	return &dynamodb.TransactWriteItemsInput{
		TransactItems: puts,
	}
}

func (p *PreparedInsert) toPutItem(item map[string]*dynamodb.AttributeValue) *dynamodb.PutItemInput {
	conditionExpr, exprAttrNames := p.conditionExpr()
	return &dynamodb.PutItemInput{
		ConditionExpression:      conditionExpr,
		ExpressionAttributeNames: exprAttrNames,
		Item:                     item,
		TableName:                &p.Table.Name,
		ReturnValues:             p.Returning,
	}
}

// In REPLACE, return nil. In INSERT, return a condition that will fail Puts on existing rows.
func (p *PreparedInsert) conditionExpr() (conditionExpr *string, exprAttrNames map[string]*string) {
	if p.Replace {
		return nil, nil
	}
	ctx := Context{}
	key := ctx.substitute(p.Table.HashKey)
	conditionExpr = aws.String(fmt.Sprintf("attribute_not_exists(%s)", key))
	exprAttrNames = ctx.ExpressionAttributeNames()
	return
}

func jsonStringToDynamodbMap(v string) (map[string]*dynamodb.AttributeValue, error) {
	var asMap map[string]interface{}
	if err := json.Unmarshal([]byte(v), &asMap); err != nil {
		return nil, err
	}
	return dynamodbattribute.MarshalMap(asMap)
}

func argToListOfMaps(v interface{}) ([]map[string]*dynamodb.AttributeValue, error) {
	t := reflect.ValueOf(v)
	if t.Kind() != reflect.Slice {
		av, err := marshalDocument(v)
		if err != nil {
			return nil, err
		}
		return []map[string]*dynamodb.AttributeValue{av}, nil
	}
	m := make([]map[string]*dynamodb.AttributeValue, t.Len())
	for i := 0; i < t.Len(); i++ {
		av, err := marshalDocument(t.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		m[i] = av
	}
	return m, nil
}

func marshalDocument(v interface{}) (map[string]*dynamodb.AttributeValue, error) {
	t := reflect.ValueOf(v)
	if t.Kind() == reflect.Ptr {
		v = t.Elem().Interface()
		t = reflect.ValueOf(v)
	}
	switch v := v.(type) {
	case string:
		// If we get a string, it must be a JSON string
		return jsonStringToDynamodbMap(v)
	default:
		// Otherwise, use dynamodbattribute to marshal, and expect a map
		av, err := dynamodbattribute.Marshal(v)
		if err != nil {
			return nil, err
		}
		if av.M == nil {
			return nil, errors.New("failed to marshal value into a map")
		}
		return av.M, nil
	}
}
