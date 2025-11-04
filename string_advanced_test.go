package distance

import (
	"testing"
)

func TestSorensenDice(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected float64
	}{
		{"identical", "hello", "hello", 1.0},
		{"similar", "night", "nacht", 0.25},
		{"different", "abc", "xyz", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := SorensenDice(tt.a, tt.b)
			if !almostEqualTolerance(result, tt.expected, 0.1) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCosineSimilarityStrings(t *testing.T) {
	result, _ := CosineSimilarityStrings("hello", "hello")
	if !almostEqual(result, 1.0) {
		t.Errorf("identical strings should have similarity 1.0, got %v", result)
	}
}

func TestRatcliffObershelp(t *testing.T) {
	result, _ := RatcliffObershelp("hello", "hallo")
	if result < 0 || result > 1 {
		t.Errorf("result should be in [0,1], got %v", result)
	}
}

func TestEditDistance(t *testing.T) {
	result, _ := EditDistance("kitten", "sitting", 1, 1, 1)
	if result <= 0 {
		t.Errorf("expected positive edit distance")
	}
}

func TestSoundex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Robert", "Robert", "R163"},
		{"Rupert", "Rupert", "R163"},
		{"empty", "", "0000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Soundex(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMetaphone(t *testing.T) {
	result := Metaphone("hello")
	if len(result) == 0 {
		t.Errorf("expected non-empty result")
	}
}

func TestPhoneticDistance(t *testing.T) {
	result := PhoneticDistance("Robert", "Rupert", Soundex)
	if result != 0 {
		t.Errorf("Robert and Rupert should have same soundex, distance should be 0, got %d", result)
	}
}

func TestNormalizedLevenshtein(t *testing.T) {
	result, _ := NormalizedLevenshtein("hello", "hello")
	if !almostEqual(result, 0.0) {
		t.Errorf("identical strings should have distance 0, got %v", result)
	}
}

func TestLCSRatio(t *testing.T) {
	result, _ := LCSRatio("hello", "hello")
	if !almostEqual(result, 1.0) {
		t.Errorf("identical strings should have LCS ratio 1.0, got %v", result)
	}
}

func TestTokenSortRatio(t *testing.T) {
	result, _ := TokenSortRatio("fuzzy wuzzy was a bear", "wuzzy fuzzy was a bear")
	if result < 0.9 {
		t.Errorf("expected high similarity for reordered words, got %v", result)
	}
}

func TestTokenSetRatio(t *testing.T) {
	result, _ := TokenSetRatio("fuzzy was a bear", "fuzzy fuzzy was a bear")
	if result < 0.5 {
		t.Errorf("expected reasonable similarity, got %v", result)
	}
}

func TestQGramDistance(t *testing.T) {
	result, err := QGramDistance("hello", "hallo", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result < 0 {
		t.Errorf("distance should be non-negative, got %d", result)
	}
}

func TestJaccardIndex(t *testing.T) {
	result, _ := JaccardIndex("hello", "hello", 2)
	if !almostEqual(result, 1.0) {
		t.Errorf("identical strings should have Jaccard index 1.0, got %v", result)
	}
}

func TestTverskyIndex(t *testing.T) {
	result, _ := TverskyIndex("hello", "hallo", 0.5, 0.5)
	if result < 0 || result > 1 {
		t.Errorf("result should be in [0,1], got %v", result)
	}
}

func BenchmarkSoundex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Soundex("Washington")
	}
}

func BenchmarkTokenSortRatio(b *testing.B) {
	s1 := "fuzzy wuzzy was a bear"
	s2 := "wuzzy fuzzy was a bear"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = TokenSortRatio(s1, s2)
	}
}
