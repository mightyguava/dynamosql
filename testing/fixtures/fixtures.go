package fixtures

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/require"
)

// SetUp creates and tears down a dynamodb table fixture.
func SetUp(t *testing.T, fixtures ...Fixture) *session.Session {
	t.Helper()

	cfg := aws.NewConfig().
		WithEndpoint("http://localhost:8000").
		WithRegion("us-west-2").
		WithCredentials(credentials.NewStaticCredentials("fake", "secret", ""))
	sess := session.Must(session.NewSession(cfg))
	client := dynamodb.New(sess)

	// Try to delete the tables before and after tests
	cleanup := func() {
		for _, fixture := range fixtures {
			_, err := client.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(fixture.Table)})
			if err != nil {
				if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == dynamodb.ErrCodeResourceNotFoundException {
					continue
				}
				t.Errorf("unexpected error deleting table %q: %s", fixture.Table, err)
			}
		}
	}
	cleanup()

	for _, fixture := range fixtures {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		_, err := client.CreateTable(fixture.Create)
		require.NoError(t, err, "could not create table %q", fixture.Table)

		for {
			_, err = client.DescribeTable(&dynamodb.DescribeTableInput{TableName: &fixture.Table})
			if err == nil {
				break
			}
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == dynamodb.ErrCodeResourceNotFoundException {
				// continue
			} else {
				t.Fatalf("error waiting for table creation %q: %s", fixture.Table, err)
			}
			select {
			case <-ctx.Done():
				t.Fatalf("timed out waiting for table creation %q", fixture.Table)
			case <-time.After(100 * time.Millisecond):
				// continue
			}
		}

		if fixture.Data != nil {
			fixture.Data(t, client)
		}
	}
	return sess
}

// Fixture defines a set of DynamoDB tables and data for testing.
type Fixture struct {
	Table  string
	Create *dynamodb.CreateTableInput
	Data   func(t *testing.T, client *dynamodb.DynamoDB)
}
