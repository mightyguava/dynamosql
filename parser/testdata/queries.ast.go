[]parser.row{
  parser.row{
    Query: "SELECT * FROM movies",
    AST: parser.Select{
      Expression: &parser.SelectExpression{
        All: true,
      },
      From: &parser.From{
        Table: "movies",
      },
    },
  },
  parser.row{
    Query: "SELECT title, year FROM movies",
    AST: parser.Select{
      Expression: &parser.SelectExpression{
        Projections: []string{
          "title",
          "year",
        },
      },
      From: &parser.From{
        Table: "movies",
      },
    },
  },
  parser.row{
    Query: "SELECT title, year FROM movies WHERE title = \"The Dark Knight\"",
    AST: parser.Select{
      Expression: &parser.SelectExpression{
        Projections: []string{
          "title",
          "year",
        },
      },
      From: &parser.From{
        Table: "movies",
        Where: &parser.ConditionExpression{
          Condition: &parser.Condition{
            Operand: &parser.ConditionOperand{
              Operand: &parser.Operand{
                Value: parser.Value{
                },
                SymbolRef: &parser.SymbolRef{
                  Symbol: "title",
                },
              },
              ConditionRHS: &parser.ConditionRHS{
                Compare: &parser.Compare{
                  Operator: "=",
                  Operand: parser.Operand{
                    Value: parser.Value{
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
  parser.row{
    Query: "SELECT title, year FROM movies WHERE title = \"The Dark Knight\" AND year >= 2009",
    AST: parser.Select{
      Expression: &parser.SelectExpression{
        Projections: []string{
          "title",
          "year",
        },
      },
      From: &parser.From{
        Table: "movies",
        Where: &parser.ConditionExpression{
          Condition: &parser.Condition{
            Operand: &parser.ConditionOperand{
              Operand: &parser.Operand{
                Value: parser.Value{
                },
                SymbolRef: &parser.SymbolRef{
                  Symbol: "title",
                },
              },
              ConditionRHS: &parser.ConditionRHS{
                Compare: &parser.Compare{
                  Operator: "=",
                  Operand: parser.Operand{
                    Value: parser.Value{
                      String: &"The Dark Knight",
                    },
                  },
                },
              },
            },
          },
          More: []parser.MoreConditions{
            parser.MoreConditions{
              LogicalOperator: "AND",
              Condition: &parser.Condition{
                Operand: &parser.ConditionOperand{
                  Operand: &parser.Operand{
                    Value: parser.Value{
                    },
                    SymbolRef: &parser.SymbolRef{
                      Symbol: "year",
                    },
                  },
                  ConditionRHS: &parser.ConditionRHS{
                    Compare: &parser.Compare{
                      Operator: ">=",
                      Operand: parser.Operand{
                        Value: parser.Value{
                          Number: &2009,
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
  parser.row{
    Query: "SELECT title, year FROM movies WHERE title = \"The Dark Knight\" AND year BETWEEN 2009 AND 2015",
    AST: parser.Select{
      Expression: &parser.SelectExpression{
        Projections: []string{
          "title",
          "year",
        },
      },
      From: &parser.From{
        Table: "movies",
        Where: &parser.ConditionExpression{
          Condition: &parser.Condition{
            Operand: &parser.ConditionOperand{
              Operand: &parser.Operand{
                Value: parser.Value{
                },
                SymbolRef: &parser.SymbolRef{
                  Symbol: "title",
                },
              },
              ConditionRHS: &parser.ConditionRHS{
                Compare: &parser.Compare{
                  Operator: "=",
                  Operand: parser.Operand{
                    Value: parser.Value{
                      String: &"The Dark Knight",
                    },
                  },
                },
              },
            },
          },
          More: []parser.MoreConditions{
            parser.MoreConditions{
              LogicalOperator: "AND",
              Condition: &parser.Condition{
                Operand: &parser.ConditionOperand{
                  Operand: &parser.Operand{
                    Value: parser.Value{
                    },
                    SymbolRef: &parser.SymbolRef{
                      Symbol: "year",
                    },
                  },
                  ConditionRHS: &parser.ConditionRHS{
                    Between: &parser.Between{
                      Start: &parser.Operand{
                        Value: parser.Value{
                          Number: &2009,
                        },
                      },
                      End: &parser.Operand{
                        Value: parser.Value{
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
  },
  parser.row{
    Query: "SELECT title, year FROM movies WHERE title = \"The Dark Knight\" AND year BETWEEN 2009 AND 2015 OR actor = \"Will Smith\"",
    AST: parser.Select{
      Expression: &parser.SelectExpression{
        Projections: []string{
          "title",
          "year",
        },
      },
      From: &parser.From{
        Table: "movies",
        Where: &parser.ConditionExpression{
          Condition: &parser.Condition{
            Operand: &parser.ConditionOperand{
              Operand: &parser.Operand{
                Value: parser.Value{
                },
                SymbolRef: &parser.SymbolRef{
                  Symbol: "title",
                },
              },
              ConditionRHS: &parser.ConditionRHS{
                Compare: &parser.Compare{
                  Operator: "=",
                  Operand: parser.Operand{
                    Value: parser.Value{
                      String: &"The Dark Knight",
                    },
                  },
                },
              },
            },
          },
          More: []parser.MoreConditions{
            parser.MoreConditions{
              LogicalOperator: "AND",
              Condition: &parser.Condition{
                Operand: &parser.ConditionOperand{
                  Operand: &parser.Operand{
                    Value: parser.Value{
                    },
                    SymbolRef: &parser.SymbolRef{
                      Symbol: "year",
                    },
                  },
                  ConditionRHS: &parser.ConditionRHS{
                    Between: &parser.Between{
                      Start: &parser.Operand{
                        Value: parser.Value{
                          Number: &2009,
                        },
                      },
                      End: &parser.Operand{
                        Value: parser.Value{
                          Number: &2015,
                        },
                      },
                    },
                  },
                },
              },
            },
            parser.MoreConditions{
              LogicalOperator: "OR",
              Condition: &parser.Condition{
                Operand: &parser.ConditionOperand{
                  Operand: &parser.Operand{
                    Value: parser.Value{
                    },
                    SymbolRef: &parser.SymbolRef{
                      Symbol: "actor",
                    },
                  },
                  ConditionRHS: &parser.ConditionRHS{
                    Compare: &parser.Compare{
                      Operator: "=",
                      Operand: parser.Operand{
                        Value: parser.Value{
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
}