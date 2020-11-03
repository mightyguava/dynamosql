querybuilder.item{
  Query: "SELECT * FROM gamescores WHERE UserId = :UserId",
  Prepared: &querybuilder.PreparedQuery{
    Query: &dynamodb.QueryInput{
      _: struct {}{      },
      KeyConditionExpression: &"UserId = :UserId",
      TableName: &"gamescores",
    },
    NamedParams: querybuilder.NamedParams{
      ":UserId": querybuilder.Empty{      },
    },
    PositionalParams: map[int]string{    },
    FixedParams: map[string]interface {}{    },
  },
}