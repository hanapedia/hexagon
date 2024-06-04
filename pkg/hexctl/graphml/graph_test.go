package graphml

import (
	"testing"
)

func TestAddEdge(t *testing.T) {
	g := NewGraph()
	err := g.AddEdge("node1", "node2")

	if err != nil {
		t.Error("Expected no error when adding a valid edge")
	}

	if _, exists := g.Edges[EdgeKey{"node1", "node2"}]; !exists {
		t.Error("Expected edge between 'node1' and 'node2' to exist")
	}
}

func TestGetNode(t *testing.T) {
	g := NewGraph()
	err := g.AddEdge("node1", "node2")

	node, err := g.GetNode("node1")
	if err != nil {
		t.Error("Expected no error when getting an existing node")
	}
	if node.ID != "node1" {
		t.Error("Expected to get 'node1'")
	}
}

func TestGetEdge(t *testing.T) {
	g := NewGraph()
	g.AddEdge("node1", "node2")

	edge, err := g.GetEdge("node1", "node2")
	if err != nil {
		t.Error("Expected no error when getting an existing edge")
	}
	if edge.Source != "node1" || edge.Target != "node2" {
		t.Error("Expected to get edge from 'node1' to 'node2'")
	}
}

func TestSetNodeData(t *testing.T) {
	g := NewGraph()
	g.AddEdge("node1", "node2")
	err := g.SetNodeData("node1", "key1", "value1")

	if err != nil {
		t.Error("Expected no error when setting data for an existing node")
	}
	if g.Nodes["node1"].Data.Key != "key1" || g.Nodes["node1"].Data.Value != "value1" {
		t.Error("Expected 'node1' data to be set to 'key1'='value1'")
	}
}

func TestSetEdgeData(t *testing.T) {
	g := NewGraph()
	g.AddEdge("node1", "node2")
	err := g.SetEdgeData("node1", "node2", "key1", "value1")

	if err != nil {
		t.Error("Expected no error when setting data for an existing edge")
	}
	edge := g.Edges[EdgeKey{"node1", "node2"}]
	if edge.Data.Key != "key1" || edge.Data.Value != "value1" {
		t.Error("Expected edge data to be set to 'key1'='value1'")
	}
}

func TestToGraphML(t *testing.T) {
	g := NewGraph()
	g.AddEdge("node1", "node1")
	xmlData, err := g.ToGraphML()

	if err != nil {
		t.Error("Expected no error when converting to GraphML")
	}

	if len(xmlData) == 0 {
		t.Error("Expected non-empty GraphML output")
	}
}

