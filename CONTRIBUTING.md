# Contributing to csv2json

Thank you for your interest in contributing to csv2json! We welcome contributions from the community.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Running Tests](#running-tests)
- [Making Changes](#making-changes)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Reporting Bugs](#reporting-bugs)
- [Feature Requests](#feature-requests)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Set up the development environment
4. Create a new branch for your feature or bugfix
5. Make your changes
6. Run tests to ensure everything works
7. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.25.4 or higher
- PostgreSQL 12 or higher
- Git

### Installation Steps

1. **Clone your fork:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/csv2json.git
   cd csv2json
   ```

2. **Add upstream remote:**
   ```bash
   git remote add upstream https://github.com/agileproject-gurpreet/csv2json.git
   ```

3. **Install dependencies:**
   ```bash
   go mod download
   ```

4. **Set up PostgreSQL:**
   ```bash
   # Create database
   createdb csv2json_dev
   
   # Run setup script
   psql -d csv2json_dev -f docs/setup.sql
   ```

5. **Configure environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your local database credentials
   ```

6. **Run the application:**
   ```bash
   go run cmd/api/main.go
   ```

## Running Tests

We maintain comprehensive test coverage for all critical components.

### Run all tests:
```bash
go test ./...
```

### Run tests with coverage:
```bash
go test ./... -cover
```

### Run tests with verbose output:
```bash
go test ./... -v
```

### Run specific package tests:
```bash
go test ./internal/parser/tests/
```

### Generate coverage report:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Before submitting a PR:
```bash
# Run all tests
go test ./...

# Check formatting
go fmt ./...

# Run go vet
go vet ./...

# Run golint (if installed)
golint ./...
```

## Making Changes

### Branch Naming Convention

- `feature/description` - for new features
- `fix/description` - for bug fixes
- `docs/description` - for documentation updates
- `refactor/description` - for code refactoring
- `test/description` - for test improvements

Example: `feature/add-xml-support` or `fix/handle-empty-csv`

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(parser): add support for custom delimiters

fix(api): handle empty CSV files gracefully

docs(readme): update installation instructions

test(parser): add edge case tests for malformed CSV
```

## Pull Request Process

### Before Submitting

1. **Update your branch:**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Ensure all tests pass:**
   ```bash
   go test ./...
   ```

3. **Format your code:**
   ```bash
   go fmt ./...
   ```

4. **Update documentation** if needed (README.md, API docs, etc.)

### Submitting a Pull Request

1. **Push your branch to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a Pull Request** on GitHub with:
   - Clear title describing the change
   - Detailed description of what changed and why
   - Reference any related issues (e.g., "Fixes #123")
   - Screenshots or examples if applicable

3. **PR Description Template:**
   ```markdown
   ## Description
   Brief description of the changes

   ## Type of Change
   - [ ] Bug fix
   - [ ] New feature
   - [ ] Breaking change
   - [ ] Documentation update

   ## Testing
   Describe the tests you ran

   ## Checklist
   - [ ] My code follows the project's style guidelines
   - [ ] I have performed a self-review of my code
   - [ ] I have commented my code where necessary
   - [ ] I have updated the documentation
   - [ ] My changes generate no new warnings
   - [ ] I have added tests that prove my fix/feature works
   - [ ] New and existing unit tests pass locally
   ```

### Review Process

1. At least one maintainer will review your PR
2. Address any feedback or requested changes
3. Once approved, a maintainer will merge your PR
4. Your contribution will be included in the next release!

### After Your PR is Merged

1. Delete your branch:
   ```bash
   git branch -d feature/your-feature-name
   git push origin --delete feature/your-feature-name
   ```

2. Update your main branch:
   ```bash
   git checkout main
   git pull upstream main
   ```

## Coding Standards

### Go Guidelines

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `go fmt` for formatting
- Keep functions small and focused
- Write clear, self-documenting code
- Use meaningful variable and function names
- Add comments for exported functions and types

### Code Organization

```
csv2json/
├── cmd/           # Application entry points
├── internal/      # Private application code
│   ├── database/  # Database operations
│   ├── handler/   # HTTP handlers
│   ├── parser/    # CSV parsing logic
│   └── service/   # Business logic
├── pkg/           # Public libraries
└── docs/          # Documentation
```

### Error Handling

- Always handle errors explicitly
- Provide meaningful error messages
- Use wrapped errors for context: `fmt.Errorf("operation failed: %w", err)`

### Testing

- Write tests for all new functionality
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Test edge cases and error conditions

## Reporting Bugs

### Before Submitting a Bug Report

1. Check existing issues to avoid duplicates
2. Ensure you're using the latest version
3. Verify the bug is reproducible

### Bug Report Template

```markdown
**Describe the bug**
A clear description of the bug

**To Reproduce**
Steps to reproduce:
1. 
2. 
3. 

**Expected behavior**
What you expected to happen

**Actual behavior**
What actually happened

**Environment:**
- OS: [e.g., Windows 11, Ubuntu 22.04]
- Go version: [e.g., 1.25.4]
- PostgreSQL version: [e.g., 14.2]

**Additional context**
Any other relevant information
```

## Feature Requests

We welcome feature requests! Please:

1. Check if the feature has already been requested
2. Provide a clear use case
3. Explain how it benefits the project
4. Be open to discussion and feedback

## Questions?

- Open an issue with the `question` label
- Reach out to the maintainers

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

Thank you for contributing to csv2json! 🎉
