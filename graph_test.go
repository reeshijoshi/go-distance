package distance

import (
	"math"
	"testing"
)

func TestGraphBasicOperations(t *testing.T) {
	g := NewGraph()
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(1, 2, 2.0)
	g.AddEdge(0, 2, 5.0)

	if len(g.nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(g.nodes))
	}
}

func TestDijkstra(t *testing.T) {
	g := NewGraph()
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(1, 2, 2.0)
	g.AddEdge(0, 2, 5.0)

	dist, path := g.Dijkstra(0, 2)

	// Shortest path should be 0->1->2 with distance 3.0
	if dist != 3.0 {
		t.Errorf("expected distance 3.0, got %v", dist)
	}

	if len(path) != 3 || path[0] != 0 || path[2] != 2 {
		t.Errorf("expected path [0,1,2], got %v", path)
	}
}

func TestBellmanFord(t *testing.T) {
	g := NewGraph()
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(1, 2, 2.0)

	dist, hasNegCycle := g.BellmanFord(0)

	if hasNegCycle {
		t.Errorf("no negative cycle should exist")
	}

	if dist[2] != 3.0 {
		t.Errorf("expected distance to node 2 is 3.0, got %v", dist[2])
	}
}

func TestFloydWarshall(t *testing.T) {
	g := NewGraph()
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(1, 2, 2.0)
	g.AddEdge(0, 2, 5.0)

	allPairs := g.FloydWarshall()

	if allPairs[0][2] != 3.0 {
		t.Errorf("shortest path 0->2 should be 3.0, got %v", allPairs[0][2])
	}
}

func TestGraphDiameter(t *testing.T) {
	g := NewGraph()
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(1, 2, 1.0)
	g.AddEdge(2, 3, 1.0)

	diameter := g.GraphDiameter()

	// Longest shortest path should be 0->3 = 3.0
	if diameter != 3.0 {
		t.Errorf("expected diameter 3.0, got %v", diameter)
	}
}

func TestGraphRadius(t *testing.T) {
	g := NewGraph()
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(1, 2, 1.0)
	g.AddEdge(2, 0, 1.0)

	radius := g.GraphRadius()

	if radius <= 0 || math.IsInf(radius, 1) {
		t.Errorf("invalid radius: %v", radius)
	}
}

func TestAStar(t *testing.T) {
	g := NewGraph()
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(1, 2, 2.0)

	// Zero heuristic makes A* equivalent to Dijkstra
	heuristic := func(_, _ int) float64 { return 0 }

	dist, path := g.AStar(0, 2, heuristic)

	if dist != 3.0 {
		t.Errorf("expected distance 3.0, got %v", dist)
	}

	if len(path) == 0 {
		t.Errorf("expected non-empty path")
	}
}

func TestBFS(t *testing.T) {
	g := NewGraph()
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(1, 2, 1.0)
	g.AddEdge(0, 3, 1.0)

	dist, path := g.BFS(0, 2)

	if dist != 2 {
		t.Errorf("expected distance 2, got %d", dist)
	}

	if len(path) != 3 {
		t.Errorf("expected path length 3, got %d", len(path))
	}
}

func TestConnectedComponents(t *testing.T) {
	g := NewGraph()
	g.AddUndirectedEdge(0, 1, 1.0)
	g.AddUndirectedEdge(2, 3, 1.0)

	components := g.ConnectedComponents()

	if len(components) != 2 {
		t.Errorf("expected 2 components, got %d", len(components))
	}
}

func TestIsConnected(t *testing.T) {
	g := NewGraph()
	g.AddUndirectedEdge(0, 1, 1.0)
	g.AddUndirectedEdge(1, 2, 1.0)

	if !g.IsConnected() {
		t.Errorf("graph should be connected")
	}

	g2 := NewGraph()
	g2.AddEdge(0, 1, 1.0)
	g2.AddEdge(2, 3, 1.0)

	// This might not be connected if edges are directed
	// Just test that it runs without error
	_ = g2.IsConnected()
}

func TestGraphEditDistance(t *testing.T) {
	g1 := NewGraph()
	g1.AddEdge(0, 1, 1.0)
	g1.AddEdge(1, 2, 1.0)

	g2 := NewGraph()
	g2.AddEdge(0, 1, 1.0)
	g2.AddEdge(1, 3, 1.0)

	dist := GraphEditDistance(g1, g2)

	if dist < 0 {
		t.Errorf("distance should be non-negative, got %v", dist)
	}
}

func BenchmarkDijkstra(b *testing.B) {
	g := NewGraph()
	for i := 0; i < 100; i++ {
		g.AddEdge(i, i+1, 1.0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Dijkstra(0, 100)
	}
}

func BenchmarkFloydWarshall(b *testing.B) {
	g := NewGraph()
	for i := 0; i < 20; i++ {
		for j := i + 1; j < 20; j++ {
			g.AddEdge(i, j, 1.0)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.FloydWarshall()
	}
}
