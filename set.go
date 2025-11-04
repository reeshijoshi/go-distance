package distance

import "math"

// JaccardSet computes Jaccard distance for sets.
// Distance = 1 - |A ∩ B| / |A ∪ B|
// Range [0, 1] where 0=identical, 1=completely different
// Time: O(n+m), Space: O(n)
func JaccardSet[T comparable](a, b []T) (float64, error) {
	if len(a) == 0 && len(b) == 0 {
		return 0, nil
	}

	setA := make(map[T]bool, len(a))
	for _, item := range a {
		setA[item] = true
	}

	intersection := 0
	setB := make(map[T]bool, len(b))
	for _, item := range b {
		if setA[item] {
			intersection++
		}
		setB[item] = true
	}

	union := len(setA) + len(setB) - intersection
	if union == 0 {
		return 0, nil
	}

	return 1.0 - float64(intersection)/float64(union), nil
}

// JaccardSimilarity computes Jaccard similarity coefficient.
// Similarity = |A ∩ B| / |A ∪ B|
// Range [0, 1] where 1=identical, 0=no overlap
// Time: O(n+m), Space: O(n)
func JaccardSimilarity[T comparable](a, b []T) (float64, error) {
	dist, err := JaccardSet(a, b)
	if err != nil {
		return 0, err
	}
	return 1.0 - dist, nil
}

// DiceSorensen computes Dice-Sørensen coefficient.
// Similarity = 2|A ∩ B| / (|A| + |B|)
// Range [0, 1] where 1=identical
// Time: O(n+m), Space: O(n)
func DiceSorensen[T comparable](a, b []T) (float64, error) {
	if len(a) == 0 && len(b) == 0 {
		return 1.0, nil
	}
	if len(a) == 0 || len(b) == 0 {
		return 0.0, nil
	}

	setA := make(map[T]bool, len(a))
	for _, item := range a {
		setA[item] = true
	}

	intersection := 0
	for _, item := range b {
		if setA[item] {
			intersection++
		}
	}

	return float64(2*intersection) / float64(len(a)+len(b)), nil
}

// DiceDistance computes Dice distance (1 - Dice coefficient).
// Time: O(n+m), Space: O(n)
func DiceDistance[T comparable](a, b []T) (float64, error) {
	sim, err := DiceSorensen(a, b)
	if err != nil {
		return 0, err
	}
	return 1.0 - sim, nil
}

// OverlapCoefficient computes overlap coefficient.
// Overlap = |A ∩ B| / min(|A|, |B|)
// Range [0, 1] where 1=one is subset of the other
// Time: O(n+m), Space: O(n)
func OverlapCoefficient[T comparable](a, b []T) (float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0.0, nil
	}

	setA := make(map[T]bool, len(a))
	for _, item := range a {
		setA[item] = true
	}

	intersection := 0
	for _, item := range b {
		if setA[item] {
			intersection++
		}
	}

	minSize := len(a)
	if len(b) < minSize {
		minSize = len(b)
	}

	return float64(intersection) / float64(minSize), nil
}

// TanimotoCoefficient computes Tanimoto coefficient (generalized Jaccard).
// For binary vectors: Tanimoto = dot(a,b) / (||a||² + ||b||² - dot(a,b))
// Time: O(n), Space: O(1)
func TanimotoCoefficient[T Number](a, b []T) (float64, error) {
	if err := Validate(a, b); err != nil {
		return 0, err
	}

	var dotProduct, normA, normB float64
	for i := range a {
		fa, fb := float64(a[i]), float64(b[i])
		dotProduct += fa * fb
		normA += fa * fa
		normB += fb * fb
	}

	denominator := normA + normB - dotProduct
	if denominator == 0 {
		if dotProduct == 0 {
			return 0, nil // Both vectors are zero
		}
		return 1, nil
	}

	return dotProduct / denominator, nil
}

// TanimotoDistance computes Tanimoto distance (1 - Tanimoto coefficient).
// Time: O(n), Space: O(1)
func TanimotoDistance[T Number](a, b []T) (float64, error) {
	coef, err := TanimotoCoefficient(a, b)
	if err != nil {
		return 0, err
	}
	return 1.0 - coef, nil
}

// CosineSimilaritySet computes cosine similarity for bag-of-words.
// Uses term frequency vectors constructed from sets.
// Time: O(n+m), Space: O(n+m)
func CosineSimilaritySet[T comparable](a, b []T) (float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0, nil
	}

	// Build frequency maps
	freqA := make(map[T]int)
	freqB := make(map[T]int)

	for _, item := range a {
		freqA[item]++
	}
	for _, item := range b {
		freqB[item]++
	}

	// Compute dot product and norms
	var dotProduct, normA, normB float64

	// Dot product (only over common terms)
	for term, countA := range freqA {
		if countB, exists := freqB[term]; exists {
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

// CosineDistanceSet computes cosine distance for bag-of-words sets.
// Time: O(n+m), Space: O(n+m)
func CosineDistanceSet[T comparable](a, b []T) (float64, error) {
	sim, err := CosineSimilaritySet(a, b)
	if err != nil {
		return 0, err
	}
	return 1.0 - sim, nil
}
