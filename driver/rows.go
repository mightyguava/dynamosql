package driver

import (
	"database/sql/driver"
	"errors"
	"io"
	"sort"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/mightyguava/dynamosql/parser"
)

type Rows struct {
	ast  parser.Select
	req  *dynamodb.QueryInput
	resp *dynamodb.QueryOutput

	mu      sync.Mutex
	cols    []string
	nextRow int
}

var _ driver.Rows = &Rows{}

func (r *Rows) Columns() []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cols != nil {
		return r.cols
	}

	if len(r.resp.Items) == 0 {
		return []string{}
	}

	cols := make([]string, 0, len(r.resp.Items[0]))
	for k := range r.resp.Items[0] {
		cols = append(cols, k)
	}
	sort.Strings(cols)

	r.cols = cols
	return cols
}

func (r *Rows) Close() error {
	return nil
}

func (r *Rows) Next(dest []driver.Value) error {
	cols := r.Columns()
	if r.nextRow >= len(r.resp.Items) {
		return io.EOF
	}
	row := r.resp.Items[r.nextRow]
	r.nextRow++
	for i, key := range cols {
		v, err := convertValue(row[key])
		if err != nil {
			return err
		}
		dest[i] = v
	}
	return nil
}

func convertValue(av *dynamodb.AttributeValue) (interface{}, error) {
	switch {
	case av.B != nil:
		return av.B, nil
	case av.BOOL != nil:
		return *av.BOOL, nil
	case av.BS != nil:
		return nil, errors.New("unsupported type Binary Set")
	case av.L != nil:
		return nil, errors.New("unsupported type List")
	case av.M != nil:
		return nil, errors.New("unsupported type Map")
	case av.N != nil:
		return strconv.ParseFloat(*av.N, 64)
	case av.NS != nil:
		return nil, errors.New("unsupported type Number Set")
	case av.S != nil:
		return *av.S, nil
	case av.SS != nil:
		return nil, errors.New("unsupported type String Set")
	default:
		return nil, nil
	}
}
