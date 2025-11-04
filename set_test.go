package distance

import (
	"testing"
)

func TestJaccardSet(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []string
		expected float64
		wantErr  bool
	}{
		{
			name:     "identical sets",
			a:        []string{"a", "b", "c"},
			b:        []string{"a", "b", "c"},
			expected: 0.0,
			wantErr:  false,
		},
		{
			name:     "no overlap",
			a:        []string{"a", "b", "c"},
			b:        []string{"d", "e", "f"},
			expected: 1.0,
			wantErr:  false,
		},
		{
			name:     "partial overlap",
			a:        []string{"a", "b", "c"},
			b:        []string{"b", "c", "d"},
			expected: 0.5,
			wantErr:  false,
		},
		{
			name:     "empty sets",
			a:        []string{},
			b:        []string{},
			expected: 0.0,
			wantErr:  false,
		},
		{
			name:     "one empty set",
			a:        []string{"a", "b"},
			b:        []string{},
			expected: 1.0,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := JaccardSet(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && !almostEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestJaccardSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []int
		expected float64
	}{
		{
			name:     "identical sets",
			a:        []int{1, 2, 3},
			b:        []int{1, 2, 3},
			expected: 1.0,
		},
		{
			name:     "no overlap",
			a:        []int{1, 2, 3},
			b:        []int{4, 5, 6},
			expected: 0.0,
		},
		{
			name:     "partial overlap",
			a:        []int{1, 2, 3, 4},
			b:        []int{3, 4, 5, 6},
			expected: 0.333,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := JaccardSimilarity(tt.a, tt.b)
			if !almostEqualTolerance(result, tt.expected, 0.01) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDiceSorensen(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []string
		expected float64
	}{
		{
			name:     "identical sets",
			a:        []string{"a", "b", "c"},
			b:        []string{"a", "b", "c"},
			expected: 1.0,
		},
		{
			name:     "no overlap",
			a:        []string{"a", "b"},
			b:        []string{"c", "d"},
			expected: 0.0,
		},
		{
			name:     "partial overlap",
			a:        []string{"a", "b", "c"},
			b:        []string{"b", "c", "d"},
			expected: 0.666,
		},
		{
			name:     "empty sets",
			a:        []string{},
			b:        []string{},
			expected: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := DiceSorensen(tt.a, tt.b)
			if !almostEqualTolerance(result, tt.expected, 0.01) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDiceDistance(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"b", "c", "d"}

	result, err := DiceDistance(a, b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Dice similarity is 2*2/(3+3) = 0.666, so distance is 0.333
	expected := 0.333
	if !almostEqualTolerance(result, expected, 0.01) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestOverlapCoefficient(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []int
		expected float64
	}{
		{
			name:     "identical sets",
			a:        []int{1, 2, 3},
			b:        []int{1, 2, 3},
			expected: 1.0,
		},
		{
			name:     "one is subset of other",
			a:        []int{1, 2},
			b:        []int{1, 2, 3, 4},
			expected: 1.0,
		},
		{
			name:     "partial overlap",
			a:        []int{1, 2, 3, 4},
			b:        []int{3, 4, 5, 6},
			expected: 0.5,
		},
		{
			name:     "no overlap",
			a:        []int{1, 2},
			b:        []int{3, 4},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := OverlapCoefficient(tt.a, tt.b)
			if !almostEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTanimotoCoefficient(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float64
		expected float64
	}{
		{
			name:     "identical vectors",
			a:        []float64{1, 1, 1},
			b:        []float64{1, 1, 1},
			expected: 1.0,
		},
		{
			name:     "binary vectors",
			a:        []float64{1, 0, 1, 0},
			b:        []float64{1, 1, 0, 0},
			expected: 0.333,
		},
		{
			name:     "zero vectors",
			a:        []float64{0, 0, 0},
			b:        []float64{0, 0, 0},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := TanimotoCoefficient(tt.a, tt.b)
			if !almostEqualTolerance(result, tt.expected, 0.01) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTanimotoDistance(t *testing.T) {
	a := []float64{1, 0, 1, 0}
	b := []float64{1, 1, 0, 0}

	result, err := TanimotoDistance(a, b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Tanimoto coefficient is 0.333, so distance is 0.666
	expected := 0.666
	if !almostEqualTolerance(result, expected, 0.01) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestCosineSimilaritySet(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []string
		expected float64
	}{
		{
			name:     "identical sets",
			a:        []string{"a", "b", "c"},
			b:        []string{"a", "b", "c"},
			expected: 1.0,
		},
		{
			name:     "no overlap",
			a:        []string{"a", "b"},
			b:        []string{"c", "d"},
			expected: 0.0,
		},
		{
			name:     "with frequencies",
			a:        []string{"a", "a", "b"},
			b:        []string{"a", "b", "b"},
			expected: 0.83,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := CosineSimilaritySet(tt.a, tt.b)
			if !almostEqualTolerance(result, tt.expected, 0.1) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCosineDistanceSet(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"a", "b", "c"}

	result, err := CosineDistanceSet(a, b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Identical sets should have distance 0
	expected := 0.0
	if !almostEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Test with different types
func TestJaccardSetWithInts(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{4, 5, 6, 7, 8}

	result, err := JaccardSet(a, b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Intersection: {4, 5} = 2
	// Union: {1,2,3,4,5,6,7,8} = 8
	// Jaccard distance = 1 - 2/8 = 0.75
	expected := 0.75
	if !almostEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Helper function for looser tolerance
func almostEqualTolerance(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < tolerance
}

// Benchmarks
func BenchmarkJaccardSet(b *testing.B) {
	a := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	b2 := []string{"e", "f", "g", "h", "i", "j", "k", "l", "m", "n"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = JaccardSet(a, b2)
	}
}

func BenchmarkDiceSorensen(b *testing.B) {
	a := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	b2 := []string{"e", "f", "g", "h", "i", "j", "k", "l", "m", "n"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DiceSorensen(a, b2)
	}
}

func BenchmarkCosineSimilaritySet(b *testing.B) {
	a := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	b2 := []string{"e", "f", "g", "h", "i", "j", "k", "l", "m", "n"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = CosineSimilaritySet(a, b2)
	}
}
