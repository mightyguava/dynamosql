parser.row{
  Query: "SELECT UserId, document(TopScore) FROM gamescores WHERE UserId = :UserId",
  AST: parser.Select{
    Projection: &parser.ProjectionExpression{
      Columns: []*parser.ProjectionColumn{
        {
          DocumentPath: &parser.DocumentPath{
            Fragment: []parser.PathFragment{
              {
                Symbol: "UserId",
              },
            },
          },
        },
        {
          Function: &parser.FunctionExpression{
            Function: "document",
            Args: []*parser.FunctionArgument{
              {
                DocumentPath: &parser.DocumentPath{
                  Fragment: []parser.PathFragment{
                    {
                      Symbol: "TopScore",
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
        {
          Operand: &parser.ConditionOperand{
            Operand: &parser.DocumentPath{
              Fragment: []parser.PathFragment{
                {
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