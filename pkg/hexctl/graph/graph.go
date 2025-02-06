// Graph representation of service units
// Uses more efficient data structure, Adjacency List, compared to graphml package that uses list of nodes and edges.
package graph

type Graph struct {
	Nodes     map[string]*Node    // Node ID -> Node
	Adjacency map[string][]string // Node ID -> List of neighbors
}

type Node struct {
	ID   string
	Data map[string]interface{}
}

// Adds a node to the graph
func (g *Graph) AddNode(nodeID string, data map[string]interface{}) {
	if g.Nodes == nil {
		g.Nodes = make(map[string]*Node)
	}
	if g.Adjacency == nil {
		g.Adjacency = make(map[string][]string)
	}
	g.Nodes[nodeID] = &Node{ID: nodeID, Data: data}
}

// Adds a directed edge to the graph
func (g *Graph) AddEdge(source, target string) {
	if g.Adjacency == nil {
		g.Adjacency = make(map[string][]string)
	}
	g.Adjacency[source] = append(g.Adjacency[source], target)
}

// DFS to count reachable nodes
func (g *Graph) CountReachable(nodeID string, visited map[string]bool) int {
	if visited[nodeID] {
		return 0
	}
	visited[nodeID] = true

	count := 0
	for _, neighbor := range g.Adjacency[nodeID] {
		if !visited[neighbor] {
			count += 1 + g.CountReachable(neighbor, visited)
		}
	}
	return count
}

// Wrapper function to calculate recursive calls for all nodes
func (g *Graph) CalculateRecursiveCalls() map[string]int {
	result := make(map[string]int)
	for nodeID := range g.Nodes {
		visited := make(map[string]bool)
		result[nodeID] = g.CountReachable(nodeID, visited)
	}
	return result
}

func (g *Graph) CountNodesAndEdges() (int, int) {
	numNodes := len(g.Nodes)
	numEdges := 0
	for _, neighbors := range g.Adjacency {
		numEdges += len(neighbors)
	}
	return numNodes, numEdges
}

// GetSourceNodesMap returns a map of node IDs to Nodes where each node has no incoming edges.
func (g *Graph) GetSourceNodesMap() map[string]*Node {
	// Create a map to keep track of the in-degree for each node.
	inDegree := make(map[string]int)

	// Initialize the in-degree for every node in the graph to 0.
	for id := range g.Nodes {
		inDegree[id] = 0
	}

	// Iterate over the Adjacency map and count incoming edges for each node.
	for _, neighbors := range g.Adjacency {
		for _, neighborID := range neighbors {
			inDegree[neighborID]++
		}
	}

	// Build the result map containing only nodes with an in-degree of 0.
	sourceNodes := make(map[string]*Node)
	for id, count := range inDegree {
		if count == 0 {
			sourceNodes[id] = g.Nodes[id]
		}
	}

	return sourceNodes
}
