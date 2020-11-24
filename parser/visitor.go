package parser

import (
	"fmt"
	"reflect"

	"github.com/alecthomas/repr"
)

// Visit nodes in the AST
//
// The visitor can call "next()" to continue traversal of child nodes.
func Visit(node Node, visitor func(node Node, next func() error) error) error {
	return visitor(node, func() error {
		if reflect.ValueOf(node).IsNil() { // Workaround for Go's typed nil interfaces.
			return nil
		}
		switch node := node.(type) {
		case *CreateTable:
			for _, entry := range node.Entries {
				if err := Visit(entry, visitor); err != nil {
					return err
				}
			}
			return nil
		case *CreateTableEntry:
			switch {
			case node.Attr != nil:
				return Visit(node.Attr, visitor)
			case node.ProvisionedThroughput != nil:
				return Visit(node.ProvisionedThroughput, visitor)
			case node.GlobalSecondaryIndex != nil:
				return Visit(node.GlobalSecondaryIndex, visitor)
			case node.LocalSecondaryIndex != nil:
				return Visit(node.LocalSecondaryIndex, visitor)
			default:
				panic(repr.String(node))
			}
		case *TableAttr, *GlobalSecondaryIndex, *LocalSecondaryIndex, *ProvisionedThroughput:
			return nil
		case *Select:
			if err := Visit(node.Projection, visitor); err != nil {
				return err
			}
			return Visit(node.Where, visitor)
		case *ProjectionExpression:
			for _, e := range node.Columns {
				if err := Visit(e, visitor); err != nil {
					return err
				}
			}
			return nil
		case *ProjectionColumn:
			if node.DocumentPath != nil {
				return Visit(node.DocumentPath, visitor)
			}
			return Visit(node.Function, visitor)
		case *ConditionExpression:
			for _, entry := range node.Or {
				if err := Visit(entry, visitor); err != nil {
					return err
				}
			}
			return nil
		case *AndExpression:
			for _, entry := range node.And {
				if err := Visit(entry, visitor); err != nil {
					return err
				}
			}
			return nil
		case *Condition:
			switch {
			case node.Parenthesized != nil:
				return Visit(node.Parenthesized, visitor)
			case node.Not != nil:
				return Visit(node.Not, visitor)
			case node.Operand != nil:
				return Visit(node.Operand, visitor)
			case node.Function != nil:
				return Visit(node.Function, visitor)
			default:
				panic(fmt.Sprintf("invalid Condition %v", node))
			}
		case *NotCondition:
			return Visit(node.Condition, visitor)
		case *ConditionOperand:
			if err := Visit(node.Operand, visitor); err != nil {
				return err
			}
			return Visit(node.ConditionRHS, visitor)
		case *ConditionRHS:
			switch {
			case node.Between != nil:
				return Visit(node.Between, visitor)
			case node.Compare != nil:
				return Visit(node.Compare, visitor)
			case node.In != nil:
				return Visit(node.In, visitor)
			default:
				panic(fmt.Sprintf("invalid ConditionRHS %v", node))
			}
		case *Between:
			if err := Visit(node.Start, visitor); err != nil {
				return err
			}
			return Visit(node.End, visitor)
		case *Compare:
			return Visit(node.Operand, visitor)
		case *In:
			for _, entry := range node.Values {
				if err := Visit(entry, visitor); err != nil {
					return err
				}
			}
			return nil
		case *Operand:
			if node.SymbolRef != nil {
				return Visit(node.SymbolRef, visitor)
			}
			return Visit(node.Value, visitor)
		case *DocumentPath:
			for _, frag := range node.Fragment {
				if err := Visit(frag, visitor); err != nil {
					return Visit(frag, visitor)
				}
			}
			return nil
		case *FunctionExpression:
			for _, arg := range node.Args {
				if err := Visit(arg, visitor); err != nil {
					return Visit(arg, visitor)
				}
			}
			return nil
		case *FunctionArgument:
			if node.DocumentPath != nil {
				return Visit(node.DocumentPath, visitor)
			}
			return Visit(node.Value, visitor)
		case *JSONObject:
			for _, entry := range node.Entries {
				if err := Visit(entry, visitor); err != nil {
					return err
				}
			}
			return nil
		case *JSONObjectEntry:
			return Visit(node.Value, visitor)
		case *JSONArray:
			for _, entry := range node.Entries {
				if err := Visit(entry, visitor); err != nil {
					return err
				}
			}
			return nil
		case *JSONValue:
			switch {
			case node.Object != nil:
				return Visit(node.Object, visitor)
			case node.Array != nil:
				return Visit(node.Array, visitor)
			}
			return Visit(&node.Scalar, visitor)
		case *Value, *PathFragment:
			// Leaf nodes
			return nil
		default:
			panic(fmt.Sprintf("invalid node type %s", reflect.TypeOf(node)))
		}
	})
}
