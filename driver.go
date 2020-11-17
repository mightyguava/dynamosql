package dynamosql

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/mightyguava/dynamosql/schema"
)

func init() {
	sql.Register("dynamodb", &Driver{})
}

// Driver is the DynamoDB SQL driver.
type Driver struct {
	cfg Config
}

var _ driver.Driver = &Driver{}
var _ driver.DriverContext = &Driver{}

// Config provides optional settings to the DynamoDB driver.
type Config struct {
	// If set, the driver will use this DynamoDB client. The Session param and the connection string will be ignored
	DynamoDB dynamodbiface.DynamoDBAPI
	// If set, and DynamoDB is not set, the driver will try to create a DynamoDB client using this session.
	Session *session.Session
	// If set, the wrapper collections []*dynamodb.AttributeValue, map[string]*dynamodb.AttributeValue will be mapped
	// unmarshaled into using the dynamodbattribute package into []interface{} and map[string]interface{}, respectively.
	AlwaysConvertCollectionsToGoType bool
}

// New creates a Driver instance using a custom config. This may be easier to use than via sql.Open.
// A driver created using this function can be used to obtain a sql.DB like so
//  driver, err := New(Config{Session: sess}).OpenConnector("")
//  db := sql.OpenDB(driver)
func New(cfg Config) *Driver {
	return &Driver{
		cfg: cfg,
	}
}

// NewDBWithClient returns a sql.DB backed by dynamosql using the given DynamoDB client.
func NewDBWithClient(ddb dynamodbiface.DynamoDBAPI) *sql.DB {
	return newDB(Config{DynamoDB: ddb})
}

// NewDBWithSession returns a sql.DB backed by dynamosql with a DynamoDB client constructed from the given Session.
func NewDBWithSession(sess *session.Session) *sql.DB {
	return newDB(Config{Session: sess})
}

func newDB(cfg Config) *sql.DB {
	conn, err := New(cfg).OpenConnector("")
	if err != nil {
		// OpenConnector does not return error if a Session or DynamoDB is provided.
		panic("unexpected error")
	}
	return sql.OpenDB(conn)
}

// Open a connection using the connection string.
func (d *Driver) Open(connStr string) (driver.Conn, error) {
	c, err := d.OpenConnector(connStr)
	if err != nil {
		return nil, err
	}
	return c.Connect(context.Background())
}

// OpenConnector initializes and returns a Connector. The db/sql package will call this exactly once
// per sql.Open() call. New connections to the database will use the returned Connector.
func (d *Driver) OpenConnector(connStr string) (driver.Connector, error) {
	var dynamo dynamodbiface.DynamoDBAPI
	if d.cfg.DynamoDB != nil {
		dynamo = d.cfg.DynamoDB
	} else {
		var err error
		sess := d.cfg.Session
		if sess == nil {
			sess, err = session.NewSession(nil)
			if err != nil {
				return nil, err
			}
		}
		dynamo = dynamodb.New(sess)
	}
	return &connector{
		dynamo:      dynamo,
		driver:      d,
		tables:      schema.NewTableLoader(dynamo),
		mapToGoType: d.cfg.AlwaysConvertCollectionsToGoType,
	}, nil
}

type connector struct {
	driver      *Driver
	dynamo      dynamodbiface.DynamoDBAPI
	tables      *schema.TableLoader
	mapToGoType bool
}

var _ driver.Connector = &connector{}

func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	return &conn{dynamo: c.dynamo, tables: c.tables, mapToGoType: c.mapToGoType}, nil
}

func (c *connector) Driver() driver.Driver {
	return c.driver
}
