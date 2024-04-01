package graphical

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	defaultNodeShape     = ShapeEllipse
	defaultLineWidth     = 2
	defaultLineWidthBold = 4
)

// WriteMermaidToString writes the given graphical representation as a
// Mermaid.js file to a string.
func WriteMermaidToString(representation *Representation) string {
	buf := bytes.NewBufferString("")
	WriteMermaidToWriter(buf, representation)
	return buf.String()
}

// WriteMermaidToFileName writes the given graphical representation as a
// Mermaid.js file to a the given filename.  Returns an error when there was a
// problem writing the file.
func WriteMermaidToFileName(filename string, representation *Representation) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return WriteMermaidToWriter(file, representation)
}

// WriteMermaidToFile writes the given graphical representation as a
// Mermaid.js file to a the given file.  Returns an error when there was a
// problem writing the file.
func WriteMermaidToFile(file *os.File, representation *Representation) error {
	defer file.Close()
	return WriteMermaidToWriter(file, representation)
}

// WriteMermaidToWriter writes the given graphical representation as a
// Mermaid.js file to a the writer.  Returns an error when there was a
// problem writing to the writer.
func WriteMermaidToWriter(writer io.Writer, representation *Representation) error {
	err := writeMermaidPreamble(writer)
	if err != nil {
		return err
	}
	err = writeMermaidNodes(writer, representation)
	if err != nil {
		return err
	}
	err = writeMermaidEdges(writer, representation)
	return err
}

func writeMermaidPreamble(writer io.Writer) error {
	_, err := fmt.Fprint(writer, "graph TD\n")
	return err
}

func writeMermaidNodes(writer io.Writer, representation *Representation) error {
	for _, node := range representation.nodes {
		fmt.Fprintf(writer, "  %s\n", mermaidNodeString(node))
		nodeStyleString := mermaidNodeStyleString(node.id, node.style)
		if nodeStyleString != "" {
			fmt.Fprint(writer, nodeStyleString)
			fmt.Fprintln(writer)
		}
	}
	_, err := fmt.Fprintln(writer)
	return err
}

func writeMermaidEdges(writer io.Writer, representation *Representation) error {
	var err error
	for i, edge := range representation.edges {
		edgeSymbol := mermaidEdgeSymbolString(edge, representation.directed)
		fmt.Fprintf(writer, "  %s %s %s\n", edge.source.id, edgeSymbol, edge.destination.id)
		edgeStyleString := mermaidEdgeStyleString(i, edge.style)
		if edgeStyleString != "" {
			fmt.Fprint(writer, edgeStyleString)
			_, err = fmt.Fprintln(writer)
		}
	}
	return err
}

func mermaidNodeString(node *Node) string {
	var start, end string
	shape := defaultNodeShape
	if node.style != nil && node.style.shape != ShapeDefault {
		shape = node.style.shape
	}
	switch shape {
	case ShapeRectangle:
		start = "["
		end = "]"
	case ShapeCircle:
		start = "(("
		end = "))"
	default:
		start = "(["
		end = "])"
	}
	return fmt.Sprintf("%s%s\"%s\"%s", node.id, start, node.label, end)
}

func mermaidEdgeSymbolString(edge *Edge, directed bool) string {
	var edgeConnector string
	if directed {
		edgeConnector = "-->"
	} else {
		edgeConnector = "---"
	}
	if edge.label == "" {
		return edgeConnector
	} else {
		return fmt.Sprintf("%s|\"%s\"|", edgeConnector, edge.label)
	}
}

func mermaidNodeStyleString(id string, style *NodeStyle) string {
	if !style.HasStyle() || (style.strokeColor == "" && style.textColor == "" && style.backgroundColor == "") {
		return ""
	}
	attributes := make([]string, 0)
	if style.strokeColor != "" {
		attributes = append(attributes, fmt.Sprintf("stroke:%s", style.strokeColor))
	}
	if style.textColor != "" {
		attributes = append(attributes, fmt.Sprintf("color:%s", style.textColor))
	}
	if style.backgroundColor != "" {
		attributes = append(attributes, fmt.Sprintf("fill:%s", style.backgroundColor))
	}
	return fmt.Sprintf("    style %s %s", id, strings.Join(attributes, ","))
}

func mermaidEdgeStyleString(counter int, style *EdgeStyle) string {
	if !style.HasStyle() {
		return ""
	}
	attributes := make([]string, 0)
	if style.color != "" {
		attributes = append(attributes, fmt.Sprintf("stroke:%s", style.color))
	}
	switch style.edgeType {
	case EdgeSolid:
		attributes = append(attributes, fmt.Sprintf("stroke-width:%d", defaultLineWidth))
	case EdgeBold:
		attributes = append(attributes, fmt.Sprintf("stroke-width:%d", defaultLineWidthBold))
	case EdgeDotted:
		attributes = append(attributes, fmt.Sprintf("stroke-width:%d,stroke-dasharray:3", defaultLineWidth))
	}
	return fmt.Sprintf("    linkStyle %d %s", counter, strings.Join(attributes, ","))
}
