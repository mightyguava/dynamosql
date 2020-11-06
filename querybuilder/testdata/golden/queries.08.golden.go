querybuilder.item{
  Query: "SELECT `foo.bar`, `foo`.`bar` FROM movies WHERE title = :title",
  Prepared: &querybuilder.PreparedQuery{
    Query: &dynamodb.QueryInput{
      _: struct {}{      },
      ExpressionAttributeNames: map[string]*string{
        "#_gen1": &"foo.bar",
      },
      KeyConditionExpression: &"title = :title",
      ProjectionExpression: &"#_gen1, foo.bar",
      TableName: &"movies",
    },
    Columns: []*parser.ProjectionColumn{
      {
        DocumentPath: &parser.DocumentPath{
          Fragment: []*parser.PathFragment{
            {
              Symbol: "foo.bar",
            },
          },
        },
      },
      {
        DocumentPath: &parser.DocumentPath{
          Fragment: []*parser.PathFragment{
            {
              Symbol: "foo",
            },
            {
              Symbol: "bar",
            },
          },
        },
      },
    },
    NamedParams: querybuilder.NamedParams{
      ":title": querybuilder.Empty{      },
    },
    PositionalParams: map[int]string{    },
    FixedParams: map[string]interface {}{    },
  },
}