parser.row{
  Query: "SELECT * FROM movies",
  AST: &parser.AST{
    Select: &parser.Select{
      Projection: &parser.ProjectionExpression{
        All: true,
      },
      From: "movies",
    },
  },
}