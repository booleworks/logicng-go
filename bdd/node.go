package bdd

import "fmt"

// Node represents a node in a BDD.
type Node interface {
	// Label returns the label of the node.  Can either be a variable or a constant.
	Label() string

	// Reports whether the node is an inner node.
	InnerNode() bool

	// Low returns the node of the low edge or nil for a terminal node.
	Low() Node

	// High returns the node of the high edge or nil for a terminal node.
	High() Node

	String() string
}

// ConstantNode represents a terminal node in a BDD.
type ConstantNode struct {
	Value bool
}

// Label returns the label of the node.  Can either be a variable or a constant.
func (c ConstantNode) Label() string {
	if c.Value {
		return "$true"
	}
	return "$false"
}

// Reports whether the node is an inner node.
func (c ConstantNode) InnerNode() bool { return false }

// Low returns the node of the low edge or nil for a terminal node.
func (c ConstantNode) Low() Node { return nil }

// High returns the node of the high edge or nil for a terminal node.
func (c ConstantNode) High() Node { return nil }

func (c ConstantNode) String() string { return fmt.Sprintf("<%s>", c.Label()) }

// InnerNode represents an inner node in a BDD.
type InnerNode struct {
	Variable string
	LowNode  Node
	HighNode Node
}

// Label returns the label of the node.  Can either be a variable or a constant.
func (n InnerNode) Label() string { return n.Variable }

// Reports whether the node is an inner node.
func (n InnerNode) InnerNode() bool { return true }

// Low returns the node of the low edge or nil for a terminal node.
func (n InnerNode) Low() Node { return n.LowNode }

// High returns the node of the high edge or nil for a terminal node.
func (n InnerNode) High() Node { return n.HighNode }

func (n InnerNode) String() string {
	return fmt.Sprintf("<%s | low=%s high=%s>", n.Variable, n.LowNode.String(), n.HighNode.String())
}
