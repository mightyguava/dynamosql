parser.row{
  Query: "SELECT Scores[3], Scores[3][2] FROM gamescores WHERE UserId = :UserId",
  AST: parser.Select{
    Projection: &parser.ProjectionExpression{
      Columns: []*parser.ProjectionColumn{
        &parser.ProjectionColumn{
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
        &parser.ProjectionColumn{
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