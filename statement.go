package dynamosql

import (
	"context"
	"database/sql/driver"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/mightyguava/dynamosql/querybuilder"
)

type execStmt struct {
	legacyStmtMixin
	preparedStmt querybuilder.ExecStmt
	dynamo       dynamodbiface.DynamoDBAPI
	mapToGoType  bool
}

func (s *execStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	return s.preparedStmt.Do(ctx, s.dynamo, args)
}

func (s *execStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	result, err := s.preparedStmt.Do(ctx, s.dynamo, args)
	if err != nil {
		return nil, err
	}
	return &oneRow{item: result.Item()}, nil
}

type queryStmt struct {
	legacyStmtMixin
	preparedStmt *querybuilder.PreparedQuery
	dynamo       dynamodbiface.DynamoDBAPI
	mapToGoType  bool
}

func (s *queryStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	return nil, errors.New("called Exec() called SELECT")
}

func (s *queryStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	q := s.preparedStmt
	req, err := q.NewRequest(args)
	if err != nil {
		return nil, err
	}
	resp, err := s.dynamo.QueryWithContext(ctx, req)
	if err != nil {
		return nil, err
	}
	return &rows{
		nextPage: func(lastEvaluatedKey map[string]*dynamodb.AttributeValue) (*dynamodb.QueryOutput, error) {
			for lastEvaluatedKey != nil {
				req.ExclusiveStartKey = lastEvaluatedKey
				// nolint: govet
				resp, err := s.dynamo.QueryWithContext(ctx, req)
				if err != nil {
					return nil, err
				}
				if len(resp.Items) > 0 {
					return resp, nil
				}
				// An empty response does not necessarily indicate there are no more results. It's possible the
				// filter expression filtered out all values in this range. Need to keep paging until LastEvaluatedKey
				// is nil.
				if resp.LastEvaluatedKey != nil {
					lastEvaluatedKey = resp.LastEvaluatedKey
				}
			}
			return nil, io.EOF
		},
		cols:        q.Columns,
		resp:        resp,
		mapToGoType: s.mapToGoType,
		limit:       q.Limit,
	}, nil
}

// wrapper type just for compile time type checking.
type fullStmt interface {
	driver.Stmt
	driver.StmtExecContext
	driver.StmtQueryContext
}

var (
	_ fullStmt = &execStmt{}
	_ fullStmt = &queryStmt{}
)

// mixin to provide no-op/panic implementations of useless db/sql methods
type legacyStmtMixin struct{}

func (d legacyStmtMixin) Close() error  { return nil }
func (d legacyStmtMixin) NumInput() int { return -1 }
func (d legacyStmtMixin) Exec(args []driver.Value) (driver.Result, error) {
	panic("unexpected call to Exec")
}
func (d legacyStmtMixin) Query(args []driver.Value) (driver.Rows, error) {
	panic("unexpected call to Query")
}
