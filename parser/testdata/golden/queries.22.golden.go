parser.row{
  Query: "SELECT * FROM gamescores WHERE UserId = ? AND TopScore > ?",
  AST: parser.Select{
    Projection: &parser.ProjectionExpression{
      All: true,
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
                    PositionalPlaceholder: &true,
                  },
                },
              },
            },
          },
        },
        &parser.Condition{
          Operand: &parser.ConditionOperand{
            Operand: &parser.DocumentPath{
              Fragment: []parser.PathFragment{
                parser.PathFragment{
                  Symbol: "TopScore",
                },
              },
            },
            ConditionRHS: &parser.ConditionRHS{
              Compare: &parser.Compare{
                Operator: ">",
                Operand: &parser.Operand{
                  Value: &parser.Value{
                    PositionalPlaceholder: &true,
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