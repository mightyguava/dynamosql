package fixtures

import (
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/require"
)

// GameScores is a fixture with a video game high score table
var GameScores = Fixture{
	Table: *gameScoresTable,
	Create: &dynamodb.CreateTableInput{
		TableName: gameScoresTable,
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("UserId"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("GameTitle"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Wins"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("TopScore"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("UserId"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("GameTitle"),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		LocalSecondaryIndexes: []*dynamodb.LocalSecondaryIndex{
			{
				IndexName: aws.String("UserWinsIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						KeyType:       aws.String(dynamodb.KeyTypeHash),
						AttributeName: aws.String("UserId"),
					},
					{
						KeyType:       aws.String(dynamodb.KeyTypeRange),
						AttributeName: aws.String("Wins"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String(dynamodb.ProjectionTypeKeysOnly),
				},
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("GameTitleIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						KeyType:       aws.String(dynamodb.KeyTypeHash),
						AttributeName: aws.String("GameTitle"),
					},
					{
						KeyType:       aws.String(dynamodb.KeyTypeRange),
						AttributeName: aws.String("TopScore"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String(dynamodb.ProjectionTypeKeysOnly),
				},
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
	},
	Data: func(t *testing.T, client *dynamodb.DynamoDB) {
		for _, s := range gameScoresData {
			_, err := client.PutItem(&dynamodb.PutItemInput{
				TableName: gameScoresTable,
				Item: map[string]*dynamodb.AttributeValue{
					"UserId":    {S: &s.UserID},
					"GameTitle": {S: &s.GameTitle},
					"TopScore":  {N: aws.String(strconv.Itoa(s.TopScore))},
					"Wins":      {N: aws.String(strconv.Itoa(s.Wins))},
					"Losses":    {N: aws.String(strconv.Itoa(s.Losses))},
				},
			})
			require.NoError(t, err)
		}
	},
}

// gameScoresTable is the name of the table that contains the video game high scores table fixture.
var gameScoresTable = aws.String("gamescores")

type GameScore struct {
	UserID    string
	GameTitle string
	TopScore  int
	Wins      int
	Losses    int
}

var gameScoresData = []GameScore{
	{"101", "Galaxy Invaders", 5842, 21, 72},
	{"101", "Meteor Blasters", 1000, 12, 3},
	{"101", "Starship X", 24, 4, 9},
	{"102", "Alien Adventure", 192, 32, 192},
	{"102", "Galaxy Invaders", 0, 0, 5},
	{"103", "Attack Ships", 3, 1, 8},
	{"103", "Galaxy Invaders", 2317, 40, 3},
	{"103", "Meteor Blasters", 723, 22, 12},
	{"103", "Starship X", 42, 4, 19},
}
