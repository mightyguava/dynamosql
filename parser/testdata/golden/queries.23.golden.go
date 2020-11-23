parser.row{
  Query: "CREATE TABLE movies (title STRING, year NUMBER);",
  AST: &parser.AST{
    CreateTable: &parser.CreateTable{
      Table: "movies",
      Entries: []*parser.CreateTableEntry{
        {
          Attr: &parser.TableAttr{
            Name: "title",
            Type: "STRING",
          },
        },
        {
          Attr: &parser.TableAttr{
            Name: "year",
            Type: "NUMBER",
          },
        },
      },
    },
  },
}