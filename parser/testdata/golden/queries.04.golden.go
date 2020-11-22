parser.row{
  Query: "SELECT title, year FROM movies WHERE title = \"The Dark Knight\" AND year BETWEEN 2009 AND 2015",
  AST: &parser.AST{
    Select: &parser.Select{
      Projection: &parser.ProjectionExpression{
        Columns: []*parser.ProjectionColumn{
          {
            DocumentPath: &parser.DocumentPath{
              Fragment: []*parser.PathFragment{
                {
                  Symbol: "title",
                },
              },
            },
          },
          {
            DocumentPath: &parser.DocumentPath{
              Fragment: []*parser.PathFragment{
                {
                  Symbol: "year",
                },
              },
            },
          },
        },
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
                        Str: &"The Dark Knight",
                      },
                    },
                  },
                },
              },
            },
          },
          {
            Operand: &parser.ConditionOperand{
              Operand: &parser.DocumentPath{
                Fragment: []*parser.PathFragment{
                  {
                    Symbol: "year",
                  },
                },
              },
              ConditionRHS: &parser.ConditionRHS{
                Between: &parser.Between{
                  Start: &parser.Operand{
                    Value: &parser.Value{
                      Scalar: parser.Scalar{
                        Number: &2009,
                      },
                    },
                  },
                  End: &parser.Operand{
                    Value: &parser.Value{
                      Scalar: parser.Scalar{
                        Number: &2015,
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
  },
}