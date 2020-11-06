querybuilder.item{
  Query: "SELECT UserId, TopScore, Scores[3], Scores[3][2], Studio.Name, Studio.Location.Country, Studio.Employees[3] FROM gamescores WHERE UserId = :UserId",
  Prepared: &querybuilder.PreparedQuery{
    Query: &dynamodb.QueryInput{
      _: struct {}{      },
      ExpressionAttributeNames: map[string]*string{
        "#Location": &"Location",
        "#Name": &"Name",
      },
      KeyConditionExpression: &"UserId = :UserId",
      ProjectionExpression: &"UserId, TopScore, Scores[3], Scores[3][2], Studio.#Name, Studio.#Location.Country, Studio.Employees[3]",
      TableName: &"gamescores",
    },
    Columns: []*parser.ProjectionColumn{
      {
        DocumentPath: &parser.DocumentPath{
          Fragment: []*parser.PathFragment{
            {
              Symbol: "UserId",
            },
          },
        },
      },
      {
        DocumentPath: &parser.DocumentPath{
          Fragment: []*parser.PathFragment{
            {
              Symbol: "TopScore",
            },
          },
        },
      },
      {
        DocumentPath: &parser.DocumentPath{
          Fragment: []*parser.PathFragment{
            {
              Symbol: "Scores",
              Indexes: []int{
                3,
              },
            },
          },
        },
      },
      {
        DocumentPath: &parser.DocumentPath{
          Fragment: []*parser.PathFragment{
            {
              Symbol: "Scores",
              Indexes: []int{
                3,
                2,
              },
            },
          },
        },
      },
      {
        DocumentPath: &parser.DocumentPath{
          Fragment: []*parser.PathFragment{
            {
              Symbol: "Studio",
            },
            {
              Symbol: "Name",
            },
          },
        },
      },
      {
        DocumentPath: &parser.DocumentPath{
          Fragment: []*parser.PathFragment{
            {
              Symbol: "Studio",
            },
            {
              Symbol: "Location",
            },
            {
              Symbol: "Country",
            },
          },
        },
      },
      {
        DocumentPath: &parser.DocumentPath{
          Fragment: []*parser.PathFragment{
            {
              Symbol: "Studio",
            },
            {
              Symbol: "Employees",
              Indexes: []int{
                3,
              },
            },
          },
        },
      },
    },
    NamedParams: querybuilder.NamedParams{
      ":UserId": querybuilder.Empty{      },
    },
    PositionalParams: map[int]string{    },
    FixedParams: map[string]interface {}{    },
  },
}