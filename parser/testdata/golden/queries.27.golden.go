parser.row{
  Query: "DROP TABLE movies;",
  AST: &parser.AST{
    DropTable: &parser.DropTable{
      Table: "movies",
    },
  },
}