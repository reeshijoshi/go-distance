package distance

import (
	"container/heap"
	"math"
)

// Graph represents a weighted graph for distance calculations
type Graph struct {
	adjacency map[int]map[int]float64 // adjacency[from][to] = weight
	nodes     map[int]bool
}

// NewGraph creates a new graph
func NewGraph() *Graph {
	return &Graph{
		adjacency: make(map[int]map[int]float64),
		nodes:     make(map[int]bool),
	}
}

// AddEdge adds a weighted edge between two nodes
func (g *Graph) AddEdge(from, to int, weight float64) {
	g.nodes[from] = true
	g.nodes[to] = true

	if g.adjacency[from] == nil {
		g.adjacency[from] = make(map[int]float64)
	}
	g.adjacency[from][to] = weight
}

// AddUndirectedEdge adds an undirected edge
func (g *Graph) AddUndirectedEdge(a, b int, weight float64) {
	g.AddEdge(a, b, weight)
	g.AddEdge(b, a, weight)
}

// Dijkstra computes shortest path distance from source to target
// Returns distance and path. Returns inf if no path exists.
// Time: O((V+E)logV), Space: O(V)
func (g *Graph) Dijkstra(source, target int) (float64, []int) {
	dist := make(map[int]float64)
	prev := make(map[int]int)
	visited := make(map[int]bool)

	// Initialize distances
	for node := range g.nodes {
		dist[node] = math.Inf(1)
	}
	dist[source] = 0

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &item{node: source, priority: 0})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*item)
		node := current.node

		if visited[node] {
			continue
		}
		visited[node] = true

		if node == target {
			break
		}

		for neighbor, weight := range g.adjacency[node] {
			if visited[neighbor] {
				continue
			}

			newDist := dist[node] + weight
			if newDist < dist[neighbor] {
				dist[neighbor] = newDist
				prev[neighbor] = node
				heap.Push(pq, &item{node: neighbor, priority: newDist})
			}
		}
	}

	// Reconstruct path
	path := []int{}
	if dist[target] != math.Inf(1) {
		node := target
		for node != source {
			path = append([]int{node}, path...)
			node = prev[node]
		}
		path = append([]int{source}, path...)
	}

	return dist[target], path
}

// BellmanFord computes shortest paths handling negative weights
// Returns distances and whether negative cycle exists
// Time: O(VE), Space: O(V)
func (g *Graph) BellmanFord(source int) (map[int]float64, bool) {
	dist := make(map[int]float64)
	for node := range g.nodes {
		dist[node] = math.Inf(1)
	}
	dist[source] = 0

	// Relax edges V-1 times
	nodeCount := len(g.nodes)
	for i := 0; i < nodeCount-1; i++ {
		for from := range g.adjacency {
			for to, weight := range g.adjacency[from] {
				if dist[from]+weight < dist[to] {
					dist[to] = dist[from] + weight
				}
			}
		}
	}

	// Check for negative cycles
	hasNegativeCycle := false
	for from := range g.adjacency {
		for to, weight := range g.adjacency[from] {
			if dist[from]+weight < dist[to] {
				hasNegativeCycle = true
				break
			}
		}
		if hasNegativeCycle {
			break
		}
	}

	return dist, hasNegativeCycle
}

// FloydWarshall computes all-pairs shortest paths
// Time: O(V³), Space: O(V²)
func (g *Graph) FloydWarshall() map[int]map[int]float64 {
	dist := make(map[int]map[int]float64)
	nodes := []int{}

	// Initialize
	for node := range g.nodes {
		nodes = append(nodes, node)
		dist[node] = make(map[int]float64)
		for other := range g.nodes {
			if node == other {
				dist[node][other] = 0
			} else {
				dist[node][other] = math.Inf(1)
			}
		}
	}

	// Set edge weights
	for from := range g.adjacency {
		for to, weight := range g.adjacency[from] {
			dist[from][to] = weight
		}
	}

	// Floyd-Warshall
	for _, k := range nodes {
		for _, i := range nodes {
			for _, j := range nodes {
				if dist[i][k]+dist[k][j] < dist[i][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}

	return dist
}

// GraphDiameter computes the diameter (maximum shortest path)
// Time: O(V³), Space: O(V²)
func (g *Graph) GraphDiameter() float64 {
	allPairs := g.FloydWarshall()
	diameter := 0.0

	for from := range allPairs {
		for to := range allPairs[from] {
			d := allPairs[from][to]
			if !math.IsInf(d, 1) && d > diameter {
				diameter = d
			}
		}
	}

	return diameter
}

// GraphRadius computes the radius (minimum eccentricity)
// Time: O(V³), Space: O(V²)
func (g *Graph) GraphRadius() float64 {
	allPairs := g.FloydWarshall()
	radius := math.Inf(1)

	for from := range allPairs {
		eccentricity := 0.0
		for to := range allPairs[from] {
			d := allPairs[from][to]
			if !math.IsInf(d, 1) && d > eccentricity {
				eccentricity = d
			}
		}
		if eccentricity < radius {
			radius = eccentricity
		}
	}

	return radius
}

// AStar computes shortest path using A* with heuristic
// Time: O(E log V) with good heuristic, Space: O(V)
func (g *Graph) AStar(source, target int, heuristic func(int, int) float64) (float64, []int) {
	dist := make(map[int]float64)
	prev := make(map[int]int)
	visited := make(map[int]bool)

	for node := range g.nodes {
		dist[node] = math.Inf(1)
	}
	dist[source] = 0

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &item{node: source, priority: heuristic(source, target)})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*item)
		node := current.node

		if visited[node] {
			continue
		}
		visited[node] = true

		if node == target {
			break
		}

		for neighbor, weight := range g.adjacency[node] {
			if visited[neighbor] {
				continue
			}

			newDist := dist[node] + weight
			if newDist < dist[neighbor] {
				dist[neighbor] = newDist
				prev[neighbor] = node
				priority := newDist + heuristic(neighbor, target)
				heap.Push(pq, &item{node: neighbor, priority: priority})
			}
		}
	}

	// Reconstruct path
	path := []int{}
	if dist[target] != math.Inf(1) {
		node := target
		for node != source {
			path = append([]int{node}, path...)
			node = prev[node]
		}
		path = append([]int{source}, path...)
	}

	return dist[target], path
}

// ResistanceDistance computes approximate effective resistance between nodes.
// WARNING: This is a simplified approximation using shortest path distance.
// A full implementation requires computing the Moore-Penrose pseudoinverse
// of the graph Laplacian matrix, which is computationally expensive.
// For accurate resistance distance, use a specialized linear algebra library.
// Time: O((V+E)logV), Space: O(V)
func (g *Graph) ResistanceDistance(source, target int) float64 {
	// Return shortest path distance as approximation
	// This provides a lower bound on the true resistance distance
	dist, _ := g.Dijkstra(source, target)
	return dist
}

// CommuteTime computes approximate expected commute time for random walk.
// WARNING: This is a simplified approximation using shortest path distance.
// True commute time requires computing hitting times using the fundamental
// matrix of the random walk, which involves matrix inversion.
// For accurate commute time, use a specialized graph analysis library.
// Time: O((V+E)logV), Space: O(V)
func (g *Graph) CommuteTime(source, target int) float64 {
	// Return twice the shortest path as a rough approximation
	// This provides a lower bound estimate of the actual commute time
	dist, _ := g.Dijkstra(source, target)
	return dist * 2
}

// GraphEditDistance computes graph edit distance between two graphs
// Time: Exponential (NP-hard), Space: O(V²)
func GraphEditDistance(g1, g2 *Graph) float64 {
	// Simplified version: count node/edge differences
	nodeCount := 0.0
	for node := range g1.nodes {
		if !g2.nodes[node] {
			nodeCount++
		}
	}
	for node := range g2.nodes {
		if !g1.nodes[node] {
			nodeCount++
		}
	}

	edgeCount := 0.0
	for from := range g1.adjacency {
		for to := range g1.adjacency[from] {
			if g2.adjacency[from] == nil || g2.adjacency[from][to] == 0 {
				edgeCount++
			}
		}
	}
	for from := range g2.adjacency {
		for to := range g2.adjacency[from] {
			if g1.adjacency[from] == nil || g1.adjacency[from][to] == 0 {
				edgeCount++
			}
		}
	}

	return nodeCount + edgeCount
}

// Priority queue for Dijkstra/A*
type item struct {
	node     int
	priority float64
	index    int
}

type priorityQueue []*item

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// BFS computes shortest path in unweighted graph
// Time: O(V+E), Space: O(V)
func (g *Graph) BFS(source, target int) (int, []int) {
	if source == target {
		return 0, []int{source}
	}

	visited := make(map[int]bool)
	queue := []int{source}
	parent := make(map[int]int)
	visited[source] = true
	distance := make(map[int]int)
	distance[source] = 0

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node == target {
			break
		}

		for neighbor := range g.adjacency[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = node
				distance[neighbor] = distance[node] + 1
				queue = append(queue, neighbor)
			}
		}
	}

	// Reconstruct path
	path := []int{}
	if visited[target] {
		node := target
		for node != source {
			path = append([]int{node}, path...)
			node = parent[node]
		}
		path = append([]int{source}, path...)
		return distance[target], path
	}

	return -1, nil
}

// ConnectedComponents finds connected components
// Time: O(V+E), Space: O(V)
func (g *Graph) ConnectedComponents() [][]int {
	visited := make(map[int]bool)
	components := [][]int{}

	var dfs func(int, *[]int)
	dfs = func(node int, component *[]int) {
		visited[node] = true
		*component = append(*component, node)
		for neighbor := range g.adjacency[node] {
			if !visited[neighbor] {
				dfs(neighbor, component)
			}
		}
	}

	for node := range g.nodes {
		if !visited[node] {
			component := []int{}
			dfs(node, &component)
			components = append(components, component)
		}
	}

	return components
}

// IsConnected checks if graph is connected
// Time: O(V+E), Space: O(V)
func (g *Graph) IsConnected() bool {
	components := g.ConnectedComponents()
	return len(components) == 1
}
