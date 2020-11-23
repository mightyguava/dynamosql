parser.row{
  Query: "CREATE TABLE movies (title STRING, year NUMBER, GLOBAL SECONDARY INDEX year_title HASH(year) RANGE(title) PROJECTION ALL PROVISIONED THROUGHPUT READ 1 WRITE 1);",
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
          GlobalSecondaryIndex: &parser.GlobalSecondaryIndex{
            Name: "year_title",
            PartitionKey: "year",
            SortKey: "title",
            Projection: &parser.Projection{
              All: true,
            },
            ProvisionedThroughput: &parser.ProvisionedThroughput{
              ReadCapacityUnits: 1,
              WriteCapacityUnits: 1,
            },
          },
        },
      },
    },
  },
}