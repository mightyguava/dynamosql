querybuilder.item{
  Query: "SELECT * FROM gamescores WHERE UserId = \"103\" AND begins_with(GameTitle, \"Galaxy\")",
  Prepared: &querybuilder.PreparedQuery{
    Query: &dynamodb.QueryInput{
      _: struct {}{      },
      KeyConditionExpression: &"UserId = :_gen1 AND begins_with(GameTitle, :_gen2)",
      TableName: &"gamescores",
    },
    NamedParams: querybuilder.NamedParams{    },
    PositionalParams: map[int]string{    },
    FixedParams: map[string]interface {}{
      ":_gen1": "103",
      ":_gen2": "Galaxy",
    },
  },
}