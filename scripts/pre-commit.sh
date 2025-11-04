#!/bin/bash
# Pre-commit hook for go-distance
#
# To install this hook, run from the project root:
#   cp scripts/pre-commit.sh .git/hooks/pre-commit
#   chmod +x .git/hooks/pre-commit
#
# Or use the install command:
#   make install-hooks

set -e

echo "Running pre-commit checks..."

# Run the pre-commit target from Makefile
make pre-commit

echo ""
echo "✓ All pre-commit checks passed!"
echo "Proceeding with commit..."
