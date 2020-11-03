querybuilder.item{
  Query: "SELECT * FROM gamescores USE INDEX (GameTitleIndex) WHERE GameTitle = :title AND UserId > \"45\"",
  Prepared: &querybuilder.PreparedQuery{
    Query: &dynamodb.QueryInput{
      _: struct {}{      },
      FilterExpression: &"UserId > :_gen1",
      IndexName: &"GameTitleIndex",
      KeyConditionExpression: &"GameTitle = :title",
      TableName: &"gamescores",
    },
    NamedParams: querybuilder.NamedParams{
      ":title": querybuilder.Empty{      },
    },
    PositionalParams: map[int]string{    },
    FixedParams: map[string]interface {}{
      ":_gen1": "45",
    },
  },
}