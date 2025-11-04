package distance

import "math"

// DTW computes Dynamic Time Warping distance between two time series.
// Allows matching sequences of different lengths.
// Time: O(mn), Space: O(min(m,n)) with optimization
func DTW[T Number](a, b []T) (float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0, ErrEmptyInput
	}

	// Ensure a is shorter for space optimization
	if len(a) > len(b) {
		a, b = b, a
	}

	n, m := len(a), len(b)

	// Use two rows instead of full matrix
	prev := make([]float64, n+1)
	curr := make([]float64, n+1)

	// Initialize first row
	prev[0] = 0
	for i := 1; i <= n; i++ {
		prev[i] = math.Inf(1)
	}

	// Fill matrix
	for j := 1; j <= m; j++ {
		curr[0] = math.Inf(1)
		for i := 1; i <= n; i++ {
			cost := math.Abs(float64(a[i-1]) - float64(b[j-1]))
			curr[i] = cost + math.Min(
				math.Min(prev[i], curr[i-1]),
				prev[i-1],
			)
		}
		prev, curr = curr, prev
	}

	return prev[n], nil
}

// DTWWithWindow computes DTW with Sakoe-Chiba band constraint.
// Window limits how far sequences can deviate (improves performance).
// Time: O(mn), Space: O(min(m,n))
func DTWWithWindow[T Number](a, b []T, window int) (float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0, ErrEmptyInput
	}
	if window < 0 {
		return 0, ErrInvalidParameter
	}

	n, m := len(a), len(b)
	prev := make([]float64, n+1)
	curr := make([]float64, n+1)

	// Initialize
	for i := range prev {
		prev[i] = math.Inf(1)
	}
	prev[0] = 0

	// Fill with window constraint
	for j := 1; j <= m; j++ {
		for i := range curr {
			curr[i] = math.Inf(1)
		}

		// Sakoe-Chiba band
		start := max(1, j-window)
		end := min(n, j+window)

		for i := start; i <= end; i++ {
			cost := math.Abs(float64(a[i-1]) - float64(b[j-1]))
			curr[i] = cost + math.Min(
				math.Min(prev[i], curr[i-1]),
				prev[i-1],
			)
		}
		prev, curr = curr, prev
	}

	return prev[n], nil
}

// Frechet computes discrete Fréchet distance between two curves.
// Measures similarity considering the flow of the curves.
// Time: O(mn), Space: O(mn)
func Frechet[T Number](a, b [][]T) (float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0, ErrEmptyInput
	}

	n, m := len(a), len(b)
	ca := make([][]float64, n)
	for i := range ca {
		ca[i] = make([]float64, m)
		for j := range ca[i] {
			ca[i][j] = -1
		}
	}

	var frechetRecursive func(i, j int) float64
	frechetRecursive = func(i, j int) float64 {
		if ca[i][j] > -1 {
			return ca[i][j]
		}

		// Euclidean distance between points
		dist, _ := Euclidean(a[i], b[j])

		//nolint:gocritic // Frechet algorithm requires boundary condition checks, if-else is most readable
		if i == 0 && j == 0 {
			ca[i][j] = dist
		} else if i > 0 && j == 0 {
			ca[i][j] = math.Max(frechetRecursive(i-1, 0), dist)
		} else if i == 0 && j > 0 {
			ca[i][j] = math.Max(frechetRecursive(0, j-1), dist)
		} else {
			ca[i][j] = math.Max(
				math.Min(
					math.Min(frechetRecursive(i-1, j), frechetRecursive(i-1, j-1)),
					frechetRecursive(i, j-1),
				),
				dist,
			)
		}
		return ca[i][j]
	}

	return frechetRecursive(n-1, m-1), nil
}

// Hausdorff computes Hausdorff distance between two point sets.
// Measures maximum distance from a point in one set to the nearest point in the other.
// Time: O(nm), Space: O(1)
func Hausdorff[T Number](a, b [][]T) (float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0, ErrEmptyInput
	}

	directedHausdorff := func(from, to [][]T) float64 {
		maxMin := 0.0
		for _, p1 := range from {
			minDist := math.Inf(1)
			for _, p2 := range to {
				dist, _ := Euclidean(p1, p2)
				if dist < minDist {
					minDist = dist
				}
			}
			if minDist > maxMin {
				maxMin = minDist
			}
		}
		return maxMin
	}

	// Hausdorff is the maximum of both directions
	h1 := directedHausdorff(a, b)
	h2 := directedHausdorff(b, a)

	return math.Max(h1, h2), nil
}

// LongestCommonSubstring computes longest common substring length for sequences.
// Time: O(mn), Space: O(mn)
func LongestCommonSubstring[T comparable](a, b []T) int {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}

	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	maxLen := 0
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				if dp[i][j] > maxLen {
					maxLen = dp[i][j]
				}
			}
		}
	}

	return maxLen
}

// SmithWaterman computes local sequence alignment score.
// Used for DNA/protein sequence comparison.
// Time: O(mn), Space: O(mn)
func SmithWaterman[T comparable](a, b []T, match, mismatch, gap int) (int, error) {
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
			var matchScore int
			if a[i-1] == b[j-1] {
				matchScore = match
			} else {
				matchScore = mismatch
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

// NeedlemanWunsch computes global sequence alignment score.
// Time: O(mn), Space: O(mn)
func NeedlemanWunsch[T comparable](a, b []T, match, mismatch, gap int) (int, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0, ErrEmptyInput
	}

	m, n := len(a), len(b)
	F := make([][]int, m+1)
	for i := range F {
		F[i] = make([]int, n+1)
	}

	// Initialize
	for i := 0; i <= m; i++ {
		F[i][0] = i * gap
	}
	for j := 0; j <= n; j++ {
		F[0][j] = j * gap
	}

	// Fill matrix
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			var matchScore int
			if a[i-1] == b[j-1] {
				matchScore = match
			} else {
				matchScore = mismatch
			}

			F[i][j] = max(
				F[i-1][j-1]+matchScore,
				max(F[i-1][j]+gap, F[i][j-1]+gap),
			)
		}
	}

	return F[m][n], nil
}

// SoftDTW computes differentiable DTW using soft-min.
// Useful for machine learning applications.
// gamma controls smoothness (smaller = closer to DTW).
// Time: O(mn), Space: O(mn)
func SoftDTW[T Number](a, b []T, gamma float64) (float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0, ErrEmptyInput
	}
	if gamma <= 0 {
		return 0, ErrInvalidParameter
	}

	n, m := len(a), len(b)
	R := make([][]float64, n+1)
	for i := range R {
		R[i] = make([]float64, m+1)
	}

	// Initialize
	for i := 1; i <= n; i++ {
		R[i][0] = math.Inf(1)
	}
	for j := 1; j <= m; j++ {
		R[0][j] = math.Inf(1)
	}
	R[0][0] = 0

	// Soft-min function
	softMin := func(a, b, c float64) float64 {
		return -gamma * math.Log(
			math.Exp(-a/gamma)+
				math.Exp(-b/gamma)+
				math.Exp(-c/gamma),
		)
	}

	// Fill matrix
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cost := math.Pow(float64(a[i-1])-float64(b[j-1]), 2)
			R[i][j] = cost + softMin(R[i-1][j], R[i][j-1], R[i-1][j-1])
		}
	}

	return R[n][m], nil
}

// Autocorrelation computes autocorrelation at lag k.
// Measures correlation of a signal with a delayed copy of itself.
// Time: O(n), Space: O(1)
func Autocorrelation[T Number](data []T, lag int) (float64, error) {
	if len(data) == 0 {
		return 0, ErrEmptyInput
	}
	if lag < 0 || lag >= len(data) {
		return 0, ErrInvalidParameter
	}

	n := len(data)
	mean := 0.0
	for _, v := range data {
		mean += float64(v)
	}
	mean /= float64(n)

	var c0, ck float64
	for i := 0; i < n; i++ {
		diff := float64(data[i]) - mean
		c0 += diff * diff
		if i+lag < n {
			ck += diff * (float64(data[i+lag]) - mean)
		}
	}

	if c0 == 0 {
		return 0, nil
	}

	return ck / c0, nil
}
