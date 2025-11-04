package distance

import (
	"math"
	"testing"
)

func TestKLDivergence(t *testing.T) {
	tests := []struct {
		name     string
		p, q     []float64
		expected float64
		wantErr  bool
		isInf    bool
	}{
		{
			name:     "identical distributions",
			p:        []float64{0.5, 0.5},
			q:        []float64{0.5, 0.5},
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "different distributions",
			p:        []float64{0.9, 0.1},
			q:        []float64{0.1, 0.9},
			expected: 1.76,
			wantErr:  false,
		},
		{
			name:    "zero in q causes infinite divergence",
			p:       []float64{0.5, 0.5},
			q:       []float64{1.0, 0.0},
			wantErr: false,
			isInf:   true,
		},
		{
			name:    "negative values",
			p:       []float64{-0.5, 0.5},
			q:       []float64{0.5, 0.5},
			wantErr: true,
		},
		{
			name:    "dimension mismatch",
			p:       []float64{0.5, 0.5},
			q:       []float64{0.5},
			wantErr: true,
		},
		{
			name:    "empty input",
			p:       []float64{},
			q:       []float64{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := KLDivergence(tt.p, tt.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr {
				if tt.isInf {
					if !math.IsInf(result, 1) {
						t.Errorf("expected infinite divergence, got %v", result)
					}
				} else if math.Abs(result-tt.expected) > 0.01 {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

//nolint:dupl // Test structure duplication is acceptable and improves test clarity
func TestJensenShannonDivergence(t *testing.T) {
	tests := []struct {
		name     string
		p, q     []float64
		expected float64
		wantErr  bool
	}{
		{
			name:     "identical distributions",
			p:        []float64{0.5, 0.5},
			q:        []float64{0.5, 0.5},
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "different distributions",
			p:        []float64{0.9, 0.1},
			q:        []float64{0.1, 0.9},
			expected: 0.368,
			wantErr:  false,
		},
		{
			name:    "negative values",
			p:       []float64{-0.5, 0.5},
			q:       []float64{0.5, 0.5},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := JensenShannonDivergence(tt.p, tt.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

//nolint:dupl // Test structure duplication is acceptable and improves test clarity
func TestBhattacharyya(t *testing.T) {
	tests := []struct {
		name     string
		p, q     []float64
		expected float64
		wantErr  bool
	}{
		{
			name:     "identical distributions",
			p:        []float64{0.5, 0.5},
			q:        []float64{0.5, 0.5},
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "similar distributions",
			p:        []float64{0.6, 0.4},
			q:        []float64{0.5, 0.5},
			expected: 0.005,
			wantErr:  false,
		},
		{
			name:    "negative values",
			p:       []float64{-0.5, 0.5},
			q:       []float64{0.5, 0.5},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Bhattacharyya(tt.p, tt.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHellinger(t *testing.T) {
	tests := []struct {
		name     string
		p, q     []float64
		expected float64
		wantErr  bool
	}{
		{
			name:     "identical distributions",
			p:        []float64{0.5, 0.5},
			q:        []float64{0.5, 0.5},
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "completely different",
			p:        []float64{1.0, 0.0},
			q:        []float64{0.0, 1.0},
			expected: 1.0,
			wantErr:  false,
		},
		{
			name:     "similar distributions",
			p:        []float64{0.6, 0.4},
			q:        []float64{0.5, 0.5},
			expected: 0.05,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Hellinger(tt.p, tt.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && math.Abs(result-tt.expected) > 0.1 {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestChiSquare(t *testing.T) {
	tests := []struct {
		name     string
		p, q     []float64
		expected float64
		wantErr  bool
	}{
		{
			name:     "identical distributions",
			p:        []float64{0.5, 0.5},
			q:        []float64{0.5, 0.5},
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "different distributions",
			p:        []float64{0.4, 0.6},
			q:        []float64{0.6, 0.4},
			expected: 0.08,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ChiSquare(tt.p, tt.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTotalVariation(t *testing.T) {
	tests := []struct {
		name     string
		p, q     []float64
		expected float64
		wantErr  bool
	}{
		{
			name:     "identical distributions",
			p:        []float64{0.5, 0.5},
			q:        []float64{0.5, 0.5},
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "completely different",
			p:        []float64{1.0, 0.0},
			q:        []float64{0.0, 1.0},
			expected: 1.0,
			wantErr:  false,
		},
		{
			name:     "similar distributions",
			p:        []float64{0.6, 0.4},
			q:        []float64{0.5, 0.5},
			expected: 0.1,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TotalVariation(tt.p, tt.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && !almostEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCrossEntropy(t *testing.T) {
	tests := []struct {
		name     string
		p, q     []float64
		expected float64
		wantErr  bool
		isInf    bool
	}{
		{
			name:     "uniform distributions",
			p:        []float64{0.5, 0.5},
			q:        []float64{0.5, 0.5},
			expected: 0.693,
			wantErr:  false,
		},
		{
			name:    "zero in q causes infinite entropy",
			p:       []float64{0.5, 0.5},
			q:       []float64{1.0, 0.0},
			wantErr: false,
			isInf:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CrossEntropy(tt.p, tt.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr {
				if tt.isInf {
					if !math.IsInf(result, 1) {
						t.Errorf("expected infinite entropy, got %v", result)
					}
				} else if math.Abs(result-tt.expected) > 0.01 {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

//nolint:dupl // Test structure duplication is acceptable and improves test clarity
func TestPearsonCorrelation(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float64
		expected float64
		wantErr  bool
	}{
		{
			name:     "perfect positive correlation",
			a:        []float64{1, 2, 3, 4, 5},
			b:        []float64{2, 4, 6, 8, 10},
			expected: 1.0,
			wantErr:  false,
		},
		{
			name:     "perfect negative correlation",
			a:        []float64{1, 2, 3, 4, 5},
			b:        []float64{10, 8, 6, 4, 2},
			expected: -1.0,
			wantErr:  false,
		},
		{
			name:     "no correlation",
			a:        []float64{1, 2, 3, 4, 5},
			b:        []float64{3, 3, 3, 3, 3},
			expected: 0.0,
			wantErr:  true, // zero variance
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PearsonCorrelation(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && !almostEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestPearsonDistance(t *testing.T) {
	a := []float64{1, 2, 3, 4, 5}
	b := []float64{2, 4, 6, 8, 10}

	result, err := PearsonDistance(a, b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Perfect correlation (r=1) should give distance 0
	expected := 0.0
	if !almostEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

//nolint:dupl // Test structure duplication is acceptable and improves test clarity
func TestSpearmanCorrelation(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float64
		expected float64
		wantErr  bool
	}{
		{
			name:     "perfect monotonic positive",
			a:        []float64{1, 2, 3, 4, 5},
			b:        []float64{2, 4, 6, 8, 10},
			expected: 1.0,
			wantErr:  false,
		},
		{
			name:     "perfect monotonic negative",
			a:        []float64{1, 2, 3, 4, 5},
			b:        []float64{10, 8, 6, 4, 2},
			expected: -1.0,
			wantErr:  false,
		},
		{
			name:     "non-linear monotonic",
			a:        []float64{1, 2, 3, 4, 5},
			b:        []float64{1, 4, 9, 16, 25},
			expected: 1.0,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SpearmanCorrelation(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && !almostEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestWasserstein1D(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float64
		expected float64
		wantErr  bool
	}{
		{
			name:     "identical distributions",
			a:        []float64{1, 2, 3},
			b:        []float64{1, 2, 3},
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "shifted distribution",
			a:        []float64{1, 2, 3},
			b:        []float64{2, 3, 4},
			expected: 1.0,
			wantErr:  false,
		},
		{
			name:    "empty input",
			a:       []float64{},
			b:       []float64{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Wasserstein1D(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && !almostEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestComputeRanks(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected []float64
	}{
		{
			name:     "no ties",
			values:   []float64{3, 1, 4, 2},
			expected: []float64{3, 1, 4, 2},
		},
		{
			name:     "with ties",
			values:   []float64{1, 2, 2, 3},
			expected: []float64{1, 2.5, 2.5, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := computeRanks(tt.values)
			for i := range result {
				if !almostEqual(result[i], tt.expected[i]) {
					t.Errorf("at index %d: expected %v, got %v", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

// Benchmarks
func BenchmarkKLDivergence(b *testing.B) {
	p := []float64{0.1, 0.2, 0.3, 0.4}
	q := []float64{0.2, 0.2, 0.3, 0.3}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = KLDivergence(p, q)
	}
}

func BenchmarkJensenShannonDivergence(b *testing.B) {
	p := []float64{0.1, 0.2, 0.3, 0.4}
	q := []float64{0.2, 0.2, 0.3, 0.3}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = JensenShannonDivergence(p, q)
	}
}

func BenchmarkPearsonCorrelation(b *testing.B) {
	a := make([]float64, 1000)
	b2 := make([]float64, 1000)
	for i := range a {
		a[i] = float64(i)
		b2[i] = float64(i) * 1.5
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PearsonCorrelation(a, b2)
	}
}

func BenchmarkSpearmanCorrelation(b *testing.B) {
	a := make([]float64, 100)
	b2 := make([]float64, 100)
	for i := range a {
		a[i] = float64(i)
		b2[i] = float64(i) * float64(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = SpearmanCorrelation(a, b2)
	}
}
