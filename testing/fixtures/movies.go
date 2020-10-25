package fixtures

import (
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/require"
)

// Movies is a fixture with a movies table and a few top movies.
var Movies = Fixture{
	Table: *moviesTable,
	Create: &dynamodb.CreateTableInput{
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
	},
	Data: func(t *testing.T, client *dynamodb.DynamoDB) {
		for _, m := range TopMovies {
			_, err := client.PutItem(&dynamodb.PutItemInput{
				TableName: moviesTable,
				Item: map[string]*dynamodb.AttributeValue{
					"title": {S: aws.String(m.Title)},
					"year":  {N: aws.String(strconv.Itoa(m.Year))},
				},
			})
			require.NoError(t, err)
		}
	},
}

// moviesTable is the name of the movies table in the Movies
var moviesTable = aws.String("movies")

// Movie is a container for movie data
type Movie struct {
	Title string
	Year  int
}

// TopMovies is the list of top movies that are added to the Movies fixture.
var TopMovies = []Movie{
	{"The Shawshank Redemption", 1994},
	{"The Godfather", 1972},
	{"The Godfather: Part II", 1974},
	{"The Dark Knight", 2008},
	{"12 Angry Men", 1957},
}
