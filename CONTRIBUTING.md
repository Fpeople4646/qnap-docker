# Contributing to qnap-docker

Thank you for your interest in contributing to qnap-docker! This document provides guidelines for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)

## Code of Conduct

By participating in this project, you are expected to uphold our Code of Conduct. Please report unacceptable behavior through GitHub issues.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone git@github.com:your-username/qnap-docker.git
   cd qnap-docker
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream git@github.com:scttfrdmn/qnap-docker.git
   ```

## Development Setup

### Prerequisites

- Go 1.21 or later
- Make
- Git
- Access to a QNAP NAS for testing (recommended)

### Installation

1. Install development dependencies:
   ```bash
   make deps
   ```

2. Build the project:
   ```bash
   make build
   ```

3. Run quality checks:
   ```bash
   make quality-check
   ```

## Making Changes

### Branch Naming

- Feature branches: `feature/description`
- Bug fixes: `fix/description`
- Documentation: `docs/description`

### Commit Messages

Follow conventional commit format:
```
type(scope): description

Longer description if needed.

- List of changes
- Another change

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

### Code Style

- Follow Go best practices
- Run `make fmt` before committing
- Ensure all linters pass with `make quality-check`
- Add comments for public functions and complex logic
- Use meaningful variable and function names

### QNAP-Specific Guidelines

- Test changes on actual QNAP hardware when possible
- Consider different QNAP models and QTS versions
- Handle CACHEDEV volume variations gracefully
- Provide clear error messages for QNAP-specific issues

## Testing

### Unit Tests

```bash
make test
```

### Integration Tests

If you have access to QNAP hardware:

```bash
NAS_HOST=your-qnap.local NAS_USER=admin make integration-test
```

### Quality Checks

```bash
make quality-check
```

This runs:
- `gofmt` - Code formatting
- `go vet` - Static analysis
- `staticcheck` - Advanced static analysis
- `ineffassign` - Ineffectual assignment detection
- `misspell` - Spelling mistakes
- `gocyclo` - Cyclomatic complexity
- `golint` - Go Report Card compliance

## Submitting Changes

### Pull Request Process

1. Ensure your branch is up to date:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. Run all checks:
   ```bash
   make clean && make all
   ```

3. Push your changes:
   ```bash
   git push origin your-branch-name
   ```

4. Create a pull request on GitHub

### Pull Request Requirements

- [ ] All tests pass
- [ ] Quality checks pass
- [ ] Documentation updated (if applicable)
- [ ] CHANGELOG.md updated (for significant changes)
- [ ] Tested on QNAP hardware (if possible)

## QNAP Testing Matrix

When testing, consider these environments:

### QNAP Models
- TS-x64 series (ARM64)
- TS-x53 series (AMD64)
- TS-x32 series (ARM)

### QTS Versions
- QTS 5.1.x (latest)
- QTS 5.0.x
- QTS 4.5.x (minimum supported)

### Container Station Versions
- Container Station 3.x
- Container Station 2.x (legacy)

## Release Process

Releases are handled by maintainers:

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create release tag
4. GitHub Actions builds and releases binaries
5. Update Homebrew tap (if applicable)

## Getting Help

- 📖 [Documentation](docs/)
- 🐛 [Issue Tracker](https://github.com/scttfrdmn/qnap-docker/issues)
- 💬 [Discussions](https://github.com/scttfrdmn/qnap-docker/discussions)

## Recognition

Contributors will be recognized in:
- README.md contributors section
- Release notes
- Git commit history

Thank you for contributing to qnap-docker! 🎉