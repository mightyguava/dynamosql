parser.row{
  Query: "SELECT * FROM movies WHERE title = :title AND attribute_exists(year)",
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
                  Symbol: "title",
                },
              },
            },
            ConditionRHS: &parser.ConditionRHS{
              Compare: &parser.Compare{
                Operator: "=",
                Operand: &parser.Operand{
                  Value: &parser.Value{
                    PlaceHolder: &":title",
                  },
                },
              },
            },
          },
        },
        &parser.Condition{
          Function: &parser.FunctionExpression{
            Function: "attribute_exists",
            Args: []*parser.FunctionArgument{
              &parser.FunctionArgument{
                DocumentPath: &parser.DocumentPath{
                  Fragment: []parser.PathFragment{
                    parser.PathFragment{
                      Symbol: "year",
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