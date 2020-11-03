parser.row{
  Query: "SELECT * FROM movies",
  AST: parser.Select{
    Projection: &parser.ProjectionExpression{
      All: true,
    },
    From: "movies",
  },
}