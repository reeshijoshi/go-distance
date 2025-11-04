package distance

// Levenshtein computes the Levenshtein edit distance between two strings.
// Counts minimum insertions, deletions, and substitutions.
// Time: O(mn), Space: O(min(m,n)) with optimization
func Levenshtein(a, b string) (int, error) {
	if len(a) == 0 {
		return len(b), nil
	}
	if len(b) == 0 {
		return len(a), nil
	}

	// Ensure a is the shorter string to optimize space
	if len(a) > len(b) {
		a, b = b, a
	}

	// Use two rows instead of full matrix
	prevRow := make([]int, len(a)+1)
	currRow := make([]int, len(a)+1)

	// Initialize first row
	for i := range prevRow {
		prevRow[i] = i
	}

	for j := 1; j <= len(b); j++ {
		currRow[0] = j
		for i := 1; i <= len(a); i++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}

			currRow[i] = min3(
				prevRow[i]+1,      // deletion
				currRow[i-1]+1,    // insertion
				prevRow[i-1]+cost, // substitution
			)
		}
		prevRow, currRow = currRow, prevRow
	}

	return prevRow[len(a)], nil
}

// DamerauLevenshtein computes Damerau-Levenshtein distance.
// Includes transposition of adjacent characters (ab -> ba).
// Time: O(mn), Space: O(mn)
func DamerauLevenshtein(a, b string) (int, error) {
	if len(a) == 0 {
		return len(b), nil
	}
	if len(b) == 0 {
		return len(a), nil
	}

	lenA, lenB := len(a), len(b)
	maxDist := lenA + lenB

	// Create distance matrix with extra row/col
	h := make([][]int, lenA+2)
	for i := range h {
		h[i] = make([]int, lenB+2)
	}

	h[0][0] = maxDist
	for i := 0; i <= lenA; i++ {
		h[i+1][0] = maxDist
		h[i+1][1] = i
	}
	for j := 0; j <= lenB; j++ {
		h[0][j+1] = maxDist
		h[1][j+1] = j
	}

	for i := 1; i <= lenA; i++ {
		for j := 1; j <= lenB; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}

			h[i+1][j+1] = min3(
				h[i][j+1]+1,  // deletion
				h[i+1][j]+1,  // insertion
				h[i][j]+cost, // substitution
			)

			// Transposition
			if i > 1 && j > 1 && a[i-1] == b[j-2] && a[i-2] == b[j-1] {
				h[i+1][j+1] = min(h[i+1][j+1], h[i-1][j-1]+1)
			}
		}
	}

	return h[lenA+1][lenB+1], nil
}

// Jaro computes the Jaro similarity between two strings.
// Returns similarity in [0, 1] where 1=identical
// Time: O(mn), Space: O(max(m,n))
func Jaro(a, b string) (float64, error) {
	if len(a) == 0 && len(b) == 0 {
		return 1.0, nil
	}
	if len(a) == 0 || len(b) == 0 {
		return 0.0, nil
	}

	matchWindow := max(len(a), len(b))/2 - 1
	if matchWindow < 0 {
		matchWindow = 0
	}

	aMatches := make([]bool, len(a))
	bMatches := make([]bool, len(b))

	matches := 0
	transpositions := 0

	// Find matches
	for i := 0; i < len(a); i++ {
		start := max(0, i-matchWindow)
		end := min(i+matchWindow+1, len(b))

		for j := start; j < end; j++ {
			if bMatches[j] || a[i] != b[j] {
				continue
			}
			aMatches[i] = true
			bMatches[j] = true
			matches++
			break
		}
	}

	if matches == 0 {
		return 0.0, nil
	}

	// Count transpositions
	k := 0
	for i := 0; i < len(a); i++ {
		if !aMatches[i] {
			continue
		}
		for !bMatches[k] {
			k++
		}
		if a[i] != b[k] {
			transpositions++
		}
		k++
	}

	m := float64(matches)
	t := float64(transpositions) / 2.0

	return (m/float64(len(a)) + m/float64(len(b)) + (m-t)/m) / 3.0, nil
}

// JaroWinkler computes Jaro-Winkler similarity (Jaro with prefix bonus).
// prefixScale: scaling factor for prefix (standard: 0.1)
// Returns similarity in [0, 1] where 1=identical
// Time: O(mn), Space: O(max(m,n))
func JaroWinkler(a, b string, prefixScale float64) (float64, error) {
	jaroSim, err := Jaro(a, b)
	if err != nil {
		return 0, err
	}

	// Find common prefix up to 4 characters
	prefixLen := 0
	for i := 0; i < min(min(len(a), len(b)), 4); i++ {
		if a[i] == b[i] {
			prefixLen++
		} else {
			break
		}
	}

	return jaroSim + float64(prefixLen)*prefixScale*(1.0-jaroSim), nil
}

// HammingString computes Hamming distance for strings (must be equal length).
// Time: O(n), Space: O(1)
func HammingString(a, b string) (int, error) {
	if len(a) != len(b) {
		return 0, ErrDimensionMismatch
	}

	count := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			count++
		}
	}
	return count, nil
}

// LongestCommonSubsequence computes the length of LCS.
// Time: O(mn), Space: O(min(m,n))
func LongestCommonSubsequence(a, b string) (int, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0, nil
	}

	// Ensure a is shorter to optimize space
	if len(a) > len(b) {
		a, b = b, a
	}

	prev := make([]int, len(a)+1)
	curr := make([]int, len(a)+1)

	for j := 1; j <= len(b); j++ {
		for i := 1; i <= len(a); i++ {
			if a[i-1] == b[j-1] {
				curr[i] = prev[i-1] + 1
			} else {
				curr[i] = max(prev[i], curr[i-1])
			}
		}
		prev, curr = curr, prev
	}

	return prev[len(a)], nil
}

// LCSDistance computes distance based on LCS.
// Returns: len(a) + len(b) - 2*LCS(a,b)
// Time: O(mn), Space: O(min(m,n))
func LCSDistance(a, b string) (int, error) {
	lcs, err := LongestCommonSubsequence(a, b)
	if err != nil {
		return 0, err
	}
	return len(a) + len(b) - 2*lcs, nil
}

// NGramDistance computes distance based on n-gram overlap.
// Time: O(m+n), Space: O(m+n)
func NGramDistance(a, b string, n int) (float64, error) {
	if n <= 0 {
		return 0, ErrInvalidParameter
	}

	if len(a) < n && len(b) < n {
		if a == b {
			return 0, nil
		}
		return 1, nil
	}

	aNgrams := extractNGrams(a, n)
	bNgrams := extractNGrams(b, n)

	if len(aNgrams) == 0 && len(bNgrams) == 0 {
		return 0, nil
	}

	intersection := 0
	for ngram := range aNgrams {
		if bNgrams[ngram] > 0 {
			intersection += min(aNgrams[ngram], bNgrams[ngram])
		}
	}

	totalNgrams := 0
	for _, count := range aNgrams {
		totalNgrams += count
	}
	for _, count := range bNgrams {
		totalNgrams += count
	}

	if totalNgrams == 0 {
		return 0, nil
	}

	return 1.0 - float64(2*intersection)/float64(totalNgrams), nil
}

// extractNGrams extracts n-grams from a string with padding
func extractNGrams(s string, n int) map[string]int {
	ngrams := make(map[string]int)

	// Pad string
	padded := string(make([]byte, n-1)) + s + string(make([]byte, n-1))

	for i := 0; i <= len(padded)-n; i++ {
		ngram := padded[i : i+n]
		ngrams[ngram]++
	}

	return ngrams
}

// Helper functions
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func min3(a, b, c int) int {
	return minInt(minInt(a, b), c)
}
