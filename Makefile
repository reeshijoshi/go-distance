.PHONY: help test test-verbose test-coverage test-race bench lint fmt vet tidy check build clean install-tools pre-commit ci install-hooks uninstall-hooks

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=go-distance
GO=go
GOTEST=$(GO) test
GOVET=$(GO) vet
GOFMT=gofmt
GOLANGCI_LINT=golangci-lint
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Color output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

##@ General

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make $(BLUE)<target>$(NC)\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(BLUE)%-20s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

fmt: ## Format code with gofmt
	@echo "$(BLUE)Running gofmt...$(NC)"
	@$(GOFMT) -s -w .
	@echo "$(GREEN)✓ Formatting complete$(NC)"

fmt-check: ## Check if code is formatted
	@echo "$(BLUE)Checking formatting...$(NC)"
	@unformatted=$$($(GOFMT) -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "$(RED)✗ The following files need formatting:$(NC)"; \
		echo "$$unformatted"; \
		exit 1; \
	fi
	@echo "$(GREEN)✓ All files are properly formatted$(NC)"

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	@$(GOVET) ./...
	@echo "$(GREEN)✓ Vet complete$(NC)"

tidy: ## Run go mod tidy
	@echo "$(BLUE)Running go mod tidy...$(NC)"
	@$(GO) mod tidy
	@echo "$(GREEN)✓ go mod tidy complete$(NC)"

tidy-check: ## Check if go.mod are tidy
	@echo "$(BLUE)Checking go mod tidy...$(NC)"
	@$(GO) mod tidy
	@if [ -n "$$(git status --porcelain go.mod)" ]; then \
		echo "$(RED)✗ go.mod needs tidying$(NC)"; \
		git diff go.mod; \
		exit 1; \
	fi
	@echo "$(GREEN)✓ go.mod are tidy$(NC)"

lint: ## Run golangci-lint
	@echo "$(BLUE)Running golangci-lint...$(NC)"
	@if ! command -v $(GOLANGCI_LINT) > /dev/null; then \
		echo "$(RED)✗ golangci-lint not installed. Run 'make install-tools'$(NC)"; \
		exit 1; \
	fi
	@$(GOLANGCI_LINT) run --timeout=5m ./...
	@echo "$(GREEN)✓ Linting complete$(NC)"

lint-fix: ## Run golangci-lint with autofix
	@echo "$(BLUE)Running golangci-lint with --fix...$(NC)"
	@if ! command -v $(GOLANGCI_LINT) > /dev/null; then \
		echo "$(RED)✗ golangci-lint not installed. Run 'make install-tools'$(NC)"; \
		exit 1; \
	fi
	@$(GOLANGCI_LINT) run --fix --timeout=5m ./...
	@echo "$(GREEN)✓ Linting with fixes complete$(NC)"

##@ Testing

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	@$(GOTEST) -timeout=2m ./...
	@echo "$(GREEN)✓ Tests passed$(NC)"

test-verbose: ## Run tests with verbose output
	@echo "$(BLUE)Running tests (verbose)...$(NC)"
	@$(GOTEST) -v -timeout=2m ./...

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@$(GOTEST) -timeout=2m -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@$(GO) tool cover -func=$(COVERAGE_FILE)
	@echo "$(GREEN)✓ Coverage report generated: $(COVERAGE_FILE)$(NC)"

test-coverage-html: test-coverage ## Generate HTML coverage report
	@echo "$(BLUE)Generating HTML coverage report...$(NC)"
	@$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "$(GREEN)✓ HTML coverage report: $(COVERAGE_HTML)$(NC)"

test-race: ## Run tests with race detector
	@echo "$(BLUE)Running tests with race detector...$(NC)"
	@$(GOTEST) -race -timeout=5m ./...
	@echo "$(GREEN)✓ Race tests passed$(NC)"

bench: ## Run benchmarks
	@echo "$(BLUE)Running benchmarks...$(NC)"
	@$(GOTEST) -bench=. -benchmem -timeout=10m ./...

bench-compare: ## Run benchmarks and save for comparison
	@echo "$(BLUE)Running benchmarks...$(NC)"
	@$(GOTEST) -bench=. -benchmem -timeout=10m ./... | tee bench.txt
	@echo "$(GREEN)✓ Benchmark results saved to bench.txt$(NC)"

##@ Build

build: ## Build the project (validates compilation)
	@echo "$(BLUE)Building...$(NC)"
	@$(GO) build -v ./...
	@echo "$(GREEN)✓ Build complete$(NC)"

clean: ## Clean build artifacts and test cache
	@echo "$(BLUE)Cleaning...$(NC)"
	@$(GO) clean
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML) bench.txt
	@echo "$(GREEN)✓ Clean complete$(NC)"

##@ Tools

install-tools: ## Install development tools
	@echo "$(BLUE)Installing development tools...$(NC)"
	@echo "Installing golangci-lint..."
	@if ! command -v $(GOLANGCI_LINT) > /dev/null; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2; \
	else \
		echo "golangci-lint already installed"; \
	fi
	@echo "$(GREEN)✓ Tools installed$(NC)"

update-tools: ## Update development tools
	@echo "$(BLUE)Updating golangci-lint...$(NC)"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2
	@echo "$(GREEN)✓ Tools updated$(NC)"

##@ Quality Checks

check: fmt-check vet tidy-check lint test ## Run all checks (fmt, vet, tidy, lint, test)
	@echo "$(GREEN)✓ All checks passed!$(NC)"

pre-commit: ## Run all pre-commit checks (recommended before git commit)
	@echo "$(YELLOW)================================$(NC)"
	@echo "$(YELLOW)Running pre-commit checks...$(NC)"
	@echo "$(YELLOW)================================$(NC)"
	@$(MAKE) fmt
	@$(MAKE) tidy
	@$(MAKE) vet
	@$(MAKE) lint
	@$(MAKE) test
	@$(MAKE) test-race
	@echo ""
	@echo "$(GREEN)================================$(NC)"
	@echo "$(GREEN)✓ All pre-commit checks passed!$(NC)"
	@echo "$(GREEN)================================$(NC)"
	@echo ""
	@echo "$(BLUE)Ready to commit!$(NC)"

ci: ## Run CI checks (used in continuous integration)
	@echo "$(BLUE)Running CI checks...$(NC)"
	@$(MAKE) fmt-check
	@$(MAKE) tidy-check
	@$(MAKE) vet
	@$(MAKE) lint
	@$(MAKE) test-coverage
	@$(MAKE) test-race
	@echo "$(GREEN)✓ CI checks passed$(NC)"

##@ Verification

verify: ## Verify dependencies
	@echo "$(BLUE)Verifying dependencies...$(NC)"
	@$(GO) mod verify
	@echo "$(GREEN)✓ Dependencies verified$(NC)"

security: ## Run security checks with gosec
	@echo "$(BLUE)Running security checks...$(NC)"
	@if ! command -v gosec > /dev/null; then \
		echo "$(YELLOW)Installing gosec...$(NC)"; \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
	fi
	@gosec -quiet ./...
	@echo "$(GREEN)✓ Security checks passed$(NC)"

##@ Documentation

doc: ## Generate and serve documentation
	@echo "$(BLUE)Starting documentation server at http://localhost:6060$(NC)"
	@echo "$(YELLOW)Visit: http://localhost:6060/pkg/github.com/reeshijoshi/go-distance/$(NC)"
	@godoc -http=:6060

##@ Release

version: ## Show current version from git tag
	@git describe --tags --always --dirty

tag: ## Create a new git tag (usage: make tag VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)✗ VERSION is required. Usage: make tag VERSION=v1.0.0$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)Creating tag $(VERSION)...$(NC)"
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "$(GREEN)✓ Tag $(VERSION) created$(NC)"
	@echo "$(YELLOW)Push with: git push origin $(VERSION)$(NC)"

##@ Maintenance

install-hooks: ## Install git pre-commit hooks
	@echo "$(BLUE)Installing git hooks...$(NC)"
	@if [ ! -d .git ]; then \
		echo "$(RED)✗ Not a git repository$(NC)"; \
		exit 1; \
	fi
	@cp scripts/pre-commit.sh .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "$(GREEN)✓ Pre-commit hook installed$(NC)"
	@echo "$(YELLOW)The hook will run automatically before each commit$(NC)"

uninstall-hooks: ## Uninstall git pre-commit hooks
	@echo "$(BLUE)Uninstalling git hooks...$(NC)"
	@rm -f .git/hooks/pre-commit
	@echo "$(GREEN)✓ Pre-commit hook uninstalled$(NC)"

deps: ## Download dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@$(GO) mod download
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

upgrade: ## Upgrade dependencies to latest versions
	@echo "$(BLUE)Upgrading dependencies...$(NC)"
	@$(GO) get -u ./...
	@$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies upgraded$(NC)"

list-deps: ## List all dependencies
	@echo "$(BLUE)Project dependencies:$(NC)"
	@$(GO) list -m all

why: ## Show why a package is needed (usage: make why PKG=golang.org/x/tools)
	@if [ -z "$(PKG)" ]; then \
		echo "$(RED)✗ PKG is required. Usage: make why PKG=golang.org/x/tools$(NC)"; \
		exit 1; \
	fi
	@$(GO) mod why $(PKG)
