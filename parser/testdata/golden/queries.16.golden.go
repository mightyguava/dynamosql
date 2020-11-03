parser.row{
  Query: "SELECT `SELECT`.`foo` FROM gamescores WHERE UserId = :UserId",
  AST: parser.Select{
    Projection: &parser.ProjectionExpression{
      Columns: []*parser.ProjectionColumn{
        &parser.ProjectionColumn{
          DocumentPath: &parser.DocumentPath{
            Fragment: []parser.PathFragment{
              parser.PathFragment{
                Symbol: "SELECT",
              },
              parser.PathFragment{
                Symbol: "foo",
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