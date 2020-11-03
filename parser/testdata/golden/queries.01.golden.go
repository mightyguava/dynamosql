parser.row{
  Query: "SELECT title, year FROM movies",
  AST: parser.Select{
    Projection: &parser.ProjectionExpression{
      Columns: []*parser.ProjectionColumn{
        &parser.ProjectionColumn{
          DocumentPath: &parser.DocumentPath{
            Fragment: []parser.PathFragment{
              parser.PathFragment{
                Symbol: "title",
              },
            },
          },
        },
        &parser.ProjectionColumn{
          DocumentPath: &parser.DocumentPath{
            Fragment: []parser.PathFragment{
              parser.PathFragment{
                Symbol: "year",
              },
            },
          },
        },
      },
    },
    From: "movies",
  },
}