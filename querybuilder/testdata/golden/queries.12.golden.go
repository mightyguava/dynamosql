querybuilder.item{
  Query: "SELECT * FROM gamescores WHERE TopScore > ? AND UserId = ?",
  Prepared: &querybuilder.PreparedQuery{
    Query: &dynamodb.QueryInput{
      _: struct {}{      },
      FilterExpression: &"TopScore > :_pos1",
      KeyConditionExpression: &"UserId = :_pos2",
      TableName: &"gamescores",
    },
    NamedParams: querybuilder.NamedParams{    },
    PositionalParams: map[int]string{
      1: ":_pos1",
      2: ":_pos2",
    },
    FixedParams: map[string]interface {}{    },
  },
}