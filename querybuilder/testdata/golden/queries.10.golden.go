querybuilder.item{
  Query: "SELECT * FROM gamescores WHERE UserId = \"103\" DESC LIMIT 1",
  Prepared: &querybuilder.PreparedQuery{
    Query: &dynamodb.QueryInput{
      _: struct {}{      },
      KeyConditionExpression: &"UserId = :_gen1",
      Limit: &1,
      ScanIndexForward: &false,
      TableName: &"gamescores",
    },
    Limit: 1,
    NamedParams: querybuilder.NamedParams{    },
    PositionalParams: map[int]string{    },
    FixedParams: map[string]interface {}{
      ":_gen1": "103",
    },
  },
}