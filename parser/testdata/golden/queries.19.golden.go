parser.row{
  Query: "SELECT * FROM movies WHERE UserId = True",
  AST: parser.Select{
    Projection: &parser.ProjectionExpression{
      All: true,
    },
    From: "movies",
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
                    Boolean: &parser.Boolean(true),
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