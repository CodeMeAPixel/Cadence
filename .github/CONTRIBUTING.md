# Contributing to Cadence

Thank you for your interest in contributing to Cadence! We welcome contributions from everyone.

## Code of Conduct

This project adheres to the [Code of Conduct](.github/CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to [hey@codemeapixel.dev](mailto:hey@codemeapixel.dev).

## How to Contribute

### Reporting Bugs

Before reporting a bug, please:

1. Check the [GitHub Issues](https://github.com/CodeMeAPixel/cadence/issues) to see if it's already reported
2. Check the [Security Policy](.github/SECURITY.md) if it's security-related

When reporting a bug, include:
- Clear title and description
- Steps to reproduce
- Expected vs actual behavior
- Go version and OS
- Relevant error messages or logs

### Suggesting Features

1. Check if the feature is already suggested in [Issues](https://github.com/CodeMeAPixel/cadence/issues)
2. Provide a clear use case and benefits
3. Include examples if possible
4. Be open to feedback and discussion

### Submitting Changes

#### Setup

```bash
# Clone the repository
git clone https://github.com/CodeMeAPixel/cadence.git
cd cadence

# Create a feature branch
git checkout -b feature/your-feature-name
```

#### Development

```bash
# Install dependencies
go mod download

# Build the project
make build              # Linux/macOS
.\scripts\build.ps1     # Windows

# Run tests
make test

# Format code
make fmt

# Run linter
make lint

# Run vet
make vet
```

#### Commit Guidelines

- Use clear, descriptive commit messages
- Start with a verb (Add, Fix, Update, Remove, etc.)
- Reference issues when applicable: `Fixes #123`
- Keep commits focused and atomic

Example:
```
Fix: correct version injection in Windows builds

- Update Makefile OS detection logic
- Add proper error handling in build scripts
- Closes #456
```

#### Code Style

- Follow Go conventions and idioms
- Use `gofmt` for formatting (enforced by `make fmt`)
- Write meaningful variable names
- Add comments for exported functions
- Keep functions focused and testable

#### Testing

- Add tests for new features
- Update tests when behavior changes
- Aim for reasonable coverage (not 100% required)
- Run tests before submitting: `make test`

```bash
# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestFunctionName ./...
```

#### Documentation

- Update README.md if behavior changes
- Add comments to exported functions
- Update CHANGELOG.md with your changes
- Include examples in commit messages when helpful

### Pull Request Process

1. **Fork** the repository
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Make your changes** and test thoroughly
4. **Commit with clear messages** (`git commit -m "Add amazing feature"`)
5. **Push to your fork** (`git push origin feature/amazing-feature`)
6. **Open a Pull Request** with:
   - Clear title and description
   - Link to related issues (`Fixes #123`)
   - Summary of changes
   - Any breaking changes clearly noted

#### PR Template

When you open a PR, please include:

```markdown
## Description
Brief description of changes

## Fixes
Fixes #(issue number)

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests added
- [ ] All tests passing
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] No new warnings generated
```

### Review Process

1. At least one maintainer will review your PR
2. Changes may be requested
3. Once approved, your PR will be merged
4. You'll be credited in CHANGELOG.md and commit history

## Development Resources

- [README.md](../README.md) - Project overview and usage
- [Project Structure](../README.md#project-structure) - Code organization
- [Security Policy](.github/SECURITY.md) - Vulnerability reporting

## Questions?

- **GitHub Issues** - For bug reports and feature requests
- **Email** - [hey@codemeapixel.dev](mailto:hey@codemeapixel.dev)
- **Social Media** - [@CodeMeAPixel](https://twitter.com/CodeMeAPixel) on Twitter/X

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (see LICENSE file).

## Recognition

Contributors are recognized in:
- Pull request comments
- CHANGELOG.md entries
- GitHub's contributor graph
- Release notes for significant contributions

Thank you for making Cadence better! 🎉
