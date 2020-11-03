querybuilder.item{
  Query: "SELECT document(UserId, TopScore, Scores[3], Scores[3][2], Studio.Name, Studio.Location.Country, Studio.Employees[3]) FROM gamescores WHERE UserId = :UserId",
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
      &parser.ProjectionColumn{
        Function: &parser.FunctionExpression{
          Function: "document",
          Args: []*parser.FunctionArgument{
            &parser.FunctionArgument{
              DocumentPath: &parser.DocumentPath{
                Fragment: []parser.PathFragment{
                  parser.PathFragment{
                    Symbol: "UserId",
                  },
                },
              },
            },
            &parser.FunctionArgument{
              DocumentPath: &parser.DocumentPath{
                Fragment: []parser.PathFragment{
                  parser.PathFragment{
                    Symbol: "TopScore",
                  },
                },
              },
            },
            &parser.FunctionArgument{
              DocumentPath: &parser.DocumentPath{
                Fragment: []parser.PathFragment{
                  parser.PathFragment{
                    Symbol: "Scores",
                    Indexes: []int{
                      3,
                    },
                  },
                },
              },
            },
            &parser.FunctionArgument{
              DocumentPath: &parser.DocumentPath{
                Fragment: []parser.PathFragment{
                  parser.PathFragment{
                    Symbol: "Scores",
                    Indexes: []int{
                      3,
                      2,
                    },
                  },
                },
              },
            },
            &parser.FunctionArgument{
              DocumentPath: &parser.DocumentPath{
                Fragment: []parser.PathFragment{
                  parser.PathFragment{
                    Symbol: "Studio",
                  },
                  parser.PathFragment{
                    Symbol: "Name",
                  },
                },
              },
            },
            &parser.FunctionArgument{
              DocumentPath: &parser.DocumentPath{
                Fragment: []parser.PathFragment{
                  parser.PathFragment{
                    Symbol: "Studio",
                  },
                  parser.PathFragment{
                    Symbol: "Location",
                  },
                  parser.PathFragment{
                    Symbol: "Country",
                  },
                },
              },
            },
            &parser.FunctionArgument{
              DocumentPath: &parser.DocumentPath{
                Fragment: []parser.PathFragment{
                  parser.PathFragment{
                    Symbol: "Studio",
                  },
                  parser.PathFragment{
                    Symbol: "Employees",
                    Indexes: []int{
                      3,
                    },
                  },
                },
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