#!/bin/bash

# Git commit message hook to enforce conventional commits
# This validates commit messages follow the format: type(scope): subject

COMMIT_MSG_FILE=$1
COMMIT_MSG=$(cat "$COMMIT_MSG_FILE")

# Skip for merge commits
if echo "$COMMIT_MSG" | grep -qE "^Merge"; then
    exit 0
fi

# Conventional commit pattern
PATTERN="^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\(.+\))?: .{1,}$"

if ! echo "$COMMIT_MSG" | grep -qE "$PATTERN"; then
    echo ""
    echo "❌ Invalid commit message format!"
    echo ""
    echo "Commit messages must follow Conventional Commits:"
    echo "  <type>[optional scope]: <description>"
    echo ""
    echo "Types:"
    echo "  feat:     A new feature"
    echo "  fix:      A bug fix"
    echo "  docs:     Documentation only changes"
    echo "  style:    Code style changes (formatting, etc)"
    echo "  refactor: Code change that neither fixes a bug nor adds a feature"
    echo "  perf:     Performance improvement"
    echo "  test:     Adding or updating tests"
    echo "  build:    Changes to build system or dependencies"
    echo "  ci:       Changes to CI configuration"
    echo "  chore:    Other changes that don't modify src or test files"
    echo "  revert:   Revert a previous commit"
    echo ""
    echo "Examples:"
    echo "  feat: add Euclidean distance metric"
    echo "  feat(optimization): implement Adam optimizer"
    echo "  fix: correct Vincenty formula edge case"
    echo "  docs: update README with new examples"
    echo "  chore: update dependencies"
    echo ""
    echo "Breaking changes (major version bump):"
    echo "  feat!: change API signature"
    echo "  feat(api)!: rename Distance to Compute"
    echo ""
    exit 1
fi

exit 0
