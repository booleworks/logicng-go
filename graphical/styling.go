package graphical

import (
	"fmt"

	"github.com/booleworks/logicng-go/errorx"
)

// Color used in graphical representations.
type Color string

const (
	ColorBlack     Color = "#000000"
	ColorWhite     Color = "#ffffff"
	ColorDarkGray  Color = "#777777"
	ColorLightGray Color = "#e4e4e4"
	ColorRed       Color = "#ea2027"
	ColorGreen     Color = "#009432"
	ColorBlue      Color = "#004f93"
	ColorYellow    Color = "#ffc612"
	ColorOrange    Color = "#f79f1f"
	ColorCyan      Color = "#1289a7"
	ColorPurple    Color = "#5758bb"
	ColorTurquoise Color = "#006266"
)

// ColorFromHex returns a color from a hex value.
func ColorFromHex(hexValue string) Color {
	return Color(hexValue)
}

// ColorFromRgb returns a color from red, green, blue values or an
// error if it cannot be converted to hex.
func ColorFromRgb(red, green, blue int) (Color, error) {
	if !isValidRgbValue(red) || !isValidRgbValue(green) || !isValidRgbValue(blue) {
		return ColorBlack, errorx.BadInput("invalid RGB value")
	}
	hex := fmt.Sprintf("#%02x%02x%02x", red, green, blue)
	return Color(hex), nil
}

func isValidRgbValue(rgbValue int) bool {
	return 0 <= rgbValue && rgbValue <= 255
}

// Shape of a node.
type Shape byte

const (
	ShapeDefault Shape = iota
	ShapeRectangle
	ShapeEllipse
	ShapeCircle
)

// NodeStyle gathers shape and stroke, text, and background color of a shape.
type NodeStyle struct {
	shape           Shape
	strokeColor     Color
	textColor       Color
	backgroundColor Color
}

// NoNodeStyle returns an empty node style defaulting to the backend's
// (Dot/Mermaid.js) default value.
func NoNodeStyle() *NodeStyle {
	return &NodeStyle{ShapeDefault, "", "", ""}
}

// NewNodeStyle generates a new node style with the given shape and stroke,
// text, and background color.
func NewNodeStyle(shape Shape, strokeColor, textColor, backgroundColor Color) *NodeStyle {
	if shape == ShapeDefault && strokeColor == "" && textColor == "" && backgroundColor == "" {
		return NoNodeStyle()
	} else {
		return &NodeStyle{shape, strokeColor, textColor, backgroundColor}
	}
}

// Circle generates a new node style for a circle with the given stroke, text,
// and background color.
func Circle(strokeColor, textColor, backgroundColor Color) *NodeStyle {
	return &NodeStyle{ShapeCircle, strokeColor, textColor, backgroundColor}
}

// Ellipse generates a new node style for an ellipse with the given stroke,
// text, and background color.
func Ellipse(strokeColor, textColor, backgroundColor Color) *NodeStyle {
	return &NodeStyle{ShapeEllipse, strokeColor, textColor, backgroundColor}
}

// Rectangle generates a new node style for a rectangle with the given stroke,
// text, and background color.
func Rectangle(strokeColor, textColor, backgroundColor Color) *NodeStyle {
	return &NodeStyle{ShapeRectangle, strokeColor, textColor, backgroundColor}
}

// HasStyle reports whether the node has any style.
func (s *NodeStyle) HasStyle() bool {
	return s.shape != ShapeDefault || s.strokeColor != "" || s.textColor != "" || s.backgroundColor != ""
}

// EdgeType describes the type of edge.
type EdgeType byte

const (
	EdgeDefault EdgeType = iota
	EdgeSolid
	EdgeDotted
	EdgeBold
)

// EdgeStyle gathers type and color of an edge.
type EdgeStyle struct {
	edgeType EdgeType
	color    Color
}

// NoEdgeStyle returns an empty edge style defaulting to the backend's
// (Dot/Mermaid.js) default value.
func NoEdgeStyle() *EdgeStyle {
	return &EdgeStyle{EdgeDefault, ""}
}

// NewEdgeStyle generates a new edge style with the given type and color.
func NewEdgeStyle(edgeType EdgeType, color Color) *EdgeStyle {
	if edgeType == EdgeDefault && color == "" {
		return NoEdgeStyle()
	} else {
		return &EdgeStyle{edgeType, color}
	}
}

// Solid generates a new solid edge with the given color.
func Solid(color Color) *EdgeStyle {
	return &EdgeStyle{EdgeSolid, color}
}

// Dotted generates a new dotted edge with the given color.
func Dotted(color Color) *EdgeStyle {
	return &EdgeStyle{EdgeDotted, color}
}

// Bold generates a new bold edge with the given color.
func Bold(color Color) *EdgeStyle {
	return &EdgeStyle{EdgeBold, color}
}

// HasStyle reports whether the edge has any style.
func (s *EdgeStyle) HasStyle() bool {
	return s.edgeType != EdgeDefault || s.color != ""
}
