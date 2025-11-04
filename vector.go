package distance

import (
	"math"
)

// Euclidean computes the L2 norm (straight-line distance) between two vectors.
// Time: O(n), Space: O(1)
func Euclidean[T Number](a, b []T) (float64, error) {
	if err := Validate(a, b); err != nil {
		return 0, err
	}

	var sum float64
	for i := range a {
		diff := float64(a[i]) - float64(b[i])
		sum += diff * diff
	}
	return math.Sqrt(sum), nil
}

// EuclideanSquared computes squared Euclidean distance (faster, avoids sqrt).
// Time: O(n), Space: O(1)
func EuclideanSquared[T Number](a, b []T) (float64, error) {
	if err := Validate(a, b); err != nil {
		return 0, err
	}

	var sum float64
	for i := range a {
		diff := float64(a[i]) - float64(b[i])
		sum += diff * diff
	}
	return sum, nil
}

// Manhattan computes the L1 norm (sum of absolute differences).
// Also known as: Taxicab distance, City Block distance
// Time: O(n), Space: O(1)
func Manhattan[T Number](a, b []T) (float64, error) {
	if err := Validate(a, b); err != nil {
		return 0, err
	}

	var sum float64
	for i := range a {
		diff := float64(a[i]) - float64(b[i])
		if diff < 0 {
			diff = -diff
		}
		sum += diff
	}
	return sum, nil
}

// Chebyshev computes the L-infinity norm (maximum absolute difference).
// Also known as: Chessboard distance
// Time: O(n), Space: O(1)
func Chebyshev[T Number](a, b []T) (float64, error) {
	if err := Validate(a, b); err != nil {
		return 0, err
	}

	var maxDiff float64
	for i := range a {
		diff := float64(a[i]) - float64(b[i])
		if diff < 0 {
			diff = -diff
		}
		if diff > maxDiff {
			maxDiff = diff
		}
	}
	return maxDiff, nil
}

// Minkowski computes the Lp norm with parameter p.
// p=1: Manhattan, p=2: Euclidean, p=inf: Chebyshev
// Time: O(n), Space: O(1)
func Minkowski[T Number](a, b []T, p float64) (float64, error) {
	if err := Validate(a, b); err != nil {
		return 0, err
	}
	if p <= 0 {
		return 0, ErrInvalidParameter
	}
	if math.IsInf(p, 1) {
		return Chebyshev(a, b)
	}

	var sum float64
	for i := range a {
		diff := float64(a[i]) - float64(b[i])
		if diff < 0 {
			diff = -diff
		}
		sum += math.Pow(diff, p)
	}
	return math.Pow(sum, 1/p), nil
}

// Cosine computes the cosine distance (1 - cosine similarity).
// Measures angle between vectors, range [0, 2] (0=identical direction, 2=opposite)
// Time: O(n), Space: O(1)
func Cosine[T Number](a, b []T) (float64, error) {
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

	if normA == 0 || normB == 0 {
		return 0, ErrZeroVector
	}

	similarity := dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
	// Clamp to [-1, 1] to handle floating point errors
	if similarity > 1 {
		similarity = 1
	} else if similarity < -1 {
		similarity = -1
	}

	return 1 - similarity, nil
}

// CosineSimilarity computes the cosine similarity (dot product of normalized vectors).
// Range [-1, 1] where 1=identical, 0=orthogonal, -1=opposite
// Time: O(n), Space: O(1)
func CosineSimilarity[T Number](a, b []T) (float64, error) {
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

	if normA == 0 || normB == 0 {
		return 0, ErrZeroVector
	}

	similarity := dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
	// Clamp to [-1, 1]
	if similarity > 1 {
		return 1, nil
	} else if similarity < -1 {
		return -1, nil
	}
	return similarity, nil
}

// Canberra computes the Canberra distance (weighted Manhattan distance).
// Sensitive to small changes near zero. Used in biology/ecology.
// Time: O(n), Space: O(1)
func Canberra[T Number](a, b []T) (float64, error) {
	if err := Validate(a, b); err != nil {
		return 0, err
	}

	var sum float64
	for i := range a {
		fa, fb := math.Abs(float64(a[i])), math.Abs(float64(b[i]))
		numerator := math.Abs(fa - fb)
		denominator := fa + fb

		if denominator != 0 {
			sum += numerator / denominator
		}
	}
	return sum, nil
}

// BrayCurtis computes the Bray-Curtis dissimilarity.
// Used in ecology and biology for compositional data.
// Range [0, 1] where 0=identical, 1=completely different
// Time: O(n), Space: O(1)
func BrayCurtis[T Number](a, b []T) (float64, error) {
	if err := Validate(a, b); err != nil {
		return 0, err
	}

	var numerator, denominator float64
	for i := range a {
		fa, fb := float64(a[i]), float64(b[i])
		numerator += math.Abs(fa - fb)
		denominator += math.Abs(fa + fb)
	}

	if denominator == 0 {
		return 0, nil // Both vectors are zero
	}
	return numerator / denominator, nil
}

// Hamming computes the Hamming distance (number of differing positions).
// Time: O(n), Space: O(1)
func Hamming[T Number](a, b []T) (float64, error) {
	if err := Validate(a, b); err != nil {
		return 0, err
	}

	var count float64
	for i := range a {
		if a[i] != b[i] {
			count++
		}
	}
	return count, nil
}

// WeightedEuclidean computes weighted Euclidean distance.
// weights[i] scales the contribution of dimension i.
// Time: O(n), Space: O(1)
func WeightedEuclidean[T Number](a, b []T, weights []float64) (float64, error) {
	if err := Validate(a, b); err != nil {
		return 0, err
	}
	if err := ValidateWeights(a, weights); err != nil {
		return 0, err
	}

	var sum float64
	for i := range a {
		diff := float64(a[i]) - float64(b[i])
		w := 1.0
		if len(weights) > 0 {
			w = weights[i]
		}
		sum += w * diff * diff
	}
	return math.Sqrt(sum), nil
}

// DotProduct computes the dot product (inner product) of two vectors.
// Time: O(n), Space: O(1)
func DotProduct[T Number](a, b []T) (float64, error) {
	if err := Validate(a, b); err != nil {
		return 0, err
	}

	var sum float64
	for i := range a {
		sum += float64(a[i]) * float64(b[i])
	}
	return sum, nil
}

// Norm computes the Lp norm of a vector.
// Time: O(n), Space: O(1)
func Norm[T Number](v []T, p float64) (float64, error) {
	if len(v) == 0 {
		return 0, ErrEmptyInput
	}
	if p <= 0 {
		return 0, ErrInvalidParameter
	}

	if math.IsInf(p, 1) {
		// L-infinity norm: maximum absolute value
		var maxAbs float64
		for _, val := range v {
			abs := math.Abs(float64(val))
			if abs > maxAbs {
				maxAbs = abs
			}
		}
		return maxAbs, nil
	}

	var sum float64
	for _, val := range v {
		abs := math.Abs(float64(val))
		sum += math.Pow(abs, p)
	}
	return math.Pow(sum, 1/p), nil
}
