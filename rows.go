package dynamosql

import (
	"database/sql/driver"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/mightyguava/dynamosql/parser"
)

type rows struct {
	resp        *dynamodb.QueryOutput
	nextPage    func(lastEvaluatedKey map[string]*dynamodb.AttributeValue) (*dynamodb.QueryOutput, error)
	cols        []*parser.ProjectionColumn
	mapToGoType bool

	nextRow int
}

var _ driver.Rows = &rows{}

// Returns the number of columns.
// Caveat: the number of columns is always equal to the number of attributes in the first returned item.
func (r *rows) Columns() []string {
	if len(r.cols) == 0 {
		return []string{"document"}
	}

	cols := make([]string, 0, len(r.cols))
	for _, col := range r.cols {
		if col.Function != nil {
			cols = append(cols, "document")
		} else {
			cols = append(cols, col.String())
		}
	}

	return cols
}

func (r *rows) Close() error {
	return nil
}

func (r *rows) Next(dest []driver.Value) error {
	if r.nextRow >= len(r.resp.Items) {
		resp, err := r.nextPage(r.resp.LastEvaluatedKey)
		if err != nil {
			return err
		}
		r.nextRow = 0
		r.resp = resp
	}
	row := r.resp.Items[r.nextRow]
	r.nextRow++

	// SELECT *
	if len(r.cols) == 0 {
		dest[0] = r.remap(row)
	}

	for i, col := range r.cols {
		if col.Function != nil {
			// SELECT document(...) returns the whole document
			dest[i] = r.remap(row)
		} else {
			dest[i] = r.remap(pluck(&dynamodb.AttributeValue{M: row}, col.DocumentPath))
		}
	}
	return nil
}

func (r *rows) remap(data interface{}) interface{} {
	if !r.mapToGoType {
		return data
	}
	switch data := data.(type) {
	case []*dynamodb.AttributeValue:
		out := make([]interface{}, len(data))
		err := dynamodbattribute.UnmarshalList(data, &out)
		if err != nil {
			panic(fmt.Sprintf("unexpected conversion error: %+v", err))
		}
		return out
	case map[string]*dynamodb.AttributeValue:
		out := make(map[string]interface{}, len(data))
		err := dynamodbattribute.UnmarshalMap(data, &out)
		if err != nil {
			panic(fmt.Sprintf("unexpected conversion error: %+v", err))
		}
		return out
	default:
		return data
	}
}

func pluck(pos *dynamodb.AttributeValue, path *parser.DocumentPath) driver.Value {
	var ok bool
	for _, frag := range path.Fragment {
		if pos.M == nil {
			return nil
		}
		pos, ok = pos.M[frag.Symbol]
		if !ok {
			return nil
		}
		for _, idx := range frag.Indexes {
			switch {
			case pos.L != nil:
				if idx >= len(pos.L) {
					return nil
				}
				pos = pos.L[idx]
			case pos.BS != nil:
				if idx >= len(pos.BS) {
					return nil
				}
				pos = &dynamodb.AttributeValue{B: pos.BS[idx]}
			case pos.NS != nil:
				if idx >= len(pos.NS) {
					return nil
				}
				pos = &dynamodb.AttributeValue{N: pos.NS[idx]}
			case pos.SS != nil:
				if idx >= len(pos.SS) {
					return nil
				}
				pos = &dynamodb.AttributeValue{S: pos.SS[idx]}
			default:
				return nil
			}
		}
	}
	return convertValue(pos)
}

func convertValue(av *dynamodb.AttributeValue) interface{} {
	switch {
	// bool
	case av.BOOL != nil:
		return *av.BOOL
	// number (returned as string since db/sql supports string conversion)
	case av.N != nil:
		return *av.N
	// string
	case av.S != nil:
		return *av.S
	// byte array
	case av.B != nil:
		return av.B
	// list
	case av.L != nil:
		return av.L
	// map
	case av.M != nil:
		return av.M
	// set of numbers
	case av.NS != nil:
		return aws.StringValueSlice(av.NS)
	// set of strings
	case av.SS != nil:
		return aws.StringValueSlice(av.SS)
	// set of byte arrays
	case av.BS != nil:
		return av.BS
	// null
	default:
		return nil
	}
}
