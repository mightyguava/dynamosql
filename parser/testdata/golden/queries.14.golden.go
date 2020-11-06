parser.row{
  Query: "SELECT document(UserId, TopScore, Scores[3], Scores[3][2], Studio.Name, Studio.Location.Country, Studio.Employees[3]) FROM gamescores WHERE UserId = :UserId",
  AST: parser.Select{
    Projection: &parser.ProjectionExpression{
      Columns: []*parser.ProjectionColumn{
        {
          Function: &parser.FunctionExpression{
            Function: "document",
            Args: []*parser.FunctionArgument{
              {
                DocumentPath: &parser.DocumentPath{
                  Fragment: []*parser.PathFragment{
                    {
                      Symbol: "UserId",
                    },
                  },
                },
              },
              {
                DocumentPath: &parser.DocumentPath{
                  Fragment: []*parser.PathFragment{
                    {
                      Symbol: "TopScore",
                    },
                  },
                },
              },
              {
                DocumentPath: &parser.DocumentPath{
                  Fragment: []*parser.PathFragment{
                    {
                      Symbol: "Scores",
                      Indexes: []int{
                        3,
                      },
                    },
                  },
                },
              },
              {
                DocumentPath: &parser.DocumentPath{
                  Fragment: []*parser.PathFragment{
                    {
                      Symbol: "Scores",
                      Indexes: []int{
                        3,
                        2,
                      },
                    },
                  },
                },
              },
              {
                DocumentPath: &parser.DocumentPath{
                  Fragment: []*parser.PathFragment{
                    {
                      Symbol: "Studio",
                    },
                    {
                      Symbol: "Name",
                    },
                  },
                },
              },
              {
                DocumentPath: &parser.DocumentPath{
                  Fragment: []*parser.PathFragment{
                    {
                      Symbol: "Studio",
                    },
                    {
                      Symbol: "Location",
                    },
                    {
                      Symbol: "Country",
                    },
                  },
                },
              },
              {
                DocumentPath: &parser.DocumentPath{
                  Fragment: []*parser.PathFragment{
                    {
                      Symbol: "Studio",
                    },
                    {
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
        },
      },
    },
    From: "gamescores",
    Where: &parser.AndExpression{
      And: []*parser.Condition{
        {
          Operand: &parser.ConditionOperand{
            Operand: &parser.DocumentPath{
              Fragment: []*parser.PathFragment{
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