package distance

import (
	"math"
	"testing"
)

const epsilon = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestEuclidean(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float64
		expected float64
		wantErr  bool
	}{
		{"identical vectors", []float64{1, 2, 3}, []float64{1, 2, 3}, 0, false},
		{"simple distance", []float64{0, 0}, []float64{3, 4}, 5, false},
		{"3D distance", []float64{1, 0, 0}, []float64{0, 1, 0}, math.Sqrt(2), false},
		{"dimension mismatch", []float64{1, 2}, []float64{1, 2, 3}, 0, true},
		{"empty vectors", []float64{}, []float64{}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Euclidean(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && !almostEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestManhattan(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float64
		expected float64
	}{
		{"identical vectors", []float64{1, 2, 3}, []float64{1, 2, 3}, 0},
		{"simple distance", []float64{0, 0}, []float64{3, 4}, 7},
		{"negative values", []float64{-1, -2}, []float64{1, 2}, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := Manhattan(tt.a, tt.b)
			if !almostEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestChebyshev(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float64
		expected float64
	}{
		{"identical vectors", []float64{1, 2, 3}, []float64{1, 2, 3}, 0},
		{"simple distance", []float64{0, 0, 0}, []float64{1, 5, 3}, 5},
		{"negative values", []float64{-5, 0}, []float64{5, 0}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := Chebyshev(tt.a, tt.b)
			if !almostEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMinkowski(t *testing.T) {
	a := []float64{0, 0}
	b := []float64{3, 4}

	// p=1 should equal Manhattan
	result, _ := Minkowski(a, b, 1)
	if !almostEqual(result, 7) {
		t.Errorf("Minkowski p=1: expected 7, got %v", result)
	}

	// p=2 should equal Euclidean
	result, _ = Minkowski(a, b, 2)
	if !almostEqual(result, 5) {
		t.Errorf("Minkowski p=2: expected 5, got %v", result)
	}

	// p=inf should equal Chebyshev
	result, _ = Minkowski(a, b, math.Inf(1))
	if !almostEqual(result, 4) {
		t.Errorf("Minkowski p=inf: expected 4, got %v", result)
	}
}

func TestCosine(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float64
		expected float64
		wantErr  bool
	}{
		{"identical vectors", []float64{1, 2, 3}, []float64{1, 2, 3}, 0, false},
		{"orthogonal vectors", []float64{1, 0}, []float64{0, 1}, 1, false},
		{"opposite vectors", []float64{1, 0}, []float64{-1, 0}, 2, false},
		{"zero vector", []float64{0, 0}, []float64{1, 1}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Cosine(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && !almostEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCosineSimilarity(t *testing.T) {
	// Same direction
	result, _ := CosineSimilarity([]float64{1, 2, 3}, []float64{2, 4, 6})
	if !almostEqual(result, 1) {
		t.Errorf("same direction: expected 1, got %v", result)
	}

	// Orthogonal
	result, _ = CosineSimilarity([]float64{1, 0, 0}, []float64{0, 1, 0})
	if !almostEqual(result, 0) {
		t.Errorf("orthogonal: expected 0, got %v", result)
	}

	// Opposite
	result, _ = CosineSimilarity([]float64{1, 0}, []float64{-1, 0})
	if !almostEqual(result, -1) {
		t.Errorf("opposite: expected -1, got %v", result)
	}
}

func TestHamming(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []int
		expected float64
	}{
		{"identical", []int{1, 2, 3}, []int{1, 2, 3}, 0},
		{"one difference", []int{1, 2, 3}, []int{1, 2, 4}, 1},
		{"all different", []int{1, 2, 3}, []int{4, 5, 6}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := Hamming(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCanberra(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{2, 3, 4}
	result, _ := Canberra(a, b)

	// Expected: |1-2|/(1+2) + |2-3|/(2+3) + |3-4|/(3+4) = 1/3 + 1/5 + 1/7
	expected := 1.0/3.0 + 1.0/5.0 + 1.0/7.0
	if !almostEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestBrayCurtis(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{2, 3, 4}
	result, _ := BrayCurtis(a, b)

	// Expected: (|1-2| + |2-3| + |3-4|) / (|1+2| + |2+3| + |3+4|) = 3 / 15 = 0.2
	expected := 3.0 / 15.0
	if !almostEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestWeightedEuclidean(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{4, 5, 6}
	weights := []float64{1, 2, 3}

	result, _ := WeightedEuclidean(a, b, weights)

	// Expected: sqrt(1*(1-4)^2 + 2*(2-5)^2 + 3*(3-6)^2) = sqrt(9 + 18 + 27) = sqrt(54)
	expected := math.Sqrt(54)
	if !almostEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestNorm(t *testing.T) {
	v := []float64{3, 4}

	// L1 norm
	result, _ := Norm(v, 1)
	if !almostEqual(result, 7) {
		t.Errorf("L1 norm: expected 7, got %v", result)
	}

	// L2 norm
	result, _ = Norm(v, 2)
	if !almostEqual(result, 5) {
		t.Errorf("L2 norm: expected 5, got %v", result)
	}

	// L-inf norm
	result, _ = Norm(v, math.Inf(1))
	if !almostEqual(result, 4) {
		t.Errorf("L-inf norm: expected 4, got %v", result)
	}
}

// Benchmarks
func BenchmarkEuclidean(b *testing.B) {
	v1 := make([]float64, 1000)
	v2 := make([]float64, 1000)
	for i := range v1 {
		v1[i] = float64(i)
		v2[i] = float64(i + 1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Euclidean(v1, v2)
	}
}

func BenchmarkCosine(b *testing.B) {
	v1 := make([]float64, 1000)
	v2 := make([]float64, 1000)
	for i := range v1 {
		v1[i] = float64(i)
		v2[i] = float64(i + 1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Cosine(v1, v2)
	}
}
