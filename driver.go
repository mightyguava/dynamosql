package dynamosql

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

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
	Session *session.Session
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
	var err error
	sess := d.cfg.Session
	if sess == nil {
		sess, err = session.NewSession(nil)
		if err != nil {
			return nil, err
		}
	}
	dynamo := dynamodb.New(sess)
	return &connector{
		dynamo: dynamo,
		driver: d,
		tables: schema.NewTableLoader(dynamo),
	}, nil
}

type connector struct {
	driver *Driver
	dynamo *dynamodb.DynamoDB
	tables *schema.TableLoader
}

var _ driver.Connector = &connector{}

func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	return &conn{dynamo: c.dynamo, tables: c.tables}, nil
}

func (c *connector) Driver() driver.Driver {
	return c.driver
}
