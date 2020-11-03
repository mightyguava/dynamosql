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
      &parser.ProjectionColumn{
        DocumentPath: &parser.DocumentPath{
          Fragment: []parser.PathFragment{
            parser.PathFragment{
              Symbol: "foo.bar",
            },
          },
        },
      },
      &parser.ProjectionColumn{
        DocumentPath: &parser.DocumentPath{
          Fragment: []parser.PathFragment{
            parser.PathFragment{
              Symbol: "foo",
            },
            parser.PathFragment{
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