package distance

import (
	"testing"
)

func TestDTW(t *testing.T) {
	tests := []struct {
		name    string
		a, b    []float64
		wantErr bool
	}{
		{"identical", []float64{1, 2, 3}, []float64{1, 2, 3}, false},
		{"different lengths", []float64{1, 2}, []float64{1, 2, 3, 4}, false},
		{"shifted", []float64{1, 2, 3}, []float64{2, 3, 4}, false},
		{"empty", []float64{}, []float64{1, 2}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DTW(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && result < 0 {
				t.Errorf("DTW distance should be non-negative, got %v", result)
			}
		})
	}
}

func TestDTWWithWindow(t *testing.T) {
	a := []float64{1, 2, 3, 4, 5}
	b := []float64{1, 2, 3, 4, 5}

	result, err := DTWWithWindow(a, b, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Identical sequences should have DTW distance close to 0
	if result > 0.1 {
		t.Errorf("expected near 0, got %v", result)
	}
}

func TestHausdorff(t *testing.T) {
	a := [][]float64{{0, 0}, {1, 0}, {0, 1}}
	b := [][]float64{{0, 0}, {1, 0}, {0, 1}}

	result, err := Hausdorff(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Identical sets should have Hausdorff distance 0
	if result > 0.1 {
		t.Errorf("expected near 0, got %v", result)
	}
}

func TestLongestCommonSubstring(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []int
		expected int
	}{
		{"identical", []int{1, 2, 3}, []int{1, 2, 3}, 3},
		{"partial match", []int{1, 2, 3, 4}, []int{2, 3, 4, 5}, 3},
		{"no match", []int{1, 2}, []int{3, 4}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LongestCommonSubstring(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestSmithWaterman(t *testing.T) {
	a := []byte{'A', 'C', 'G', 'T'}
	b := []byte{'A', 'C', 'G', 'T'}

	score, err := SmithWaterman(a, b, 2, -1, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if score <= 0 {
		t.Errorf("expected positive score for identical sequences, got %d", score)
	}
}

func TestNeedlemanWunsch(t *testing.T) {
	a := []byte{'A', 'C', 'G', 'T'}
	b := []byte{'A', 'C', 'G', 'T'}

	score, err := NeedlemanWunsch(a, b, 1, -1, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Identical sequences should have score equal to length
	if score != 4 {
		t.Errorf("expected score 4, got %d", score)
	}
}

func TestAutocorrelation(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}

	// Lag 0 should be 1.0 (perfect correlation)
	result, err := Autocorrelation(data, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !almostEqual(result, 1.0) {
		t.Errorf("lag 0 should be 1.0, got %v", result)
	}
}

func BenchmarkDTW(b *testing.B) {
	a := make([]float64, 100)
	b2 := make([]float64, 100)
	for i := range a {
		a[i] = float64(i)
		b2[i] = float64(i) + 0.5
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DTW(a, b2)
	}
}
