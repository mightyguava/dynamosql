parser.row{
  Query: "SELECT title, year FROM movies",
  AST: parser.Select{
    Projection: &parser.ProjectionExpression{
      Columns: []*parser.ProjectionColumn{
        {
          DocumentPath: &parser.DocumentPath{
            Fragment: []parser.PathFragment{
              {
                Symbol: "title",
              },
            },
          },
        },
        {
          DocumentPath: &parser.DocumentPath{
            Fragment: []parser.PathFragment{
              {
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