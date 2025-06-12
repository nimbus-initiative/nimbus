# Contributing to Nimbus

Thank you for your interest in contributing to Nimbus! We welcome all contributions, whether they're bug reports, feature requests, documentation improvements, or code contributions.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Environment](#development-environment)
- [Making Changes](#making-changes)
- [Submitting a Pull Request](#submitting-a-pull-request)
- [Reporting Issues](#reporting-issues)
- [Code Style](#code-style)
- [Testing](#testing)
- [Documentation](#documentation)
- [Community](#community)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. **Fork** the repository on GitHub
2. **Clone** the project to your own machine
3. **Commit** changes to your own branch
4. **Push** your work back up to your fork
5. Submit a **Pull Request** so that we can review your changes

## Development Environment

### Prerequisites

- Go 1.18 or later
- Git
- Make
- Docker (optional, for containerized development)

### Setting Up

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/nimbus.git
   cd nimbus
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   make
   ```

## Making Changes

1. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b bugfix/issue-number-description
   ```

2. Make your changes following the [code style](#code-style)

3. Run tests:
   ```bash
   make test
   ```

4. Run linters:
   ```bash
   make lint
   ```

5. Commit your changes with a descriptive commit message:
   ```bash
   git commit -m "Add feature: brief description of changes"
   ```

## Submitting a Pull Request

1. Push your changes to your fork:
   ```bash
   git push origin your-branch-name
   ```

2. Open a Pull Request (PR) from your fork to the main repository

3. Fill out the PR template with all relevant information

4. Ensure all CI checks pass

5. Address any review comments

## Reporting Issues

When reporting issues, please include:

- A clear title and description
- Steps to reproduce the issue
- Expected vs. actual behavior
- Any relevant logs or screenshots
- Your environment details (OS, version, etc.)

## Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `gofmt` and `goimports` on your code before committing
- Keep lines under 120 characters
- Write clear, concise commit messages
- Document exported functions and types

## Testing

- Write unit tests for new functionality
- Ensure all tests pass before submitting a PR
- Add integration tests for complex features
- Update tests when fixing bugs

## Documentation

- Update relevant documentation when adding or changing features
- Keep comments clear and concise
- Document any breaking changes
- Update README.md with new features or changes

## Community

- Join our [Discord/Slack/Mailing List] for discussions
- Be respectful and inclusive in all communications
- Help others in the community when possible

## License

By contributing to Nimbus, you agree that your contributions will be licensed under the project's [LICENSE](LICENSE) file.
