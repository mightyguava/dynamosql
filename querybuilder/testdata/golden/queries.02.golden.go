querybuilder.item{
  Query: "SELECT * FROM gamescores WHERE UserId = :UserId AND GameTitle BETWEEN :MinGameTitle AND :MaxGameTitle AND TopScore > :MinTopScore",
  Prepared: &querybuilder.PreparedQuery{
    Query: &dynamodb.QueryInput{
      _: struct {}{      },
      FilterExpression: &"TopScore > :MinTopScore",
      KeyConditionExpression: &"UserId = :UserId AND GameTitle BETWEEN :MinGameTitle AND :MaxGameTitle",
      TableName: &"gamescores",
    },
    NamedParams: querybuilder.NamedParams{
      ":MaxGameTitle": querybuilder.Empty{      },
      ":MinGameTitle": querybuilder.Empty{      },
      ":MinTopScore": querybuilder.Empty{      },
      ":UserId": querybuilder.Empty{      },
    },
    PositionalParams: map[int]string{    },
    FixedParams: map[string]interface {}{    },
  },
}