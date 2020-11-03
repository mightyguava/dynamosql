parser.row{
  Query: "SELECT Studio.Name, Studio.Name.FirstName, Studio.Employees[3] FROM gamescores WHERE UserId = :UserId",
  AST: parser.Select{
    Projection: &parser.ProjectionExpression{
      Columns: []*parser.ProjectionColumn{
        &parser.ProjectionColumn{
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
        &parser.ProjectionColumn{
          DocumentPath: &parser.DocumentPath{
            Fragment: []parser.PathFragment{
              parser.PathFragment{
                Symbol: "Studio",
              },
              parser.PathFragment{
                Symbol: "Name",
              },
              parser.PathFragment{
                Symbol: "FirstName",
              },
            },
          },
        },
        &parser.ProjectionColumn{
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