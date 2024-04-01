package graphical

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// WriteDotToString writes the given graphical representation as a Dot file to
// a string.
func WriteDotToString(representation *Representation) string {
	buf := bytes.NewBufferString("")
	WriteDotToWriter(buf, representation)
	return buf.String()
}

// WriteDotToFileName writes the given graphical representation as a Dot file
// to the given filename.  Returns an error when there was a problem writing
// the file.
func WriteDotToFileName(filename string, representation *Representation) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return WriteDotToWriter(file, representation)
}

// WriteDotToFile writes the given graphical representation as a Dot file to
// the given file.  Returns an error when there was a problem writing the file.
func WriteDotToFile(file *os.File, representation *Representation) error {
	defer file.Close()
	return WriteDotToWriter(file, representation)
}

// WriteDotToWriter writes the given graphical representation as a Dot file to
// the given writer.  Returns an error when there was a problem writing to
// the writer.
func WriteDotToWriter(writer io.Writer, representation *Representation) error {
	err := writeDotPreamble(writer, representation)
	if err != nil {
		return err
	}
	err = writeDotNodes(writer, representation)
	if err != nil {
		return err
	}
	err = writeDotEdges(writer, representation)
	if err != nil {
		return err
	}
	return writeDotClosing(writer)
}

func writeDotPreamble(writer io.Writer, representation *Representation) error {
	var graphType string
	if representation.directed {
		graphType = "digraph G"
	} else {
		graphType = "strict graph"
	}
	fmt.Fprintf(writer, "%s {\n", graphType)
	if representation.background != "" {
		fmt.Fprintf(writer, "  bgcolor=\"%s\"\n", representation.background)
	}
	_, err := fmt.Fprintln(writer)
	return err
}

func writeDotNodes(writer io.Writer, representation *Representation) error {
	if representation.alignTerminals {
		fmt.Fprintln(writer, "{ rank = same;")
		for _, terminalNode := range representation.TerminalNodes() {
			fmt.Fprintln(writer, dotNodeString(terminalNode))
		}
		fmt.Fprintln(writer, "}")
		for _, nonTerminalNode := range representation.NonTerminalNodes() {
			fmt.Fprintln(writer, dotNodeString(nonTerminalNode))
		}
	} else {
		for _, node := range representation.nodes {
			fmt.Fprintln(writer, dotNodeString(node))
		}
	}
	_, err := fmt.Fprintln(writer)
	return err
}

func writeDotEdges(writer io.Writer, representation *Representation) error {
	var err error
	for _, edge := range representation.edges {
		_, err = fmt.Fprintln(writer, dotEdgeString(edge, representation.directed))
	}
	return err
}

func writeDotClosing(writer io.Writer) error {
	_, err := fmt.Fprintln(writer, "}")
	return err
}

func dotNodeString(node *Node) string {
	style := node.style
	attributes := make([]string, 0)
	if style.shape != ShapeDefault {
		attributes = append(attributes, fmt.Sprintf("shape=%s", dotShapeString(style.shape)))
	}
	if style.strokeColor != "" {
		attributes = append(attributes, fmt.Sprintf("color=\"%s\"", style.strokeColor))
	}
	if style.textColor != "" {
		attributes = append(attributes, fmt.Sprintf("fontcolor=\"%s\"", style.textColor))
	}
	if style.backgroundColor != "" {
		attributes = append(attributes, fmt.Sprintf("style=filled, fillcolor=\"%s\"", style.backgroundColor))
	}
	attributeString := ""
	if len(attributes) > 0 {
		attributeString = ", " + strings.Join(attributes, ", ")
	}
	return fmt.Sprintf("  %s [label=\"%s\"%s]", node.id, node.label, attributeString)
}

func dotEdgeString(edge *Edge, isDirected bool) string {
	style := edge.style
	attributes := make([]string, 0)
	if style.color != "" {
		attributes = append(attributes, fmt.Sprintf("color=\"%[1]s\", fontcolor=\"%[1]s\"", style.color))
	}
	if style.edgeType != EdgeDefault {
		attributes = append(attributes, fmt.Sprintf("style=%s", dotEdgeStyleString(style.edgeType)))
	}
	if edge.label != "" {
		attributes = append(attributes, fmt.Sprintf("label=\"%s\"", edge.label))
	}
	attributeString := ""
	if len(attributes) > 0 {
		attributeString = " [" + strings.Join(attributes, ", ") + "]"
	}
	var edgeSymbol string
	if isDirected {
		edgeSymbol = "->"
	} else {
		edgeSymbol = "--"
	}
	return fmt.Sprintf("  %s %s %s%s", edge.source.id, edgeSymbol, edge.destination.id, attributeString)
}

func dotShapeString(shape Shape) string {
	switch shape {
	case ShapeRectangle:
		return "box"
	case ShapeCircle:
		return "circle"
	default:
		return "ellipse"
	}
}

func dotEdgeStyleString(edgeType EdgeType) string {
	switch edgeType {
	case EdgeDotted:
		return "dotted"
	case EdgeBold:
		return "bold"
	default:
		return "solid"
	}
}
