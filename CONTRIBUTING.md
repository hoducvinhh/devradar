# Contributing to DevRadar

Thank you for your interest in contributing to **DevRadar**! We are building tools to empower developers everywhere to take their first or next step into the open-source ecosystem.

## How Can I Contribute?

### 1. Reporting Bugs
- Search existing issues to ensure the bug hasn't already been reported.
- Use the **Bug Report** template and provide detailed reproduction steps.

### 2. Suggesting Features
- Open an issue using the **Feature Request** template.
- Explain the motivation, use case, and expected UX.

### 3. Submitting Code
1. Fork the repository on GitHub.
2. Clone your fork locally:
   ```bash
   git clone https://github.com/<your-username>/devradar.git
   cd devradar
   ```
3. Create a descriptive feature branch:
   ```bash
   git checkout -b feat/your-feature-name
   ```
4. Write your changes and add unit tests.
5. Verify tests and style:
   ```bash
   go test -v ./...
   go vet ./...
   ```
6. Commit your changes with clear semantic commit messages:
   ```bash
   git commit -m "feat(ranking): add repo freshness penalty"
   ```
7. Push to your fork and submit a Pull Request.

## Code Standards
- **Pure Go**: Do not introduce dependencies requiring CGO or GCC.
- **Test Coverage**: Maintain test coverage > 80% for new logic.
- **Formatting**: Run `gofmt` before committing.
- **Privacy First**: DevRadar never sends user telemetry or logs tokens.
