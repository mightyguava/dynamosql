package driver

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/alecthomas/repr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/require"
)

func setUp(t *testing.T) *session.Session {
	t.Helper()

	cfg := aws.NewConfig().
		WithEndpoint("http://localhost:8000").
		WithRegion("us-west-2").
		WithCredentials(credentials.NewStaticCredentials("fake", "secret", ""))
	sess := session.Must(session.NewSession(cfg))
	client := dynamodb.New(sess)

	// Try to delete the table before and after tests
	cleanup := func() {
		_, _ = client.DeleteTable(&dynamodb.DeleteTableInput{TableName: moviesTable})
	}
	cleanup()
	t.Cleanup(cleanup)

	_, err := client.CreateTable(&dynamodb.CreateTableInput{
		TableName: moviesTable,
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("title"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("year"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("title"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("year"),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
	})
	require.NoError(t, err)

	for _, m := range topMovies {
		_, err := client.PutItem(&dynamodb.PutItemInput{
			TableName: moviesTable,
			Item: map[string]*dynamodb.AttributeValue{
				"title": {S: aws.String(m.title)},
				"year":  {N: aws.String(strconv.Itoa(m.year))},
			},
		})
		require.NoError(t, err)
	}
	return sess
}

var moviesTable = aws.String("movies")

type movie struct {
	title string
	year  int
}

var topMovies = []movie{
	{"The Shawshank Redemption", 1994},
	{"The Godfather", 1972},
	{"The Godfather: Part II", 1974},
	{"The Dark Knight", 2008},
	{"12 Angry Men", 1957},
}

func TestDriver(t *testing.T) {
	sess := setUp(t)

	driver, err := New(Config{Session: sess}).OpenConnector("")
	require.NoError(t, err)
	db := sql.OpenDB(driver)
	err = db.Ping()
	require.NoError(t, err)

	rows, err := db.Query("SELECT * FROM movies")
	require.NoError(t, err)

	var movies []movie
	for rows.Next() {
		var m movie
		err = rows.Scan(&m.title, &m.year)
		require.NoError(t, err)
		movies = append(movies, m)
	}
	repr.Println(movies)
}
