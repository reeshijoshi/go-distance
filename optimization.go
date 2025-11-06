// Package distance provides distance metrics and optimization algorithms.
//
//nolint:gosec // G404: math/rand/v2 is intentionally used for optimization algorithms.
// Cryptographic randomness is not required for these mathematical optimization functions
// (simulated annealing, genetic algorithms, PSO, differential evolution).
// Using crypto/rand would be unnecessarily slow and provide no security benefit.
package distance

import (
	"math"
	"math/rand/v2"
)

// OptimizationFunc represents a function to minimize/maximize
type OptimizationFunc func([]float64) float64

// GradientFunc computes the gradient of the function
type GradientFunc func([]float64) []float64

// GradientDescent performs gradient descent optimization
// Time: O(iterations * d), Space: O(d)
func GradientDescent(
	_ OptimizationFunc, // unused but kept for API consistency
	grad GradientFunc,
	initial []float64,
	learningRate float64,
	iterations int,
) []float64 {
	x := make([]float64, len(initial))
	copy(x, initial)

	for i := 0; i < iterations; i++ {
		gradient := grad(x)
		for j := range x {
			x[j] -= learningRate * gradient[j]
		}
	}

	return x
}

// GradientDescentWithMomentum performs gradient descent with momentum
// Time: O(iterations * d), Space: O(d)
func GradientDescentWithMomentum(
	_ OptimizationFunc, // unused but kept for API consistency
	grad GradientFunc,
	initial []float64,
	learningRate float64,
	momentum float64,
	iterations int,
) []float64 {
	x := make([]float64, len(initial))
	copy(x, initial)

	velocity := make([]float64, len(initial))

	for i := 0; i < iterations; i++ {
		gradient := grad(x)
		for j := range x {
			velocity[j] = momentum*velocity[j] - learningRate*gradient[j]
			x[j] += velocity[j]
		}
	}

	return x
}

// Adam optimizer (Adaptive Moment Estimation)
// Time: O(iterations * d), Space: O(d)
func Adam(
	_ OptimizationFunc, // unused but kept for API consistency
	grad GradientFunc,
	initial []float64,
	learningRate float64,
	beta1, beta2 float64,
	epsilon float64,
	iterations int,
) []float64 {
	x := make([]float64, len(initial))
	copy(x, initial)

	m := make([]float64, len(initial)) // First moment
	v := make([]float64, len(initial)) // Second moment

	for t := 1; t <= iterations; t++ {
		gradient := grad(x)

		for j := range x {
			// Update biased first moment estimate
			m[j] = beta1*m[j] + (1-beta1)*gradient[j]
			// Update biased second moment estimate
			v[j] = beta2*v[j] + (1-beta2)*gradient[j]*gradient[j]

			// Compute bias-corrected moments
			mHat := m[j] / (1 - math.Pow(beta1, float64(t)))
			vHat := v[j] / (1 - math.Pow(beta2, float64(t)))

			// Update parameters
			x[j] -= learningRate * mHat / (math.Sqrt(vHat) + epsilon)
		}
	}

	return x
}

// SimulatedAnnealing performs simulated annealing optimization
// Time: O(iterations * d), Space: O(d)
func SimulatedAnnealing(
	f OptimizationFunc,
	initial []float64,
	initialTemp float64,
	coolingRate float64,
	iterations int,
	stepSize float64,
) []float64 {
	current := make([]float64, len(initial))
	copy(current, initial)
	currentEnergy := f(current)

	best := make([]float64, len(initial))
	copy(best, current)
	bestEnergy := currentEnergy

	temp := initialTemp

	for i := 0; i < iterations; i++ {
		// Generate neighbor solution
		neighbor := make([]float64, len(current))
		for j := range current {
			neighbor[j] = current[j] + (rand.Float64()-0.5)*2*stepSize
		}

		neighborEnergy := f(neighbor)
		delta := neighborEnergy - currentEnergy

		// Accept or reject
		if delta < 0 || rand.Float64() < math.Exp(-delta/temp) {
			copy(current, neighbor)
			currentEnergy = neighborEnergy

			if currentEnergy < bestEnergy {
				copy(best, current)
				bestEnergy = currentEnergy
			}
		}

		// Cool down
		temp *= coolingRate
	}

	return best
}

// Individual represents a genetic algorithm individual
type Individual struct {
	Genes   []float64
	Fitness float64
}

// GeneticAlgorithm performs genetic algorithm optimization
// Time: O(generations * popSize * d), Space: O(popSize * d)
func GeneticAlgorithm(
	f OptimizationFunc,
	dimensions int,
	bounds [][]float64, // [min, max] for each dimension
	popSize int,
	generations int,
	mutationRate float64,
	crossoverRate float64,
) []float64 {
	// Initialize population
	population := make([]Individual, popSize)
	for i := range population {
		genes := make([]float64, dimensions)
		for j := range genes {
			genes[j] = bounds[j][0] + rand.Float64()*(bounds[j][1]-bounds[j][0])
		}
		population[i] = Individual{
			Genes:   genes,
			Fitness: f(genes),
		}
	}

	for gen := 0; gen < generations; gen++ {
		// Selection (tournament)
		newPopulation := make([]Individual, popSize)
		for i := 0; i < popSize; i++ {
			// Tournament selection
			a := rand.IntN(popSize)
			b := rand.IntN(popSize)
			if population[a].Fitness < population[b].Fitness {
				newPopulation[i] = population[a]
			} else {
				newPopulation[i] = population[b]
			}
		}

		// Crossover
		for i := 0; i < popSize-1; i += 2 {
			if rand.Float64() < crossoverRate {
				point := rand.IntN(dimensions)
				for j := point; j < dimensions; j++ {
					newPopulation[i].Genes[j], newPopulation[i+1].Genes[j] =
						newPopulation[i+1].Genes[j], newPopulation[i].Genes[j]
				}
			}
		}

		// Mutation
		for i := range newPopulation {
			for j := range newPopulation[i].Genes {
				if rand.Float64() < mutationRate {
					newPopulation[i].Genes[j] = bounds[j][0] +
						rand.Float64()*(bounds[j][1]-bounds[j][0])
				}
			}
			newPopulation[i].Fitness = f(newPopulation[i].Genes)
		}

		population = newPopulation
	}

	// Find best
	best := population[0]
	for i := 1; i < popSize; i++ {
		if population[i].Fitness < best.Fitness {
			best = population[i]
		}
	}

	return best.Genes
}

// Particle represents a PSO particle
type Particle struct {
	Position     []float64
	Velocity     []float64
	BestPosition []float64
	BestFitness  float64
	Fitness      float64
}

// ParticleSwarmOptimization performs PSO
// Time: O(iterations * swarmSize * d), Space: O(swarmSize * d)
func ParticleSwarmOptimization(
	f OptimizationFunc,
	dimensions int,
	bounds [][]float64,
	swarmSize int,
	iterations int,
	inertia float64,
	cognitive float64,
	social float64,
) []float64 {
	// Initialize swarm
	swarm := make([]Particle, swarmSize)
	globalBest := make([]float64, dimensions)
	globalBestFitness := math.Inf(1)

	for i := range swarm {
		position := make([]float64, dimensions)
		velocity := make([]float64, dimensions)

		for j := range position {
			position[j] = bounds[j][0] + rand.Float64()*(bounds[j][1]-bounds[j][0])
			velocity[j] = (rand.Float64() - 0.5) * (bounds[j][1] - bounds[j][0])
		}

		fitness := f(position)
		swarm[i] = Particle{
			Position:     position,
			Velocity:     velocity,
			BestPosition: append([]float64{}, position...),
			BestFitness:  fitness,
			Fitness:      fitness,
		}

		if fitness < globalBestFitness {
			globalBestFitness = fitness
			copy(globalBest, position)
		}
	}

	// Iterate
	for iter := 0; iter < iterations; iter++ {
		for i := range swarm {
			for j := 0; j < dimensions; j++ {
				r1 := rand.Float64()
				r2 := rand.Float64()

				// Update velocity
				swarm[i].Velocity[j] = inertia*swarm[i].Velocity[j] +
					cognitive*r1*(swarm[i].BestPosition[j]-swarm[i].Position[j]) +
					social*r2*(globalBest[j]-swarm[i].Position[j])

				// Update position
				swarm[i].Position[j] += swarm[i].Velocity[j]

				// Clamp to bounds
				if swarm[i].Position[j] < bounds[j][0] {
					swarm[i].Position[j] = bounds[j][0]
				}
				if swarm[i].Position[j] > bounds[j][1] {
					swarm[i].Position[j] = bounds[j][1]
				}
			}

			// Evaluate fitness
			swarm[i].Fitness = f(swarm[i].Position)

			// Update personal best
			if swarm[i].Fitness < swarm[i].BestFitness {
				swarm[i].BestFitness = swarm[i].Fitness
				copy(swarm[i].BestPosition, swarm[i].Position)
			}

			// Update global best
			if swarm[i].Fitness < globalBestFitness {
				globalBestFitness = swarm[i].Fitness
				copy(globalBest, swarm[i].Position)
			}
		}
	}

	return globalBest
}

// NelderMead performs Nelder-Mead simplex optimization
// Time: O(iterations * d²), Space: O(d²)
func NelderMead(
	f OptimizationFunc,
	initial []float64,
	iterations int,
	alpha, gamma, rho, sigma float64,
) []float64 {
	n := len(initial)

	// Initialize simplex
	simplex := make([][]float64, n+1)
	values := make([]float64, n+1)

	// First vertex is the initial point
	simplex[0] = make([]float64, n)
	copy(simplex[0], initial)
	values[0] = f(simplex[0])

	// Create other vertices by perturbing initial point
	for i := 1; i <= n; i++ {
		simplex[i] = make([]float64, n)
		copy(simplex[i], initial)
		simplex[i][i-1] += 1.0
		values[i] = f(simplex[i])
	}

	for iter := 0; iter < iterations; iter++ {
		// Sort vertices by function value
		for i := 0; i < n+1; i++ {
			for j := i + 1; j < n+1; j++ {
				if values[j] < values[i] {
					simplex[i], simplex[j] = simplex[j], simplex[i]
					values[i], values[j] = values[j], values[i]
				}
			}
		}

		// Compute centroid (excluding worst point)
		centroid := make([]float64, n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				centroid[j] += simplex[i][j]
			}
		}
		for j := 0; j < n; j++ {
			centroid[j] /= float64(n)
		}

		// Reflection
		reflected := make([]float64, n)
		for j := 0; j < n; j++ {
			reflected[j] = centroid[j] + alpha*(centroid[j]-simplex[n][j])
		}
		reflectedVal := f(reflected)

		//nolint:gocritic // Nelder-Mead algorithm requires floating point comparisons, not suitable for switch
		if reflectedVal < values[0] {
			// Expansion
			expanded := make([]float64, n)
			for j := 0; j < n; j++ {
				expanded[j] = centroid[j] + gamma*(reflected[j]-centroid[j])
			}
			expandedVal := f(expanded)

			if expandedVal < reflectedVal {
				simplex[n] = expanded
				values[n] = expandedVal
			} else {
				simplex[n] = reflected
				values[n] = reflectedVal
			}
		} else if reflectedVal < values[n-1] {
			simplex[n] = reflected
			values[n] = reflectedVal
		} else {
			// Contraction
			contracted := make([]float64, n)
			if reflectedVal < values[n] {
				// Outside contraction
				for j := 0; j < n; j++ {
					contracted[j] = centroid[j] + rho*(reflected[j]-centroid[j])
				}
			} else {
				// Inside contraction
				for j := 0; j < n; j++ {
					contracted[j] = centroid[j] + rho*(simplex[n][j]-centroid[j])
				}
			}
			contractedVal := f(contracted)

			if contractedVal < values[n] {
				simplex[n] = contracted
				values[n] = contractedVal
			} else {
				// Shrink
				for i := 1; i <= n; i++ {
					for j := 0; j < n; j++ {
						simplex[i][j] = simplex[0][j] + sigma*(simplex[i][j]-simplex[0][j])
					}
					values[i] = f(simplex[i])
				}
			}
		}
	}

	// Return best point
	bestIdx := 0
	for i := 1; i <= n; i++ {
		if values[i] < values[bestIdx] {
			bestIdx = i
		}
	}

	return simplex[bestIdx]
}

// ConjugateGradient performs conjugate gradient optimization
// Time: O(iterations * d), Space: O(d)
func ConjugateGradient(
	f OptimizationFunc,
	grad GradientFunc,
	initial []float64,
	iterations int,
	tolerance float64,
) []float64 {
	x := make([]float64, len(initial))
	copy(x, initial)

	g := grad(x)
	d := make([]float64, len(g))
	for i := range d {
		d[i] = -g[i]
	}

	for iter := 0; iter < iterations; iter++ {
		// Line search (simple backtracking)
		alpha := 1.0
		xNew := make([]float64, len(x))
		for i := 0; i < 10; i++ {
			for j := range xNew {
				xNew[j] = x[j] + alpha*d[j]
			}
			if f(xNew) < f(x) {
				break
			}
			alpha *= 0.5
		}

		// Update x
		for i := range x {
			x[i] += alpha * d[i]
		}

		// Compute new gradient
		gNew := grad(x)

		// Check convergence
		norm := 0.0
		for i := range gNew {
			norm += gNew[i] * gNew[i]
		}
		if math.Sqrt(norm) < tolerance {
			break
		}

		// Compute beta (Fletcher-Reeves)
		numerator := 0.0
		denominator := 0.0
		for i := range gNew {
			numerator += gNew[i] * gNew[i]
			denominator += g[i] * g[i]
		}
		beta := numerator / denominator

		// Update search direction
		for i := range d {
			d[i] = -gNew[i] + beta*d[i]
		}

		g = gNew
	}

	return x
}

// BFGS performs BFGS quasi-Newton optimization
// Time: O(iterations * d²), Space: O(d²)
func BFGS(
	f OptimizationFunc,
	grad GradientFunc,
	initial []float64,
	iterations int,
	tolerance float64,
) []float64 {
	n := len(initial)
	x := make([]float64, n)
	copy(x, initial)

	// Initialize inverse Hessian approximation as identity
	H := make([][]float64, n)
	for i := range H {
		H[i] = make([]float64, n)
		H[i][i] = 1.0
	}

	g := grad(x)

	for iter := 0; iter < iterations; iter++ {
		// Compute search direction: d = -H * g
		d := make([]float64, n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				d[i] -= H[i][j] * g[j]
			}
		}

		// Line search
		alpha := 1.0
		xNew := make([]float64, n)
		for i := 0; i < 10; i++ {
			for j := range xNew {
				xNew[j] = x[j] + alpha*d[j]
			}
			if f(xNew) < f(x) {
				break
			}
			alpha *= 0.5
		}

		// Update x
		s := make([]float64, n)
		for i := range x {
			s[i] = alpha * d[i]
			x[i] += s[i]
		}

		// Compute new gradient
		gNew := grad(x)

		// Compute gradient difference
		y := make([]float64, n)
		for i := range y {
			y[i] = gNew[i] - g[i]
		}

		// Check convergence
		norm := 0.0
		for i := range gNew {
			norm += gNew[i] * gNew[i]
		}
		if math.Sqrt(norm) < tolerance {
			break
		}

		// BFGS update: H_{k+1} = (I - rho*s*y^T) * H_k * (I - rho*y*s^T) + rho*s*s^T
		rho := 0.0
		for i := 0; i < n; i++ {
			rho += y[i] * s[i]
		}
		if rho > 0 {
			rho = 1.0 / rho

			// Compute I - rho*s*y^T
			A := make([][]float64, n)
			for i := range A {
				A[i] = make([]float64, n)
				A[i][i] = 1.0
			}
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					A[i][j] -= rho * s[i] * y[j]
				}
			}

			// Compute A * H
			AH := make([][]float64, n)
			for i := range AH {
				AH[i] = make([]float64, n)
			}
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					for k := 0; k < n; k++ {
						AH[i][j] += A[i][k] * H[k][j]
					}
				}
			}

			// Compute A * H * A^T (where A^T = I - rho*y*s^T)
			HNew := make([][]float64, n)
			for i := range HNew {
				HNew[i] = make([]float64, n)
			}
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					// AH[i] * (I - rho*y*s^T)[:][j]
					for k := 0; k < n; k++ {
						delta := 0.0
						if k == j {
							delta = 1.0
						}
						HNew[i][j] += AH[i][k] * (delta - rho*y[k]*s[j])
					}
					// Add rho*s*s^T
					HNew[i][j] += rho * s[i] * s[j]
				}
			}

			H = HNew
		}

		g = gNew
	}

	return x
}

// DifferentialEvolution performs differential evolution
// Time: O(generations * popSize * d), Space: O(popSize * d)
func DifferentialEvolution(
	f OptimizationFunc,
	dimensions int,
	bounds [][]float64,
	popSize int,
	generations int,
	mutationFactor float64,
	crossoverProb float64,
) []float64 {
	// Initialize population
	population := make([][]float64, popSize)
	fitness := make([]float64, popSize)

	for i := range population {
		population[i] = make([]float64, dimensions)
		for j := range population[i] {
			population[i][j] = bounds[j][0] + rand.Float64()*(bounds[j][1]-bounds[j][0])
		}
		fitness[i] = f(population[i])
	}

	for gen := 0; gen < generations; gen++ {
		for i := 0; i < popSize; i++ {
			// Select three random distinct individuals
			indices := rand.Perm(popSize)
			a, b, c := indices[0], indices[1], indices[2]
			for a == i {
				a = rand.IntN(popSize)
			}
			for b == i || b == a {
				b = rand.IntN(popSize)
			}
			for c == i || c == a || c == b {
				c = rand.IntN(popSize)
			}

			// Mutation and crossover
			trial := make([]float64, dimensions)
			jrand := rand.IntN(dimensions)

			for j := 0; j < dimensions; j++ {
				if rand.Float64() < crossoverProb || j == jrand {
					trial[j] = population[a][j] +
						mutationFactor*(population[b][j]-population[c][j])

					// Clamp to bounds
					if trial[j] < bounds[j][0] {
						trial[j] = bounds[j][0]
					}
					if trial[j] > bounds[j][1] {
						trial[j] = bounds[j][1]
					}
				} else {
					trial[j] = population[i][j]
				}
			}

			// Selection
			trialFitness := f(trial)
			if trialFitness < fitness[i] {
				population[i] = trial
				fitness[i] = trialFitness
			}
		}
	}

	// Find best
	bestIdx := 0
	for i := 1; i < popSize; i++ {
		if fitness[i] < fitness[bestIdx] {
			bestIdx = i
		}
	}

	return population[bestIdx]
}
