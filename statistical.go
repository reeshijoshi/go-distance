package distance

import (
	"math"
	"sort"
)

// KLDivergence computes Kullback-Leibler divergence KL(P||Q).
// Measures how probability distribution P diverges from Q.
// NOTE: Asymmetric (KL(P||Q) ≠ KL(Q||P))
// Time: O(n), Space: O(1)
func KLDivergence[T Float](p, q []T) (float64, error) {
	if err := Validate(p, q); err != nil {
		return 0, err
	}

	var divergence float64
	for i := range p {
		pi, qi := float64(p[i]), float64(q[i])

		if pi < 0 || qi < 0 {
			return 0, ErrNegativeValue
		}

		if pi > 0 {
			if qi == 0 {
				return math.Inf(1), nil // Infinite divergence
			}
			divergence += pi * math.Log(pi/qi)
		}
	}
	return divergence, nil
}

// JensenShannonDivergence computes Jensen-Shannon divergence.
// Symmetric version of KL divergence: JS(P||Q) = JS(Q||P)
// Bounded: 0 ≤ JS ≤ log(2)
// Time: O(n), Space: O(n)
func JensenShannonDivergence[T Float](p, q []T) (float64, error) {
	if err := Validate(p, q); err != nil {
		return 0, err
	}

	// Compute average distribution M = (P+Q)/2
	m := make([]float64, len(p))
	for i := range p {
		pi, qi := float64(p[i]), float64(q[i])
		if pi < 0 || qi < 0 {
			return 0, ErrNegativeValue
		}
		m[i] = (pi + qi) / 2.0
	}

	// JS(P||Q) = (KL(P||M) + KL(Q||M)) / 2
	var klPM, klQM float64
	for i := range p {
		pi, qi := float64(p[i]), float64(q[i])
		mi := m[i]

		if pi > 0 && mi > 0 {
			klPM += pi * math.Log(pi/mi)
		}
		if qi > 0 && mi > 0 {
			klQM += qi * math.Log(qi/mi)
		}
	}

	return (klPM + klQM) / 2.0, nil
}

// Bhattacharyya computes Bhattacharyya distance.
// Measures similarity between probability distributions.
// Range [0, +∞) where 0=identical
// Time: O(n), Space: O(1)
func Bhattacharyya[T Float](p, q []T) (float64, error) {
	if err := Validate(p, q); err != nil {
		return 0, err
	}

	var bc float64 // Bhattacharyya coefficient
	for i := range p {
		pi, qi := float64(p[i]), float64(q[i])
		if pi < 0 || qi < 0 {
			return 0, ErrNegativeValue
		}
		bc += math.Sqrt(pi * qi)
	}

	if bc == 0 {
		return math.Inf(1), nil
	}
	if bc > 1 {
		bc = 1 // Clamp for numerical stability
	}

	return -math.Log(bc), nil
}

// Hellinger computes Hellinger distance.
// Related to Bhattacharyya: H² = 1 - BC
// Range [0, 1] where 0=identical, 1=completely different
// Time: O(n), Space: O(1)
func Hellinger[T Float](p, q []T) (float64, error) {
	if err := Validate(p, q); err != nil {
		return 0, err
	}

	var sum float64
	for i := range p {
		pi, qi := float64(p[i]), float64(q[i])
		if pi < 0 || qi < 0 {
			return 0, ErrNegativeValue
		}
		diff := math.Sqrt(pi) - math.Sqrt(qi)
		sum += diff * diff
	}

	return math.Sqrt(sum) / math.Sqrt(2.0), nil
}

// ChiSquare computes Chi-square distance.
// Used for histogram comparison in computer vision.
// Time: O(n), Space: O(1)
func ChiSquare[T Float](p, q []T) (float64, error) {
	if err := Validate(p, q); err != nil {
		return 0, err
	}

	var sum float64
	for i := range p {
		pi, qi := float64(p[i]), float64(q[i])
		if pi < 0 || qi < 0 {
			return 0, ErrNegativeValue
		}

		denom := pi + qi
		if denom != 0 {
			num := pi - qi
			sum += (num * num) / denom
		}
	}

	return sum, nil
}

// TotalVariation computes total variation distance.
// Maximum difference in probabilities across all events.
// Range [0, 1]
// Time: O(n), Space: O(1)
func TotalVariation[T Float](p, q []T) (float64, error) {
	if err := Validate(p, q); err != nil {
		return 0, err
	}

	var sum float64
	for i := range p {
		pi, qi := float64(p[i]), float64(q[i])
		if pi < 0 || qi < 0 {
			return 0, ErrNegativeValue
		}
		sum += math.Abs(pi - qi)
	}

	return sum / 2.0, nil
}

// CrossEntropy computes cross-entropy H(P,Q).
// Used in ML loss functions: H(P,Q) = -Σ p(x) log q(x)
// Time: O(n), Space: O(1)
func CrossEntropy[T Float](p, q []T) (float64, error) {
	if err := Validate(p, q); err != nil {
		return 0, err
	}

	var entropy float64
	for i := range p {
		pi, qi := float64(p[i]), float64(q[i])
		if pi < 0 || qi < 0 {
			return 0, ErrNegativeValue
		}
		if pi > 0 {
			if qi == 0 {
				return math.Inf(1), nil
			}
			entropy -= pi * math.Log(qi)
		}
	}

	return entropy, nil
}

// PearsonCorrelation computes Pearson correlation coefficient.
// Range [-1, 1] where 1=perfect positive, -1=perfect negative, 0=no correlation
// Time: O(n), Space: O(1)
func PearsonCorrelation[T Number](a, b []T) (float64, error) {
	if err := Validate(a, b); err != nil {
		return 0, err
	}

	n := float64(len(a))

	// Compute means
	var sumA, sumB float64
	for i := range a {
		sumA += float64(a[i])
		sumB += float64(b[i])
	}
	meanA, meanB := sumA/n, sumB/n

	// Compute correlation
	var numerator, varA, varB float64
	for i := range a {
		diffA := float64(a[i]) - meanA
		diffB := float64(b[i]) - meanB
		numerator += diffA * diffB
		varA += diffA * diffA
		varB += diffB * diffB
	}

	if varA == 0 || varB == 0 {
		return 0, ErrZeroVector
	}

	return numerator / math.Sqrt(varA*varB), nil
}

// PearsonDistance computes distance based on Pearson correlation.
// Range [0, 2] where 0=perfect correlation, 2=perfect anti-correlation
// Time: O(n), Space: O(1)
func PearsonDistance[T Number](a, b []T) (float64, error) {
	corr, err := PearsonCorrelation(a, b)
	if err != nil {
		return 0, err
	}
	return 1 - corr, nil
}

// SpearmanCorrelation computes Spearman rank correlation.
// Measures monotonic relationship between variables.
// Time: O(n log n), Space: O(n)
func SpearmanCorrelation[T Number](a, b []T) (float64, error) {
	if err := Validate(a, b); err != nil {
		return 0, err
	}

	// Convert to ranks
	ranksA := computeRanks(a)
	ranksB := computeRanks(b)

	// Compute Pearson correlation on ranks
	return PearsonCorrelation(ranksA, ranksB)
}

// computeRanks converts values to ranks (average rank for ties)
func computeRanks[T Number](values []T) []float64 {
	n := len(values)

	// Create pairs of (value, index)
	type pair struct {
		val float64
		idx int
	}
	pairs := make([]pair, n)
	for i, v := range values {
		pairs[i] = pair{float64(v), i}
	}

	// Sort by value using standard library
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].val < pairs[j].val
	})

	ranks := make([]float64, n)
	i := 0
	for i < n {
		j := i
		// Find ties
		for j < n && pairs[j].val == pairs[i].val {
			j++
		}

		// Assign average rank
		avgRank := float64(i+j+1) / 2.0
		for k := i; k < j; k++ {
			ranks[pairs[k].idx] = avgRank
		}
		i = j
	}

	return ranks
}

// Wasserstein1D computes 1D Wasserstein (Earth Mover's) distance.
// For 1D distributions, this equals the area between CDFs.
// Time: O(n log n), Space: O(n)
func Wasserstein1D[T Number](a, b []T) (float64, error) {
	if err := Validate(a, b); err != nil {
		return 0, err
	}

	// Sort both arrays
	aSorted := make([]float64, len(a))
	bSorted := make([]float64, len(b))
	for i := range a {
		aSorted[i] = float64(a[i])
		bSorted[i] = float64(b[i])
	}

	sortFloat64Slice(aSorted)
	sortFloat64Slice(bSorted)

	// Compute area between CDFs
	var distance float64
	for i := range aSorted {
		distance += math.Abs(aSorted[i] - bSorted[i])
	}

	return distance / float64(len(a)), nil
}

// sortFloat64Slice sorts a float64 slice using standard library
func sortFloat64Slice(arr []float64) {
	sort.Float64s(arr)
}
