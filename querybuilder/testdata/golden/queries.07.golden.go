querybuilder.item{
  Query: "SELECT title, year FROM movies WHERE title = :title AND year > 2009 AND escaped = TRUE",
  Prepared: &querybuilder.PreparedQuery{
    Query: &dynamodb.QueryInput{
      _: struct {}{      },
      ExpressionAttributeNames: map[string]*string{
        "#escaped": &"escaped",
        "#year": &"year",
      },
      FilterExpression: &"#escaped = :_gen2",
      KeyConditionExpression: &"title = :title AND #year > :_gen1",
      ProjectionExpression: &"title, #year",
      TableName: &"movies",
    },
    Columns: []*parser.ProjectionColumn{
      {
        DocumentPath: &parser.DocumentPath{
          Fragment: []*parser.PathFragment{
            {
              Symbol: "title",
            },
          },
        },
      },
      {
        DocumentPath: &parser.DocumentPath{
          Fragment: []*parser.PathFragment{
            {
              Symbol: "year",
            },
          },
        },
      },
    },
    NamedParams: querybuilder.NamedParams{
      ":title": querybuilder.Empty{      },
    },
    PositionalParams: map[int]string{    },
    FixedParams: map[string]interface {}{
      ":_gen1": 2009,
      ":_gen2": true,
    },
  },
}