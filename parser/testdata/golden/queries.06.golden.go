parser.row{
  Query: "SELECT title, year FROM movies WHERE title = \"The Dark Knight\" AND (year BETWEEN 2009 AND 2015 OR actor = \"Will Smith\")",
  AST: parser.Select{
    Projection: &parser.ProjectionExpression{
      Columns: []*parser.ProjectionColumn{
        &parser.ProjectionColumn{
          DocumentPath: &parser.DocumentPath{
            Fragment: []parser.PathFragment{
              parser.PathFragment{
                Symbol: "title",
              },
            },
          },
        },
        &parser.ProjectionColumn{
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
                    String: &"The Dark Knight",
                  },
                },
              },
            },
          },
        },
        &parser.Condition{
          Parenthesized: &parser.ParenthesizedExpression{
            ConditionExpression: &parser.ConditionExpression{
              Or: []*parser.AndExpression{
                &parser.AndExpression{
                  And: []*parser.Condition{
                    &parser.Condition{
                      Operand: &parser.ConditionOperand{
                        Operand: &parser.DocumentPath{
                          Fragment: []parser.PathFragment{
                            parser.PathFragment{
                              Symbol: "year",
                            },
                          },
                        },
                        ConditionRHS: &parser.ConditionRHS{
                          Between: &parser.Between{
                            Start: &parser.Operand{
                              Value: &parser.Value{
                                Number: &2009,
                              },
                            },
                            End: &parser.Operand{
                              Value: &parser.Value{
                                Number: &2015,
                              },
                            },
                          },
                        },
                      },
                    },
                  },
                },
                &parser.AndExpression{
                  And: []*parser.Condition{
                    &parser.Condition{
                      Operand: &parser.ConditionOperand{
                        Operand: &parser.DocumentPath{
                          Fragment: []parser.PathFragment{
                            parser.PathFragment{
                              Symbol: "actor",
                            },
                          },
                        },
                        ConditionRHS: &parser.ConditionRHS{
                          Compare: &parser.Compare{
                            Operator: "=",
                            Operand: &parser.Operand{
                              Value: &parser.Value{
                                String: &"Will Smith",
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
        },
      },
    },
  },
}