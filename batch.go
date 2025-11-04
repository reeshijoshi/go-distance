package distance

import (
	"context"
	"sync"
)

// BatchCompute computes distances between all pairs of vectors (distance matrix).
// Time: O(n²d), Space: O(n²) where n=vectors, d=dimensions
func BatchCompute[T Number](vectors [][]T, distFn DistanceFunc[T]) ([][]float64, error) {
	n := len(vectors)
	if n == 0 {
		return [][]float64{}, nil
	}

	result := make([][]float64, n)
	for i := range result {
		result[i] = make([]float64, n)
	}

	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			dist, err := distFn(vectors[i], vectors[j])
			if err != nil {
				return nil, err
			}
			result[i][j] = dist
			result[j][i] = dist // Symmetric
		}
	}

	return result, nil
}

// BatchComputeParallel computes distance matrix in parallel.
// Time: O(n²d/workers), Space: O(n²)
func BatchComputeParallel[T Number](vectors [][]T, distFn DistanceFunc[T], workers int) ([][]float64, error) {
	n := len(vectors)
	if n == 0 {
		return [][]float64{}, nil
	}
	if workers <= 0 {
		workers = 4
	}

	result := make([][]float64, n)
	for i := range result {
		result[i] = make([]float64, n)
	}

	type job struct {
		i, j int
	}

	jobs := make(chan job, n*n/2)
	errors := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(workers)

	// Start workers
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for j := range jobs {
				dist, err := distFn(vectors[j.i], vectors[j.j])
				if err != nil {
					select {
					case errors <- err:
					default:
					}
					return
				}
				result[j.i][j.j] = dist
				result[j.j][j.i] = dist
			}
		}()
	}

	// Send jobs
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			jobs <- job{i, j}
		}
	}
	close(jobs)

	wg.Wait()
	close(errors)

	// Check for errors
	if err := <-errors; err != nil {
		return nil, err
	}

	return result, nil
}

// KNearestNeighbors finds k nearest neighbors for each vector.
// Returns indices of k nearest neighbors for each vector.
// Time: O(n²d), Space: O(nk)
func KNearestNeighbors[T Number](vectors [][]T, k int, distFn DistanceFunc[T]) ([][]int, error) {
	n := len(vectors)
	if n == 0 || k <= 0 {
		return [][]int{}, nil
	}
	if k > n-1 {
		k = n - 1
	}

	result := make([][]int, n)

	for i := 0; i < n; i++ {
		// Compute distances to all other vectors
		distances := make([]struct {
			index int
			dist  float64
		}, 0, n-1)

		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			dist, err := distFn(vectors[i], vectors[j])
			if err != nil {
				return nil, err
			}
			distances = append(distances, struct {
				index int
				dist  float64
			}{j, dist})
		}

		// Sort by distance (simple bubble sort for k elements)
		for p := 0; p < k; p++ {
			minIdx := p
			for q := p + 1; q < len(distances); q++ {
				if distances[q].dist < distances[minIdx].dist {
					minIdx = q
				}
			}
			distances[p], distances[minIdx] = distances[minIdx], distances[p]
		}

		// Extract k nearest
		result[i] = make([]int, k)
		for p := 0; p < k; p++ {
			result[i][p] = distances[p].index
		}
	}

	return result, nil
}

// RadiusNeighbors finds all neighbors within radius for each vector.
// Time: O(n²d), Space: O(n*m) where m=avg neighbors
func RadiusNeighbors[T Number](vectors [][]T, radius float64, distFn DistanceFunc[T]) ([][]int, error) {
	n := len(vectors)
	if n == 0 || radius < 0 {
		return [][]int{}, nil
	}

	result := make([][]int, n)

	for i := 0; i < n; i++ {
		neighbors := make([]int, 0)
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			dist, err := distFn(vectors[i], vectors[j])
			if err != nil {
				return nil, err
			}
			if dist <= radius {
				neighbors = append(neighbors, j)
			}
		}
		result[i] = neighbors
	}

	return result, nil
}

// ComputeToPoint computes distances from all vectors to a single point.
// Time: O(nd), Space: O(n)
func ComputeToPoint[T Number](vectors [][]T, point []T, distFn DistanceFunc[T]) ([]float64, error) {
	n := len(vectors)
	if n == 0 {
		return []float64{}, nil
	}

	result := make([]float64, n)
	for i := 0; i < n; i++ {
		dist, err := distFn(vectors[i], point)
		if err != nil {
			return nil, err
		}
		result[i] = dist
	}

	return result, nil
}

// ComputeWithContext computes distance with cancellation support.
func ComputeWithContext[T Number](ctx context.Context, a, b []T, distFn DistanceFunc[T]) (float64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return distFn(a, b)
	}
}

// BatchComputeWithContext computes distance matrix with cancellation.
func BatchComputeWithContext[T Number](ctx context.Context, vectors [][]T, distFn DistanceFunc[T], workers int) ([][]float64, error) {
	n := len(vectors)
	if n == 0 {
		return [][]float64{}, nil
	}
	if workers <= 0 {
		workers = 4
	}

	result := make([][]float64, n)
	for i := range result {
		result[i] = make([]float64, n)
	}

	type job struct {
		i, j int
	}

	jobs := make(chan job, n*n/2)
	errors := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(workers)

	// Start workers
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for j := range jobs {
				select {
				case <-ctx.Done():
					select {
					case errors <- ctx.Err():
					default:
					}
					return
				default:
					dist, err := distFn(vectors[j.i], vectors[j.j])
					if err != nil {
						select {
						case errors <- err:
						default:
						}
						return
					}
					result[j.i][j.j] = dist
					result[j.j][j.i] = dist
				}
			}
		}()
	}

	// Send jobs with cancellation check
	go func() {
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				select {
				case <-ctx.Done():
					close(jobs)
					return
				case jobs <- job{i, j}:
				}
			}
		}
		close(jobs)
	}()

	wg.Wait()
	close(errors)

	// Check for errors
	if err := <-errors; err != nil {
		return nil, err
	}

	return result, nil
}

// NearestNeighbor finds the nearest neighbor for a query point.
// Time: O(nd), Space: O(1)
func NearestNeighbor[T Number](vectors [][]T, query []T, distFn DistanceFunc[T]) (int, float64, error) {
	if len(vectors) == 0 {
		return -1, 0, ErrEmptyInput
	}

	minIdx := 0
	minDist, err := distFn(vectors[0], query)
	if err != nil {
		return -1, 0, err
	}

	for i := 1; i < len(vectors); i++ {
		dist, err := distFn(vectors[i], query)
		if err != nil {
			return -1, 0, err
		}
		if dist < minDist {
			minDist = dist
			minIdx = i
		}
	}

	return minIdx, minDist, nil
}

// Centroid computes the centroid (mean) of a set of vectors.
// Time: O(nd), Space: O(d)
func Centroid[T Number](vectors [][]T) ([]float64, error) {
	if len(vectors) == 0 {
		return nil, ErrEmptyInput
	}

	d := len(vectors[0])
	centroid := make([]float64, d)

	for _, vec := range vectors {
		if len(vec) != d {
			return nil, ErrDimensionMismatch
		}
		for i, v := range vec {
			centroid[i] += float64(v)
		}
	}

	n := float64(len(vectors))
	for i := range centroid {
		centroid[i] /= n
	}

	return centroid, nil
}

// PairwiseDistinctCount counts pairs with distance above threshold.
// Useful for diversity metrics.
// Time: O(n²d), Space: O(1)
func PairwiseDistinctCount[T Number](vectors [][]T, threshold float64, distFn DistanceFunc[T]) (int, error) {
	if len(vectors) == 0 {
		return 0, nil
	}

	count := 0
	n := len(vectors)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dist, err := distFn(vectors[i], vectors[j])
			if err != nil {
				return 0, err
			}
			if dist > threshold {
				count++
			}
		}
	}

	return count, nil
}
