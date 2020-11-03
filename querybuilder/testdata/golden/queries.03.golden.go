querybuilder.item{
  Query: "SELECT * FROM gamescores WHERE UserId = \"103\" AND GameTitle BETWEEN \"Galaxy\" AND \"Meteor\" AND TopScore > 1000",
  Prepared: &querybuilder.PreparedQuery{
    Query: &dynamodb.QueryInput{
      _: struct {}{      },
      FilterExpression: &"TopScore > :_gen4",
      KeyConditionExpression: &"UserId = :_gen1 AND GameTitle BETWEEN :_gen2 AND :_gen3",
      TableName: &"gamescores",
    },
    NamedParams: querybuilder.NamedParams{    },
    PositionalParams: map[int]string{    },
    FixedParams: map[string]interface {}{
      ":_gen1": "103",
      ":_gen2": "Galaxy",
      ":_gen3": "Meteor",
      ":_gen4": 1000,
    },
  },
}