parser.row{
  Query: "CREATE TABLE movies (title STRING HASH KEY, year NUMBER RANGE KEY);",
  AST: &parser.AST{
    CreateTable: &parser.CreateTable{
      Table: "movies",
      Entries: []*parser.CreateTableEntry{
        {
          Attr: &parser.TableAttr{
            Name: "title",
            Type: "STRING",
            Key: "HASH",
          },
        },
        {
          Attr: &parser.TableAttr{
            Name: "year",
            Type: "NUMBER",
            Key: "RANGE",
          },
        },
      },
    },
  },
}