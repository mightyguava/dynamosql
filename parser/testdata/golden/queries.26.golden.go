parser.row{
  Query: "CREATE TABLE movies (title STRING, year NUMBER, LOCAL SECONDARY INDEX year_index RANGE(year) PROJECTION ALL) PROVISIONED THROUGHPUT READ 1 WRITE 1;",
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
        {
          LocalSecondaryIndex: &parser.LocalSecondaryIndex{
            Name: "year_index",
            SortKey: "year",
            Projection: &parser.Projection{
              All: true,
            },
          },
        },
      },
      ProvisionedThroughput: &parser.ProvisionedThroughput{
        ReadCapacityUnits: 1,
        WriteCapacityUnits: 1,
      },
    },
  },
}