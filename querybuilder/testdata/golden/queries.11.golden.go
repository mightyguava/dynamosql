querybuilder.item{
  Query: "SELECT * FROM gamescores WHERE UserId = ? AND TopScore > ?",
  Prepared: &querybuilder.PreparedQuery{
    Query: &dynamodb.QueryInput{
      _: struct {}{      },
      FilterExpression: &"TopScore > :_pos2",
      KeyConditionExpression: &"UserId = :_pos1",
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