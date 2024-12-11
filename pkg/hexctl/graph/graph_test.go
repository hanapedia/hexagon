package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountReachable_SimpleGraph(t *testing.T) {
	graph := &Graph{}
	graph.AddNode("A", map[string]interface{}{"key": "valueA"})
	graph.AddNode("B", map[string]interface{}{"key": "valueB"})
	graph.AddNode("C", map[string]interface{}{"key": "valueC"})
	graph.AddEdge("A", "B")
	graph.AddEdge("B", "C")

	result := graph.CalculateRecursiveCalls()
	expected := map[string]int{
		"A": 2, // A -> B -> C
		"B": 1, // B -> C
		"C": 0, // No neighbors
	}

	assert.Equal(t, expected, result, "Recursive counts for Simple Graph are incorrect")
}

func TestCountReachable_DisconnectedGraph(t *testing.T) {
	graph := &Graph{}
	graph.AddNode("A", map[string]interface{}{"key": "valueA"})
	graph.AddNode("B", map[string]interface{}{"key": "valueB"})
	graph.AddNode("C", map[string]interface{}{"key": "valueC"})

	// No edges

	result := graph.CalculateRecursiveCalls()
	expected := map[string]int{
		"A": 0, // No neighbors
		"B": 0, // No neighbors
		"C": 0, // No neighbors
	}

	assert.Equal(t, expected, result, "Recursive counts for Disconnected Graph are incorrect")
}

func TestCountReachable_CyclicGraph(t *testing.T) {
	graph := &Graph{}
	graph.AddNode("A", map[string]interface{}{"key": "valueA"})
	graph.AddNode("B", map[string]interface{}{"key": "valueB"})
	graph.AddNode("C", map[string]interface{}{"key": "valueC"})
	graph.AddEdge("A", "B")
	graph.AddEdge("B", "C")
	graph.AddEdge("C", "A") // Creates a cycle

	result := graph.CalculateRecursiveCalls()
	expected := map[string]int{
		"A": 2, // A -> B -> C
		"B": 2, // B -> C -> A
		"C": 2, // C -> A -> B
	}

	assert.Equal(t, expected, result, "Recursive counts for Cyclic Graph are incorrect")
}

func TestCountReachable_MultiplePaths(t *testing.T) {
	graph := &Graph{}
	graph.AddNode("A", map[string]interface{}{"key": "valueA"})
	graph.AddNode("B", map[string]interface{}{"key": "valueB"})
	graph.AddNode("C", map[string]interface{}{"key": "valueC"})
	graph.AddNode("D", map[string]interface{}{"key": "valueD"})
	graph.AddEdge("A", "B")
	graph.AddEdge("A", "C")
	graph.AddEdge("B", "D")
	graph.AddEdge("C", "D")

	result := graph.CalculateRecursiveCalls()
	expected := map[string]int{
		"A": 3, // A -> B -> D, A -> C -> D
		"B": 1, // B -> D
		"C": 1, // C -> D
		"D": 0, // No neighbors
	}

	assert.Equal(t, expected, result, "Recursive counts for Multiple Paths Graph are incorrect")
}

func TestCountReachable_SelfLoop(t *testing.T) {
	graph := &Graph{}
	graph.AddNode("A", map[string]interface{}{"key": "valueA"})
	graph.AddEdge("A", "A") // Self-loop

	result := graph.CalculateRecursiveCalls()
	expected := map[string]int{
		"A": 0, // Self-loops are not double-counted due to `visited` check
	}

	assert.Equal(t, expected, result, "Recursive counts for Self-Loop Graph are incorrect")
}
