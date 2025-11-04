// Package distance provides comprehensive distance, similarity, and divergence metrics
// for vectors, strings, sets, probability distributions, and time series.
package distance

import (
	"context"
	"errors"
)

var (
	// ErrDimensionMismatch is returned when vector dimensions don't match.
	ErrDimensionMismatch = errors.New("dimension mismatch between vectors")

	// ErrEmptyInput is returned when input is empty.
	ErrEmptyInput = errors.New("empty input provided")

	// ErrInvalidParameter is returned when a parameter value is invalid.
	ErrInvalidParameter = errors.New("invalid parameter value")

	// ErrZeroVector is returned when a zero vector is encountered.
	ErrZeroVector = errors.New("zero vector encountered")

	// ErrNegativeValue is returned when a negative value is found in input that requires non-negative values.
	ErrNegativeValue = errors.New("negative value in input")
)

// Number constraint for generic numeric types
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Float constraint for floating point types
type Float interface {
	~float32 | ~float64
}

// DistanceFunc is a function that computes distance between two vectors.
// Note: This name stutters with the package name. Consider using batch.Func or similar
// patterns in client code to avoid repetition.
//
//nolint:revive // Name stuttering is acceptable here for API clarity and consistency
type DistanceFunc[T Number] func(a, b []T) (float64, error)

// StringDistanceFunc computes distance between strings
type StringDistanceFunc func(a, b string) (int, error)

// Options for configurable distance calculations
type Options struct {
	Normalize   bool      // Normalize result to [0,1]
	Weights     []float64 // Dimension weights
	Parallel    bool      // Use parallel computation for batch operations
	MaxDistance float64   // Early termination threshold (0 means no limit)
}

// Metric interface for any distance metric
type Metric interface {
	Name() string
	Distance(a, b any) (float64, error)
	IsSymmetric() bool // Some metrics like KL divergence are asymmetric
	IsMetric() bool    // Satisfies metric space axioms (triangle inequality, etc.)
}

// BatchComputer for efficient batch distance calculations
type BatchComputer[T Number] interface {
	ComputePairwise(vectors [][]T) ([][]float64, error)
	ComputeToPoint(vectors [][]T, point []T) ([]float64, error)
	ComputeWithContext(ctx context.Context, a, b []T) (float64, error)
}

// Validate checks if two vectors have the same dimension
func Validate[T Number](a, b []T) error {
	if len(a) == 0 || len(b) == 0 {
		return ErrEmptyInput
	}
	if len(a) != len(b) {
		return ErrDimensionMismatch
	}
	return nil
}

// ValidateWeights checks if weights match vector dimensions
func ValidateWeights[T Number](v []T, weights []float64) error {
	if len(weights) == 0 {
		return nil // No weights is valid
	}
	if len(v) != len(weights) {
		return ErrInvalidParameter
	}
	for _, w := range weights {
		if w < 0 {
			return ErrNegativeValue
		}
	}
	return nil
}
