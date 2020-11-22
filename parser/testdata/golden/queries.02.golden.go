parser.row{
  Query: "SELECT title, year FROM movies WHERE title = \"The Dark Knight\"",
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
                      String: &"The Dark Knight",
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