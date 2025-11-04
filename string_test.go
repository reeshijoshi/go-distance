package distance

import "testing"

func TestLevenshtein(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected int
	}{
		{"identical", "hello", "hello", 0},
		{"classic example", "kitten", "sitting", 3},
		{"empty strings", "", "", 0},
		{"one empty", "hello", "", 5},
		{"single char diff", "cat", "bat", 1},
		{"complete difference", "abc", "xyz", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := Levenshtein(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestDamerauLevenshtein(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected int
	}{
		{"identical", "hello", "hello", 0},
		{"transposition", "ab", "ba", 1},
		{"transposition in word", "ca", "ac", 1},
		{"swap last two", "abc", "acb", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := DamerauLevenshtein(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestJaro(t *testing.T) {
	tests := []struct {
		name   string
		a, b   string
		minSim float64 // Minimum expected similarity
	}{
		{"identical", "hello", "hello", 0.99},
		{"empty", "", "", 1.0},
		{"similar", "martha", "marhta", 0.94},
		{"different", "abc", "xyz", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := Jaro(tt.a, tt.b)
			if result < tt.minSim {
				t.Errorf("expected at least %v, got %v", tt.minSim, result)
			}
		})
	}
}

func TestJaroWinkler(t *testing.T) {
	tests := []struct {
		name   string
		a, b   string
		minSim float64
	}{
		{"identical", "hello", "hello", 0.99},
		{"common prefix", "dixon", "dickson", 0.8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := JaroWinkler(tt.a, tt.b, 0.1)
			if result < tt.minSim {
				t.Errorf("expected at least %v, got %v", tt.minSim, result)
			}
		})
	}
}

func TestHammingString(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected int
		wantErr  bool
	}{
		{"identical", "hello", "hello", 0, false},
		{"one diff", "hello", "hallo", 1, false},
		{"all diff", "abc", "xyz", 3, false},
		{"length mismatch", "ab", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HammingString(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestLongestCommonSubsequence(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected int
	}{
		{"identical", "hello", "hello", 5},
		{"partial match", "abcdef", "ace", 3},
		{"no match", "abc", "xyz", 0},
		{"empty", "", "hello", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := LongestCommonSubsequence(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestLCSDistance(t *testing.T) {
	result, _ := LCSDistance("kitten", "sitting")
	// LCS("kitten", "sitting") = "ittn" = 4
	// Distance = 6 + 7 - 2*4 = 5
	expected := 5
	if result != expected {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

func TestNGramDistance(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		n    int
	}{
		{"identical", "hello", "hello", 2},
		{"similar", "hello", "hallo", 2},
		{"different", "abc", "xyz", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NGramDistance(tt.a, tt.b, tt.n)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result < 0 || result > 1 {
				t.Errorf("result out of range [0,1]: %v", result)
			}
		})
	}
}

// Benchmarks
func BenchmarkLevenshtein(b *testing.B) {
	s1 := "kitten"
	s2 := "sitting"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Levenshtein(s1, s2)
	}
}

func BenchmarkJaro(b *testing.B) {
	s1 := "martha"
	s2 := "marhta"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Jaro(s1, s2)
	}
}

func BenchmarkNGramDistance(b *testing.B) {
	s1 := "hello world"
	s2 := "hallo world"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NGramDistance(s1, s2, 2)
	}
}
