package distance

import (
	"math"
	"sort"
	"strings"
	"unicode"
)

// SorensenDice computes Sørensen-Dice coefficient for strings
// Uses bigrams (2-grams) for comparison
// Range [0, 1] where 1=identical
// Time: O(n+m), Space: O(n+m)
func SorensenDice(a, b string) (float64, error) {
	if len(a) == 0 && len(b) == 0 {
		return 1.0, nil
	}
	if len(a) == 0 || len(b) == 0 {
		return 0.0, nil
	}

	bigramsA := extractBigrams(a)
	bigramsB := extractBigrams(b)

	if len(bigramsA) == 0 && len(bigramsB) == 0 {
		return 1.0, nil
	}

	intersection := 0
	for bigram := range bigramsA {
		if bigramsB[bigram] > 0 {
			intersection += min(bigramsA[bigram], bigramsB[bigram])
		}
	}

	return float64(2*intersection) / float64(len(bigramsA)+len(bigramsB)), nil
}

func extractBigrams(s string) map[string]int {
	bigrams := make(map[string]int)
	runes := []rune(s)
	for i := 0; i < len(runes)-1; i++ {
		bigram := string(runes[i : i+2])
		bigrams[bigram]++
	}
	return bigrams
}

// CosineSimilarityStrings computes cosine similarity for strings
// Treats strings as character frequency vectors
// Range [0, 1] where 1=identical distribution
// Time: O(n+m), Space: O(n+m)
func CosineSimilarityStrings(a, b string) (float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0, nil
	}

	freqA := make(map[rune]int)
	freqB := make(map[rune]int)

	for _, r := range a {
		freqA[r]++
	}
	for _, r := range b {
		freqB[r]++
	}

	var dotProduct, normA, normB float64

	// Dot product
	for char, countA := range freqA {
		if countB, exists := freqB[char]; exists {
			dotProduct += float64(countA * countB)
		}
	}

	// Norms
	for _, count := range freqA {
		normA += float64(count * count)
	}
	for _, count := range freqB {
		normB += float64(count * count)
	}

	if normA == 0 || normB == 0 {
		return 0, nil
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB)), nil
}

// RatcliffObershelp computes Ratcliff/Obershelp similarity
// Also known as Gestalt Pattern Matching
// Range [0, 1] where 1=identical
// Time: O(n²), Space: O(n)
func RatcliffObershelp(a, b string) (float64, error) {
	if len(a) == 0 && len(b) == 0 {
		return 1.0, nil
	}
	if len(a) == 0 || len(b) == 0 {
		return 0.0, nil
	}

	matches := ratcliffMatches([]rune(a), []rune(b))
	return float64(2*matches) / float64(len(a)+len(b)), nil
}

func ratcliffMatches(a, b []rune) int {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}

	// Find longest common substring
	maxLen := 0
	maxPosA, maxPosB := 0, 0

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			k := 0
			for i+k < len(a) && j+k < len(b) && a[i+k] == b[j+k] {
				k++
			}
			if k > maxLen {
				maxLen = k
				maxPosA = i
				maxPosB = j
			}
		}
	}

	if maxLen == 0 {
		return 0
	}

	// Recursively find matches before and after
	leftMatches := ratcliffMatches(a[:maxPosA], b[:maxPosB])
	rightMatches := ratcliffMatches(a[maxPosA+maxLen:], b[maxPosB+maxLen:])

	return maxLen + leftMatches + rightMatches
}

// EditDistance computes generic edit distance with custom costs
// Time: O(mn), Space: O(min(m,n))
func EditDistance(a, b string, insertCost, deleteCost, replaceCost int) (int, error) {
	if len(a) == 0 {
		return len(b) * insertCost, nil
	}
	if len(b) == 0 {
		return len(a) * deleteCost, nil
	}

	if len(a) > len(b) {
		a, b = b, a
		insertCost, deleteCost = deleteCost, insertCost
	}

	prev := make([]int, len(a)+1)
	curr := make([]int, len(a)+1)

	for i := range prev {
		prev[i] = i * deleteCost
	}

	for j := 1; j <= len(b); j++ {
		curr[0] = j * insertCost
		for i := 1; i <= len(a); i++ {
			cost := replaceCost
			if a[i-1] == b[j-1] {
				cost = 0
			}

			curr[i] = min3(
				prev[i]+insertCost,
				curr[i-1]+deleteCost,
				prev[i-1]+cost,
			)
		}
		prev, curr = curr, prev
	}

	return prev[len(a)], nil
}

// SmithWatermanString computes Smith-Waterman local alignment for strings
// Returns alignment score
// Time: O(mn), Space: O(mn)
func SmithWatermanString(a, b string, match, mismatch, gap int) (int, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0, ErrEmptyInput
	}

	m, n := len(a), len(b)
	H := make([][]int, m+1)
	for i := range H {
		H[i] = make([]int, n+1)
	}

	maxScore := 0

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			matchScore := mismatch
			if a[i-1] == b[j-1] {
				matchScore = match
			}

			H[i][j] = max(
				0,
				max(
					H[i-1][j-1]+matchScore,
					max(H[i-1][j]+gap, H[i][j-1]+gap),
				),
			)

			if H[i][j] > maxScore {
				maxScore = H[i][j]
			}
		}
	}

	return maxScore, nil
}

// MongeElkan computes Monge-Elkan similarity
// Uses maximum token similarity
// Range [0, 1] where 1=identical
// Time: O(n²m), Space: O(n)
func MongeElkan(a, b string, tokenSim func(string, string) float64) (float64, error) {
	tokensA := strings.Fields(a)
	tokensB := strings.Fields(b)

	if len(tokensA) == 0 && len(tokensB) == 0 {
		return 1.0, nil
	}
	if len(tokensA) == 0 || len(tokensB) == 0 {
		return 0.0, nil
	}

	sum := 0.0
	for _, tokenA := range tokensA {
		maxSim := 0.0
		for _, tokenB := range tokensB {
			sim := tokenSim(tokenA, tokenB)
			if sim > maxSim {
				maxSim = sim
			}
		}
		sum += maxSim
	}

	return sum / float64(len(tokensA)), nil
}

// QGramDistance computes q-gram distance
// Time: O(n+m), Space: O(n+m)
func QGramDistance(a, b string, q int) (int, error) {
	if q <= 0 {
		return 0, ErrInvalidParameter
	}

	gramsA := extractQGrams(a, q)
	gramsB := extractQGrams(b, q)

	distance := 0
	seen := make(map[string]bool)

	for gram, countA := range gramsA {
		seen[gram] = true
		countB := gramsB[gram]
		diff := countA - countB
		if diff < 0 {
			diff = -diff
		}
		distance += diff
	}

	for gram, countB := range gramsB {
		if !seen[gram] {
			distance += countB
		}
	}

	return distance, nil
}

func extractQGrams(s string, q int) map[string]int {
	grams := make(map[string]int)
	padded := strings.Repeat(" ", q-1) + s + strings.Repeat(" ", q-1)
	runes := []rune(padded)

	for i := 0; i <= len(runes)-q; i++ {
		gram := string(runes[i : i+q])
		grams[gram]++
	}

	return grams
}

// JaccardIndex computes Jaccard index for strings (using n-grams)
// Range [0, 1] where 1=identical
// Time: O(n+m), Space: O(n+m)
func JaccardIndex(a, b string, n int) (float64, error) {
	if n <= 0 {
		return 0, ErrInvalidParameter
	}

	gramsA := extractNGrams(a, n)
	gramsB := extractNGrams(b, n)

	if len(gramsA) == 0 && len(gramsB) == 0 {
		return 1.0, nil
	}

	intersection := 0
	union := 0

	all := make(map[string]bool)
	for gram := range gramsA {
		all[gram] = true
	}
	for gram := range gramsB {
		all[gram] = true
	}

	for gram := range all {
		countA, countB := gramsA[gram], gramsB[gram]
		if countA > 0 && countB > 0 {
			intersection += min(countA, countB)
		}
		if countA > countB {
			union += countA
		} else {
			union += countB
		}
	}

	if union == 0 {
		return 0, nil
	}

	return float64(intersection) / float64(union), nil
}

// TverskyIndex computes Tversky index (asymmetric Jaccard)
// alpha and beta control asymmetry
// Time: O(n+m), Space: O(n+m)
func TverskyIndex(a, b string, alpha, beta float64) (float64, error) {
	if alpha < 0 || beta < 0 {
		return 0, ErrInvalidParameter
	}

	gramsA := extractBigrams(a)
	gramsB := extractBigrams(b)

	intersection := 0
	aMinusB := 0
	bMinusA := 0

	all := make(map[string]bool)
	for gram := range gramsA {
		all[gram] = true
	}
	for gram := range gramsB {
		all[gram] = true
	}

	for gram := range all {
		countA, countB := gramsA[gram], gramsB[gram]
		if countA > 0 && countB > 0 {
			intersection += min(countA, countB)
		}
		if countA > countB {
			aMinusB += countA - countB
		}
		if countB > countA {
			bMinusA += countB - countA
		}
	}

	denom := float64(intersection) + alpha*float64(aMinusB) + beta*float64(bMinusA)
	if denom == 0 {
		return 0, nil
	}

	return float64(intersection) / denom, nil
}

// Metaphone computes metaphone phonetic encoding
// Returns phonetic code for phonetic matching
// Time: O(n), Space: O(n)
//
//nolint:gocyclo // Phonetic algorithms are inherently complex with many rules
func Metaphone(s string) string {
	if len(s) == 0 {
		return ""
	}

	s = strings.ToUpper(s)
	result := ""

	// Simplified metaphone (not full algorithm)
	for i, r := range s {
		switch r {
		case 'A', 'E', 'I', 'O', 'U':
			if i == 0 {
				result += string(r)
			}
		case 'B':
			if i == len(s)-1 && i > 0 && s[i-1] == 'M' {
				continue
			}
			result += "B"
		case 'C':
			//nolint:gocritic // Metaphone algorithm requires character lookahead, if-else is most readable
			if i+1 < len(s) && s[i+1] == 'H' {
				result += "X"
			} else if i+1 < len(s) && (s[i+1] == 'I' || s[i+1] == 'E' || s[i+1] == 'Y') {
				result += "S"
			} else {
				result += "K"
			}
		case 'D':
			result += "T"
		case 'G':
			if i+1 < len(s) && (s[i+1] == 'H' || s[i+1] == 'N') {
				continue
			}
			result += "K"
		case 'H':
			if i == 0 || unicode.IsLetter(rune(s[i-1])) {
				result += "H"
			}
		case 'K':
			if i == 0 || s[i-1] != 'C' {
				result += "K"
			}
		case 'P':
			if i+1 < len(s) && s[i+1] == 'H' {
				result += "F"
			} else {
				result += "P"
			}
		case 'Q':
			result += "K"
		case 'S':
			if i+1 < len(s) && s[i+1] == 'H' {
				result += "X"
			} else {
				result += "S"
			}
		case 'T':
			//nolint:gocritic // Metaphone algorithm requires character lookahead, if-else is most readable
			if i+2 < len(s) && s[i+1] == 'I' && (s[i+2] == 'O' || s[i+2] == 'A') {
				result += "X"
			} else if i+1 < len(s) && s[i+1] == 'H' {
				result += "0"
			} else {
				result += "T"
			}
		case 'V':
			result += "F"
		case 'W', 'Y':
			if i+1 < len(s) && unicode.IsLetter(rune(s[i+1])) {
				result += string(r)
			}
		case 'X':
			result += "KS"
		case 'Z':
			result += "S"
		default:
			if unicode.IsLetter(r) {
				result += string(r)
			}
		}
	}

	return result
}

// Soundex computes Soundex phonetic encoding
// Returns 4-character code for phonetic matching
// Time: O(n), Space: O(1)
func Soundex(s string) string {
	if len(s) == 0 {
		return "0000"
	}

	s = strings.ToUpper(s)
	result := string(s[0])
	prev := soundexCode(rune(s[0]))

	for i := 1; i < len(s) && len(result) < 4; i++ {
		code := soundexCode(rune(s[i]))
		if code != "0" && code != prev {
			result += code
			prev = code
		} else if code == "0" {
			prev = "0"
		}
	}

	// Pad with zeros
	for len(result) < 4 {
		result += "0"
	}

	return result
}

func soundexCode(r rune) string {
	switch r {
	case 'B', 'F', 'P', 'V':
		return "1"
	case 'C', 'G', 'J', 'K', 'Q', 'S', 'X', 'Z':
		return "2"
	case 'D', 'T':
		return "3"
	case 'L':
		return "4"
	case 'M', 'N':
		return "5"
	case 'R':
		return "6"
	default:
		return "0"
	}
}

// PhoneticDistance computes distance between phonetic encodings
// Returns 0 if phonetically similar
// Time: O(1), Space: O(1)
func PhoneticDistance(a, b string, encoder func(string) string) int {
	codeA := encoder(a)
	codeB := encoder(b)

	if codeA == codeB {
		return 0
	}

	// Count differences
	distance := 0
	maxLen := len(codeA)
	if len(codeB) > maxLen {
		maxLen = len(codeB)
	}

	for i := 0; i < maxLen; i++ {
		if i >= len(codeA) || i >= len(codeB) || codeA[i] != codeB[i] {
			distance++
		}
	}

	return distance
}

// LCSRatio computes LCS-based similarity ratio
// Range [0, 1] where 1=identical
// Time: O(mn), Space: O(min(m,n))
func LCSRatio(a, b string) (float64, error) {
	lcs, err := LongestCommonSubsequence(a, b)
	if err != nil {
		return 0, err
	}

	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}

	if maxLen == 0 {
		return 1.0, nil
	}

	return float64(lcs) / float64(maxLen), nil
}

// NormalizedLevenshtein computes normalized Levenshtein distance
// Range [0, 1] where 0=identical
// Time: O(mn), Space: O(min(m,n))
func NormalizedLevenshtein(a, b string) (float64, error) {
	dist, err := Levenshtein(a, b)
	if err != nil {
		return 0, err
	}

	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}

	if maxLen == 0 {
		return 0, nil
	}

	return float64(dist) / float64(maxLen), nil
}

// TokenSortRatio computes similarity after sorting tokens
// Useful for comparing similar strings with different word order
// Range [0, 1] where 1=identical
// Time: O(n log n), Space: O(n)
func TokenSortRatio(a, b string) (float64, error) {
	tokensA := strings.Fields(strings.ToLower(a))
	tokensB := strings.Fields(strings.ToLower(b))

	sort.Strings(tokensA)
	sort.Strings(tokensB)

	sortedA := strings.Join(tokensA, " ")
	sortedB := strings.Join(tokensB, " ")

	dist, err := Levenshtein(sortedA, sortedB)
	if err != nil {
		return 0, err
	}

	maxLen := len(sortedA)
	if len(sortedB) > maxLen {
		maxLen = len(sortedB)
	}

	if maxLen == 0 {
		return 1.0, nil
	}

	return 1.0 - float64(dist)/float64(maxLen), nil
}

// TokenSetRatio computes similarity using set intersection of tokens
// Range [0, 1] where 1=identical
// Time: O(n), Space: O(n)
func TokenSetRatio(a, b string) (float64, error) {
	tokensA := strings.Fields(strings.ToLower(a))
	tokensB := strings.Fields(strings.ToLower(b))

	setA := make(map[string]bool)
	setB := make(map[string]bool)

	for _, t := range tokensA {
		setA[t] = true
	}
	for _, t := range tokensB {
		setB[t] = true
	}

	intersection := []string{}
	for t := range setA {
		if setB[t] {
			intersection = append(intersection, t)
		}
	}

	sort.Strings(intersection)
	intersectionStr := strings.Join(intersection, " ")

	diff1 := []string{}
	for t := range setA {
		if !setB[t] {
			diff1 = append(diff1, t)
		}
	}
	sort.Strings(diff1)

	diff2 := []string{}
	for t := range setB {
		if !setA[t] {
			diff2 = append(diff2, t)
		}
	}
	sort.Strings(diff2)

	combined1 := intersectionStr + " " + strings.Join(diff1, " ")
	combined2 := intersectionStr + " " + strings.Join(diff2, " ")

	dist, err := Levenshtein(combined1, combined2)
	if err != nil {
		return 0, err
	}

	maxLen := len(combined1)
	if len(combined2) > maxLen {
		maxLen = len(combined2)
	}

	if maxLen == 0 {
		return 1.0, nil
	}

	return 1.0 - float64(dist)/float64(maxLen), nil
}
