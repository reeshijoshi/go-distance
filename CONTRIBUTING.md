# Contributing to go-distance

Thank you for your interest in contributing to go-distance! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Code Style](#code-style)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/go-distance.git
   cd go-distance
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/reeshijoshi/go-distance.git
   ```

## Development Setup

### Prerequisites

- Go 1.21 or later
- Make (for running Makefile commands)
- Git

### Install Development Tools

```bash
make install-tools
```

This installs:
- `golangci-lint` - Comprehensive linter
- `gosec` - Security checker (optional)

### Install Git Hooks

Install the pre-commit hook to automatically run checks before each commit:

```bash
make install-hooks
```

## Making Changes

### Before You Start

1. **Check existing issues** - Someone might already be working on it
2. **Create an issue** for major changes to discuss the approach
3. **Keep changes focused** - One feature/fix per PR

### Development Workflow

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our code style guidelines

3. **Run pre-commit checks**:
   ```bash
   make pre-commit
   ```

4. **Commit your changes** using [Conventional Commits](https://www.conventionalcommits.org/):
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

   **Commit message format** (enforced by commit-msg hook):
   - Format: `<type>[optional scope]: <description>`
   - Types: `feat` (new feature), `fix` (bug fix), `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`
   - Breaking changes: Add `!` after type (e.g., `feat!:`)
   - Examples:
     - `feat: add Mahalanobis distance` → triggers v0.1.0 to v0.2.0
     - `fix: correct Vincenty edge case` → triggers v0.1.0 to v0.1.1
     - `feat!: change API signature` → triggers v0.1.0 to v1.0.0

   The commit-msg hook validates this format and will show examples if invalid.

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage
make test-coverage

# View coverage in browser
make test-coverage-html

# Run race detector
make test-race

# Run benchmarks
make bench
```

### Writing Tests

- All new code must include tests
- Aim for >80% code coverage
- Test edge cases and error conditions
- Use table-driven tests where appropriate

Example test structure:

```go
func TestNewFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    []float64
        expected float64
        wantErr  bool
    }{
        {"basic case", []float64{1, 2, 3}, 2.0, false},
        {"error case", []float64{}, 0.0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := NewFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
            }
            if !tt.wantErr && !almostEqual(result, tt.expected) {
                t.Errorf("expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

## Code Style

### Go Best Practices

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting (run `make fmt`)
- Pass `go vet` checks (run `make vet`)
- Pass `golangci-lint` (run `make lint`)

### Code Guidelines

1. **Function Documentation**:
   ```go
   // Euclidean computes the L2 norm (straight-line distance) between two vectors.
   // Time: O(n), Space: O(1)
   func Euclidean[T Number](a, b []T) (float64, error) {
   ```

2. **Error Handling**:
   - Return errors explicitly
   - Use descriptive error messages
   - Define package-level errors for common cases

3. **Generics**:
   - Use type constraints appropriately
   - Document type parameters

4. **Performance**:
   - Minimize allocations in hot paths
   - Document time/space complexity
   - Add benchmarks for performance-critical code

### Quality Checks

Run all quality checks before submitting:

```bash
make check
```

This runs:
- `make fmt-check` - Verify formatting
- `make vet` - Run go vet
- `make tidy-check` - Check go.mod
- `make lint` - Run golangci-lint
- `make test` - Run all tests

## Submitting Changes

### Pull Request Process

1. **Update your branch** with latest upstream:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run pre-commit checks**:
   ```bash
   make pre-commit
   ```

3. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

4. **Create Pull Request** on GitHub with:
   - Clear description of changes
   - Reference related issues
   - Screenshots/examples if applicable

### PR Requirements

- [ ] All tests pass (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] Linting passes (`make lint`)
- [ ] Documentation is updated
- [ ] Tests added for new functionality
- [ ] No breaking changes (or discussed in issue)

### Review Process

1. Maintainers will review your PR
2. Address any feedback or requested changes
3. Once approved, maintainers will merge

## Adding New Distance Metrics

When adding a new distance metric:

1. **Implementation**:
   - Add function to appropriate file (e.g., `vector.go`, `string.go`)
   - Include comprehensive documentation
   - Specify time/space complexity
   - Handle errors appropriately

2. **Testing**:
   - Add comprehensive tests in `*_test.go`
   - Include edge cases
   - Add benchmarks

3. **Documentation**:
   - Update README.md with new function
   - Add usage example
   - Document use cases

4. **Example**:
   ```go
   // MyDistance computes distance between two vectors.
   // Time: O(n), Space: O(1)
   func MyDistance[T Number](a, b []T) (float64, error) {
       if err := Validate(a, b); err != nil {
           return 0, err
       }
       // Implementation
   }
   ```

## Release Process

Releases are **fully automated** using semantic-release:

1. Commit with conventional format (e.g., `feat:`, `fix:`, `feat!:`)
2. Push/merge to `main` branch
3. GitHub Actions automatically:
   - Analyzes commits and determines version bump
   - Updates CHANGELOG.md
   - Creates git tag
   - Publishes GitHub release
   - Triggers pkg.go.dev update

**Version bumps**:
- `feat:` commits → minor version (v0.1.0 → v0.2.0)
- `fix:`, `perf:`, `docs:` commits → patch version (v0.1.0 → v0.1.1)
- `feat!:` or `BREAKING CHANGE:` → major version (v0.1.0 → v1.0.0)

## Need Help?

- 📖 Read the [README](README.md)
- 🐛 [Open an issue](https://github.com/reeshijoshi/go-distance/issues)
- 💬 [Start a discussion](https://github.com/reeshijoshi/go-distance/discussions)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

Thank you for contributing to go-distance! 🚀
