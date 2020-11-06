parser.row{
  Query: "SELECT Studio.Name, Studio.Name.FirstName, Studio.Employees[3] FROM gamescores WHERE UserId = :UserId",
  AST: parser.AST{
    Select: &parser.Select{
      Projection: &parser.ProjectionExpression{
        Columns: []*parser.ProjectionColumn{
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
                  Symbol: "Name",
                },
                {
                  Symbol: "FirstName",
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
  },
}