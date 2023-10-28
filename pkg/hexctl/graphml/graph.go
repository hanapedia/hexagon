package graphml

import (
	"encoding/xml"
)

type GraphMLStruct struct {
	XMLName xml.Name     `xml:"graphml"`
	Keys    []Key        `xml:"key"`
	Graph   GraphContent `xml:"graph"`
}

type Key struct {
	ID       string `xml:"id,attr"`
	For      string `xml:"for,attr"`
	AttrName string `xml:"attr.name,attr"`
	AttrType string `xml:"attr.type,attr"`
}

type GraphContent struct {
	Edges       []*Edge `xml:"edge"`
	Nodes       []*Node `xml:"node"`
	Edgedefault string  `xml:"edgedefault,attr"`
}

type Graph struct {
	Nodes map[string]*Node
	Edges map[EdgeKey]*Edge
}

type Node struct {
	ID   string `xml:"id,attr"`
	Data Data   `xml:"data"`
}

type Edge struct {
	Source string `xml:"source,attr"`
	Target string `xml:"target,attr"`
	Data   Data   `xml:"data"`
}

type Data struct {
	Key   string `xml:"key,attr"`
	Value string `xml:",chardata"`
}

type EdgeKey struct {
	Source, Target string
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
		Edges: make(map[EdgeKey]*Edge),
	}
}

func (g *Graph) AddEdge(source, target string) error {
	// Automatically add nodes if they don't exist
	if _, exists := g.Nodes[source]; !exists {
		g.Nodes[source] = &Node{ID: source}
	}
	if _, exists := g.Nodes[target]; !exists {
		g.Nodes[target] = &Node{ID: target}
	}

	key := EdgeKey{Source: source, Target: target}
	if _, exists := g.Edges[key]; exists {
		return ErrEdgeAlreadyExists
	}

	g.Edges[key] = &Edge{Source: source, Target: target}
	return nil
}

func (g *Graph) GetNode(id string) (*Node, error) {
	node, exists := g.Nodes[id]
	if !exists {
		return nil, ErrNodeDoesNotExist
	}
	return node, nil
}

func (g *Graph) GetEdge(source, target string) (*Edge, error) {
	edge, exists := g.Edges[EdgeKey{Source: source, Target: target}]
	if !exists {
		return nil, ErrEdgeDoesNotExist
	}
	return edge, nil
}

func (g *Graph) SetNodeData(id, key, value string) error {
	node, err := g.GetNode(id)
	if err != nil {
		return err
	}
	if node.Data.Key != "" || node.Data.Value != "" {
		return ErrNodeDataAlreadySet
	}
	node.Data = Data{Key: key, Value: value}
	return nil
}

func (g *Graph) SetEdgeData(source, target, key, value string) error {
	edge, err := g.GetEdge(source, target)
	if err != nil {
		return err
	}
	if edge.Data.Key != "" || edge.Data.Value != "" {
		return ErrEdgeDataAlreadySet
	}
	edge.Data = Data{Key: key, Value: value}
	return nil
}

func (g *Graph) ToGraphML() ([]byte, error) {
	graphContent := GraphContent{
		Edges:       make([]*Edge, 0, len(g.Edges)),
		Nodes:       make([]*Node, 0, len(g.Nodes)),
		Edgedefault: "directed",
	}
	keys := []Key{
		{ID: "type", For: "node", AttrName: "type", AttrType: "string"},
		{ID: "type", For: "edge", AttrName: "type", AttrType: "string"},
	}

	for _, node := range g.Nodes {
		graphContent.Nodes = append(graphContent.Nodes, node)
	}
	for _, edge := range g.Edges {
		graphContent.Edges = append(graphContent.Edges, edge)
	}

	graphML := GraphMLStruct{
		Graph: graphContent,
		Keys:  keys,
	}

	return xml.MarshalIndent(graphML, "", "  ")
}
