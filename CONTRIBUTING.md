# Contributing to watch-fs

Thank you for your interest in contributing to watch-fs! This document provides guidelines for contributing to the project.

## Development Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/pbouamriou/watch-fs.git
   cd watch-fs
   ```

2. **Install dependencies**

   ```bash
   make deps
   ```

3. **Build the project**
   ```bash
   make build
   ```

## Project Structure

- `cmd/watch-fs/`: Main application entry point
- `internal/ui/`: Terminal user interface implementation
- `internal/watcher/`: File system watcher wrapper
- `pkg/utils/`: Utility functions
- `docs/`: Documentation
- `test/`: Test files

## Development Workflow

1. **Create a feature branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**

   - Follow Go coding standards
   - Add tests for new functionality
   - Update documentation as needed

3. **Run tests**

   ```bash
   make test
   ```

4. **Build and test**

   ```bash
   make build
   ./bin/watch-fs -path . -tui=false
   ```

5. **Commit your changes**

   ```bash
   git add .
   git commit -m "feat: Add your feature description"
   ```

6. **Push and create a pull request**

## Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions small and focused

## Testing

- Write unit tests for new functionality
- Test both TUI and console modes
- Ensure error handling is properly tested

## Documentation

- Update README.md for user-facing changes
- Update docs/ARCHITECTURE.md for architectural changes
- Add inline comments for complex code

## Commit Messages

Use conventional commit format:

- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `refactor:` for code refactoring
- `test:` for test changes

Example:

```
feat: Add filtering by file extension

- Add new filter option for file extensions
- Update UI to show extension filter status
- Add tests for extension filtering
```

## Questions?

If you have questions about contributing, please open an issue on GitHub.
