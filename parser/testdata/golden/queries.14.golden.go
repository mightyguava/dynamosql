parser.row{
  Query: "SELECT document(UserId, TopScore, Scores[3], Scores[3][2], Studio.Name, Studio.Location.Country, Studio.Employees[3]) FROM gamescores WHERE UserId = :UserId",
  AST: parser.Select{
    Projection: &parser.ProjectionExpression{
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
    },
    From: "gamescores",
    Where: &parser.AndExpression{
      And: []*parser.Condition{
        &parser.Condition{
          Operand: &parser.ConditionOperand{
            Operand: &parser.DocumentPath{
              Fragment: []parser.PathFragment{
                parser.PathFragment{
                  Symbol: "UserId",
                },
              },
            },
            ConditionRHS: &parser.ConditionRHS{
              Compare: &parser.Compare{
                Operator: "=",
                Operand: &parser.Operand{
                  Value: &parser.Value{
                    PlaceHolder: &":UserId",
                  },
                },
              },
            },
          },
        },
      },
    },
  },
}