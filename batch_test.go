package distance

import (
	"context"
	"testing"
	"time"
)

func TestBatchCompute(t *testing.T) {
	vectors := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	result, err := BatchCompute(vectors, Euclidean[float64])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 3 || len(result[0]) != 3 {
		t.Errorf("expected 3x3 matrix, got %dx%d", len(result), len(result[0]))
	}

	// Diagonal should be zero
	for i := 0; i < 3; i++ {
		if result[i][i] != 0 {
			t.Errorf("diagonal[%d] should be 0, got %v", i, result[i][i])
		}
	}

	// Should be symmetric
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if !almostEqual(result[i][j], result[j][i]) {
				t.Errorf("not symmetric: result[%d][%d]=%v, result[%d][%d]=%v",
					i, j, result[i][j], j, i, result[j][i])
			}
		}
	}
}

func TestBatchComputeParallel(t *testing.T) {
	vectors := [][]float64{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
	}

	result, err := BatchComputeParallel(vectors, Euclidean[float64], 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 4 {
		t.Errorf("expected 4x4 matrix, got %dx%d", len(result), len(result[0]))
	}
}

func TestKNearestNeighbors(t *testing.T) {
	vectors := [][]float64{
		{0, 0},
		{1, 1},
		{2, 2},
		{10, 10},
	}

	result, err := KNearestNeighbors(vectors, 2, Euclidean[float64])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// For point 0, nearest neighbors should be 1 and 2
	if len(result[0]) != 2 {
		t.Errorf("expected 2 neighbors, got %d", len(result[0]))
	}
}

func TestRadiusNeighbors(t *testing.T) {
	vectors := [][]float64{
		{0, 0},
		{1, 0},
		{0, 1},
		{10, 10},
	}

	result, err := RadiusNeighbors(vectors, 2.0, Euclidean[float64])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Point 0 should have 2 neighbors within radius 2
	if len(result[0]) != 2 {
		t.Errorf("expected 2 neighbors, got %d", len(result[0]))
	}
}

func TestComputeToPoint(t *testing.T) {
	vectors := [][]float64{
		{1, 2},
		{3, 4},
		{5, 6},
	}
	point := []float64{0, 0}

	result, err := ComputeToPoint(vectors, point, Euclidean[float64])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 distances, got %d", len(result))
	}
}

func TestComputeWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	a := []float64{1, 2, 3}
	b := []float64{4, 5, 6}

	result, err := ComputeWithContext(ctx, a, b, Euclidean[float64])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result <= 0 {
		t.Errorf("expected positive distance, got %v", result)
	}
}

func TestNearestNeighbor(t *testing.T) {
	vectors := [][]float64{
		{0, 0},
		{1, 1},
		{5, 5},
	}
	query := []float64{0.9, 0.9}

	idx, dist, err := NearestNeighbor(vectors, query, Euclidean[float64])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Nearest should be index 1 (1,1)
	if idx != 1 {
		t.Errorf("expected nearest neighbor index 1, got %d", idx)
	}

	if dist <= 0 {
		t.Errorf("expected positive distance, got %v", dist)
	}
}

func TestCentroid(t *testing.T) {
	vectors := [][]float64{
		{0, 0},
		{2, 0},
		{0, 2},
		{2, 2},
	}

	result, err := Centroid(vectors)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Centroid should be (1, 1)
	if !almostEqual(result[0], 1.0) || !almostEqual(result[1], 1.0) {
		t.Errorf("expected centroid [1, 1], got %v", result)
	}
}

func TestPairwiseDistinctCount(t *testing.T) {
	vectors := [][]float64{
		{0, 0},
		{0, 0},
		{10, 10},
	}

	count, err := PairwiseDistinctCount(vectors, 5.0, Euclidean[float64])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only the pair (0,2) and (1,2) should have distance > 5
	if count != 2 {
		t.Errorf("expected 2 distinct pairs, got %d", count)
	}
}

// Benchmarks
func BenchmarkBatchCompute(b *testing.B) {
	vectors := make([][]float64, 50)
	for i := range vectors {
		vectors[i] = []float64{float64(i), float64(i * 2)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = BatchCompute(vectors, Euclidean[float64])
	}
}

func BenchmarkBatchComputeParallel(b *testing.B) {
	vectors := make([][]float64, 50)
	for i := range vectors {
		vectors[i] = []float64{float64(i), float64(i * 2)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = BatchComputeParallel(vectors, Euclidean[float64], 4)
	}
}
