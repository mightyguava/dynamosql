parser.row{
  Query: "SELECT * FROM movies WHERE title = :title AND begins_with(actor, \"Will\")",
  AST: &parser.AST{
    Select: &parser.Select{
      Projection: &parser.ProjectionExpression{
        All: true,
      },
      From: "movies",
      Where: &parser.AndExpression{
        And: []*parser.Condition{
          {
            Operand: &parser.ConditionOperand{
              Operand: &parser.DocumentPath{
                Fragment: []*parser.PathFragment{
                  {
                    Symbol: "title",
                  },
                },
              },
              ConditionRHS: &parser.ConditionRHS{
                Compare: &parser.Compare{
                  Operator: "=",
                  Operand: &parser.Operand{
                    Value: &parser.Value{
                      Scalar: parser.Scalar{
                      },
                      PlaceHolder: &":title",
                    },
                  },
                },
              },
            },
          },
          {
            Function: &parser.FunctionExpression{
              Function: "begins_with",
              Args: []*parser.FunctionArgument{
                {
                  DocumentPath: &parser.DocumentPath{
                    Fragment: []*parser.PathFragment{
                      {
                        Symbol: "actor",
                      },
                    },
                  },
                },
                {
                  Value: &parser.Value{
                    Scalar: parser.Scalar{
                      Str: &"Will",
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
}