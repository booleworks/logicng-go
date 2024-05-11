package bdd

import "fmt"

// BDDNode represents a node in a BDD.
type BDDNode interface {
	// Label returns the label of the node.  Can either be a variable or a constant.
	Label() string

	// Reports whether the node is an inner node.
	InnerNode() bool

	// Low returns the node of the low edge or nil for a terminal node.
	Low() BDDNode

	// High returns the node of the high edge or nil for a terminal node.
	High() BDDNode

	String() string
}

// BDDConstant represents a terminal node in a BDD.
type BDDConstant struct {
	Value bool
}

// Label returns the label of the node.  Can either be a variable or a constant.
func (c BDDConstant) Label() string {
	if c.Value {
		return "$true"
	} else {
		return "$false"
	}
}

// Reports whether the node is an inner node.
func (c BDDConstant) InnerNode() bool { return false }

// Low returns the node of the low edge or nil for a terminal node.
func (c BDDConstant) Low() BDDNode { return nil }

// High returns the node of the high edge or nil for a terminal node.
func (c BDDConstant) High() BDDNode { return nil }

func (c BDDConstant) String() string { return fmt.Sprintf("<%s>", c.Label()) }

// BDDInnerNode represents an inner node in a BDD.
type BDDInnerNode struct {
	Variable string
	LowNode  BDDNode
	HighNode BDDNode
}

// Label returns the label of the node.  Can either be a variable or a constant.
func (n BDDInnerNode) Label() string { return n.Variable }

// Reports whether the node is an inner node.
func (n BDDInnerNode) InnerNode() bool { return true }

// Low returns the node of the low edge or nil for a terminal node.
func (n BDDInnerNode) Low() BDDNode { return n.LowNode }

// High returns the node of the high edge or nil for a terminal node.
func (n BDDInnerNode) High() BDDNode { return n.HighNode }

func (n BDDInnerNode) String() string {
	return fmt.Sprintf("<%s | low=%s high=%s>", n.Variable, n.LowNode.String(), n.HighNode.String())
}
