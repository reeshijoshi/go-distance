# go-distance

[![Go Reference](https://pkg.go.dev/badge/github.com/reeshijoshi/go-distance.svg)](https://pkg.go.dev/github.com/reeshijoshi/go-distance)
[![Go Report Card](https://goreportcard.com/badge/github.com/reeshijoshi/go-distance)](https://goreportcard.com/report/github.com/reeshijoshi/go-distance)
[![Coverage](https://img.shields.io/badge/coverage-80%25-brightgreen)](https://github.com/reeshijoshi/go-distance)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive, production-ready Go library providing 100+ distance, similarity, and divergence metrics for vectors, strings, probability distributions, time series, graphs, and geographic coordinates.

## Why go-distance?

- **Comprehensive**: 100+ metrics across 8 categories - the most complete distance library in Go
- **Type-Safe**: Leverages Go 1.18+ generics for compile-time type safety
- **Zero Dependencies**: Pure Go implementation with no external dependencies
- **Production-Ready**: 80% test coverage, extensively benchmarked, battle-tested algorithms
- **High Performance**: Optimized implementations with O(1) space complexity where possible
- **Well-Documented**: Every function includes complexity analysis and usage examples
- **Versatile**: Use cases spanning ML, NLP, GIS, bioinformatics, data science, and optimization

## Installation

```bash
go get github.com/reeshijoshi/go-distance
```

**Requirements:** Go 1.21 or later

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/reeshijoshi/go-distance"
)

func main() {
    // Vector distances
    euclidean, _ := distance.Euclidean([]float64{1, 2, 3}, []float64{4, 5, 6})
    fmt.Printf("Euclidean: %.2f\n", euclidean) // 5.20

    // String distances
    levenshtein, _ := distance.Levenshtein("kitten", "sitting")
    fmt.Printf("Levenshtein: %d\n", levenshtein) // 3

    // Statistical divergence
    kl, _ := distance.KLDivergence([]float64{0.5, 0.5}, []float64{0.6, 0.4})
    fmt.Printf("KL Divergence: %.4f\n", kl) // 0.0202

    // Geographic distance
    nyc := distance.Coord{Lat: 40.7128, Lon: -74.0060}
    london := distance.Coord{Lat: 51.5074, Lon: -0.1278}
    haversine := distance.Haversine(nyc, london)
    fmt.Printf("NYC to London: %.0f km\n", haversine) // 5570 km
}
```

## Feature Categories

### 🔢 Vector Distances
For numerical data, feature vectors, embeddings

| Function | Description | Use Case |
|----------|-------------|----------|
| `Euclidean` | L2 norm, straight-line distance | General-purpose, ML features |
| `Manhattan` | L1 norm, taxicab distance | Grid-based problems, sparse data |
| `Chebyshev` | L∞ norm, maximum difference | Chess moves, uniform grids |
| `Minkowski` | Generalized Lp norm | Flexible distance metric |
| `Cosine` | Angular distance | Text embeddings, recommendations |
| `Canberra` | Weighted Manhattan | Positive data, outlier-sensitive |
| `BrayCurtis` | Ecological dissimilarity | Compositional data |
| `Hamming` | Differing positions | Binary vectors, error detection |
| `WeightedEuclidean` | Weighted L2 norm | Feature importance weighting |

**Example:**
```go
// K-nearest neighbors with custom metric
vectors := [][]float64{{1, 2}, {3, 4}, {5, 6}}
neighbors, _ := distance.KNearestNeighbors(vectors, 2, distance.Euclidean[float64])
```

### 📝 String Distances
For text comparison, spell checking, fuzzy matching

| Function | Description | Use Case |
|----------|-------------|----------|
| `Levenshtein` | Edit distance (insert/delete/substitute) | Spell checking, DNA alignment |
| `DamerauLevenshtein` | Edit distance with transpositions | OCR correction, typo detection |
| `Jaro` | String similarity | Record linkage, name matching |
| `JaroWinkler` | Jaro with prefix bonus | Names, addresses |
| `HammingString` | Substitution-only distance | Fixed-length codes |
| `LongestCommonSubsequence` | LCS length | Diff algorithms, plagiarism |
| `NGramDistance` | N-gram overlap | Language detection, fuzzy search |
| `Soundex` | Phonetic encoding | Name search, homophones |
| `TokenSortRatio` | Word-order invariant | Document similarity |

**Example:**
```go
// Fuzzy string matching for search
query := "color"
candidates := []string{"colour", "cooler", "collar", "colored"}
for _, candidate := range candidates {
    dist, _ := distance.Levenshtein(query, candidate)
    if dist <= 2 {
        fmt.Printf("Match: %s (distance: %d)\n", candidate, dist)
    }
}
```

### 📊 Statistical Distances
For probability distributions, histograms, data analysis

| Function | Description | Use Case |
|----------|-------------|----------|
| `KLDivergence` | Kullback-Leibler divergence | Information theory, ML loss |
| `JensenShannonDivergence` | Symmetric KL divergence | Distribution comparison |
| `Bhattacharyya` | Distribution overlap | Classification, image processing |
| `Hellinger` | Distribution distance | Density estimation |
| `ChiSquare` | χ² distance | Histogram comparison |
| `TotalVariation` | L1 on probabilities | Statistical testing |
| `CrossEntropy` | Information-theoretic loss | ML training |
| `PearsonCorrelation` | Linear correlation | Regression analysis |
| `SpearmanCorrelation` | Rank correlation | Non-parametric statistics |
| `Wasserstein1D` | Earth mover's distance | Optimal transport |

**Example:**
```go
// Compare distributions
observed := []float64{0.25, 0.35, 0.20, 0.20}
expected := []float64{0.25, 0.25, 0.25, 0.25}

kl, _ := distance.KLDivergence(observed, expected)
js, _ := distance.JensenShannonDivergence(observed, expected)
fmt.Printf("KL: %.4f, JS: %.4f\n", kl, js)
```

### 🔤 Set-Based Distances
For comparing sets, bags, collections

| Function | Description | Use Case |
|----------|-------------|----------|
| `JaccardSet` | Intersection over union | Document similarity |
| `DiceSorensen` | 2*intersection / (size1+size2) | Image segmentation |
| `OverlapCoefficient` | Subset similarity | Query matching |
| `TanimotoCoefficient` | Generalized Jaccard | Chemical fingerprints |
| `CosineSimilaritySet` | Bag-of-words cosine | Text retrieval |

**Example:**
```go
// Find similar documents by tags
doc1 := []string{"go", "programming", "tutorial"}
doc2 := []string{"go", "golang", "programming"}

jaccard, _ := distance.JaccardSet(doc1, doc2)
dice, _ := distance.DiceSorensen(doc1, doc2)
fmt.Printf("Jaccard: %.2f, Dice: %.2f\n", 1-jaccard, dice)
```

### 🌍 Geographic Distances
For location data, GIS, mapping applications

| Function | Description | Accuracy | Speed |
|----------|-------------|----------|-------|
| `Haversine` | Great circle distance | ±0.5% | Fast |
| `Vincenty` | Geodesic on ellipsoid (WGS-84) | ±0.5mm | Moderate |
| `GreatCircle` | Spherical law of cosines | ±0.5% | Fast |
| `Equirectangular` | Fast approximation | ±1% (short) | Very Fast |

**Example:**
```go
// Calculate delivery distance
warehouse := distance.Coord{Lat: 37.7749, Lon: -122.4194} // SF
customer := distance.Coord{Lat: 34.0522, Lon: -118.2437}  // LA

km := distance.Haversine(warehouse, customer)
miles := distance.HaversineMiles(warehouse, customer)
fmt.Printf("Distance: %.0f km (%.0f miles)\n", km, miles)
```

### ⏱️ Time Series Distances
For temporal data, signal processing, pattern matching

| Function | Description | Use Case |
|----------|-------------|----------|
| `DTW` | Dynamic Time Warping | Speech recognition, gesture matching |
| `DTWWithWindow` | DTW with Sakoe-Chiba band | Faster DTW with constraints |
| `Frechet` | Discrete Fréchet distance | Curve similarity, GPS tracks |
| `Hausdorff` | Maximum point-to-set distance | Shape matching, image comparison |
| `SmithWaterman` | Local sequence alignment | DNA/protein sequence analysis |
| `NeedlemanWunsch` | Global sequence alignment | Bioinformatics |
| `Autocorrelation` | Self-similarity at lag | Seasonality detection |

**Example:**
```go
// Find similar patterns in sensor data
signal1 := []float64{1, 2, 3, 4, 5}
signal2 := []float64{0, 1, 2, 3, 4, 5, 6}

dtwDist, _ := distance.DTW(signal1, signal2)
fmt.Printf("DTW Distance: %.2f\n", dtwDist)
```

### 📈 Graph Distances
For network analysis, shortest paths, graph theory

| Function | Description | Complexity |
|----------|-------------|------------|
| `Dijkstra` | Single-source shortest path | O((V+E)log V) |
| `BellmanFord` | Handles negative weights | O(VE) |
| `FloydWarshall` | All-pairs shortest paths | O(V³) |
| `AStar` | Heuristic search | O(E log V) |
| `BFS` | Unweighted shortest path | O(V+E) |
| `GraphDiameter` | Maximum shortest path | O(V³) |
| `GraphEditDistance` | Graph similarity | Approximate |

**Example:**
```go
// Route planning
g := distance.NewGraph()
g.AddEdge(0, 1, 10.0) // Node 0 to 1, weight 10
g.AddEdge(1, 2, 5.0)
g.AddEdge(0, 2, 25.0)

dist, path := g.Dijkstra(0, 2)
fmt.Printf("Shortest path: %v, distance: %.0f\n", path, dist)
```

### 🎯 Optimization Algorithms
For parameter tuning, function minimization, search

| Algorithm | Type | Use Case |
|-----------|------|----------|
| `GradientDescent` | First-order | Convex optimization |
| `Adam` | Adaptive moment | ML training |
| `BFGS` | Quasi-Newton | Smooth non-convex |
| `NelderMead` | Derivative-free | Black-box optimization |
| `SimulatedAnnealing` | Stochastic | Combinatorial problems |
| `GeneticAlgorithm` | Evolutionary | Multi-modal search |
| `ParticleSwarmOptimization` | Swarm intelligence | Continuous optimization |

**Example:**
```go
// Minimize a function
f := func(x []float64) float64 {
    return x[0]*x[0] + x[1]*x[1] // f(x,y) = x² + y²
}

grad := func(x []float64) []float64 {
    return []float64{2*x[0], 2*x[1]}
}

minimum := distance.GradientDescent(f, grad, []float64{5, 5}, 0.1, 100)
fmt.Printf("Minimum at: %v\n", minimum) // Near [0, 0]
```

## Advanced Features

### Batch Operations
Efficient parallel computation for large datasets:

```go
// Compute distance matrix for all pairs
vectors := [][]float64{{1, 2}, {3, 4}, {5, 6}, {7, 8}}
matrix, _ := distance.BatchComputeParallel(vectors, distance.Euclidean[float64], 4)
```

### Context-Aware Computation
Cancel long-running operations:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

matrix, _ := distance.BatchComputeWithContext(ctx, vectors, distFunc, 8)
```

## Performance Benchmarks

Benchmarks run on Apple M1, Go 1.21:

```
BenchmarkEuclidean-8         100000000    10.2 ns/op     0 B/op   0 allocs/op
BenchmarkManhattan-8         200000000     8.5 ns/op     0 B/op   0 allocs/op
BenchmarkCosine-8             50000000    24.3 ns/op     0 B/op   0 allocs/op
BenchmarkLevenshtein-8         1000000  1243.0 ns/op   512 B/op   2 allocs/op
BenchmarkHaversine-8          10000000   112.0 ns/op     0 B/op   0 allocs/op
BenchmarkDTW-8                    5000 245000.0 ns/op  8192 B/op   2 allocs/op
```

## Real-World Use Cases

### Machine Learning
```go
// K-means clustering
func assignClusters(points, centroids [][]float64) []int {
    assignments := make([]int, len(points))
    for i, point := range points {
        idx, _, _ := distance.NearestNeighbor(centroids, point, distance.Euclidean[float64])
        assignments[i] = idx
    }
    return assignments
}
```

### Information Retrieval
```go
// Document ranking by relevance
func rankDocuments(query string, docs []string) []struct{ idx int; score float64 } {
    scores := make([]struct{ idx int; score float64 }, len(docs))
    for i, doc := range docs {
        sim, _ := distance.CosineSimilarityStrings(query, doc)
        scores[i] = struct{ idx int; score float64 }{i, sim}
    }
    // Sort by score descending...
    return scores
}
```

### GIS/Location Services
```go
// Find nearby locations
func findNearby(user distance.Coord, places []distance.Coord, radiusKm float64) []int {
    var nearby []int
    for i, place := range places {
        if distance.Haversine(user, place) <= radiusKm {
            nearby = append(nearby, i)
        }
    }
    return nearby
}
```

## API Design Philosophy

1. **Consistency**: All distance functions follow `func(a, b Type) (float64, error)` pattern
2. **Type Safety**: Generic functions provide compile-time guarantees
3. **Error Handling**: Explicit error returns for invalid inputs
4. **Performance**: Zero-allocation hot paths, O(1) space complexity where possible
5. **Documentation**: Every function includes time/space complexity and use cases

## Testing & Quality

- **80% Test Coverage**: Comprehensive test suite with edge cases
- **Extensive Benchmarks**: Performance regression tracking
- **No External Dependencies**: Reduced supply chain risk
- **Semantic Versioning**: Stable public API

## Contributing

Contributions are welcome! Please:

1. Open an issue to discuss major changes
2. Add tests for new functionality
3. Update documentation
4. Follow existing code style

## License

MIT License - see [LICENSE](LICENSE) for details

## Citation

If you use go-distance in academic work, please cite:

```bibtex
@software{go_distance,
  title = {go-distance: Comprehensive distance metrics library for Go},
  author = {go-distance contributors},
  year = {2025},
  url = {https://github.com/reeshijoshi/go-distance}
}
```

## Acknowledgments

Algorithms implemented from:
- Levenshtein (1966) - Binary codes capable of correcting deletions
- Damerau (1964) - A technique for computer detection and correction of spelling errors
- Vincenty (1975) - Direct and inverse solutions of geodesics on the ellipsoid
- Sakoe & Chiba (1978) - Dynamic programming algorithm optimization for spoken word recognition

## Support

- 📖 Documentation: [pkg.go.dev](https://pkg.go.dev/github.com/reeshijoshi/go-distance)
- 🐛 Issues: [GitHub Issues](https://github.com/reeshijoshi/go-distance/issues)

---

Made with ❤️ for the Go community | Star ⭐ if you find this useful!
