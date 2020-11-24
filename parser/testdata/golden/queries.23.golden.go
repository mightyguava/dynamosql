parser.row{
  Query: "CREATE TABLE movies (title STRING, year NUMBER) PROVISIONED THROUGHPUT READ 1 WRITE 1;",
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
      ProvisionedThroughput: &parser.ProvisionedThroughput{
        ReadCapacityUnits: 1,
        WriteCapacityUnits: 1,
      },
    },
  },
}