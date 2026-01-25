# AGENT GUIDELINES FOR THE `sopra` REPOSITORY

This document outlines the conventions, commands, and code style guidelines for agentic coding agents operating within the `sopra` Go project. Adhering to these guidelines ensures consistency, maintainability, and compatibility with the existing codebase and CI/CD pipelines.

## 1. Build, Lint, and Test Commands

These commands are essential for developing, verifying, and maintaining the codebase.

### 1.1 Building the Project

*   **Local Development Build**: For local compilation and execution of the Go application.
    ```bash
    go build -o sopra .
    ```
    To run: `./sopra`

*   **Containerized Build (CI/CD)**: The project uses Cloud Native Buildpacks for containerization in CI/CD. Do not use this for local development unless specifically instructed.
    ```bash
    pack build sopra-server --builder gcr.io/buildpacks/builder:v1
    ```

### 1.2 Testing

*   **Run All Tests**: Executes all Go tests within the project.
    ```bash
    go test ./...
    ```

*   **Run Tests in a Specific Package**:
    ```bash
    go test ./path/to/package
    ```
    Example: `go test ./client`

*   **Run a Single Test Function**: To execute a specific test function in a file.
    ```bash
    go test -run <TestFunctionName> ./path/to/package/<file_test.go>
    ```
    Example: `go test -run TestFlightAwareClient_GetAircraft ./client/flightaware_test.go`

*   **Run Tests with Verbose Output**:
    ```bash
    go test -v ./...
    ```

### 1.3 Linting and Formatting

*   **Formatting (gofmt)**: This is the standard Go formatter. Always run this before committing code.
    ```bash
    gofmt -w .
    ```

*   **Import Organization (goimports)**: Automatically adds and removes Go imports, and formats the file. `gofmt` is included in `goimports`.
    ```bash
    goimports -w .
    ```
    *Note: Ensure `goimports` is installed (`go install golang.org/x/tools/cmd/goimports@latest`).*

*   **Basic Static Analysis (go vet)**: Reports suspicious constructs, like `Printf` calls whose arguments don't align with the format string.
    ```bash
    go vet ./...
    ```

*   **Comprehensive Linting (golangci-lint)**: For more robust static analysis, `golangci-lint` is recommended. While not explicitly configured in CI, it's a good practice for agents to use for thorough code quality checks.
    ```bash
    golangci-lint run ./...
    ```
    *Note: If `golangci-lint` is not installed, install it (`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`).*

## 2. Code Style Guidelines

Adherence to Go's idiomatic style is paramount.

### 2.1 Imports

*   Use `goimports` to manage and format imports automatically.
*   Imports should be grouped: standard library, then third-party packages, then internal project packages, separated by blank lines.

### 2.2 Formatting

*   **Always use `gofmt`**. Consistent formatting is critical for readability.
*   Line length: Aim for readability, typically around 100-120 characters per line, but no hard limit. Break long lines logically.

### 2.3 Types

*   Use explicit types. Go is a statically typed language.
*   Prefer interfaces for defining behavior, especially for dependency injection, rather than concrete types where flexibility is needed.
*   Structs should group related data. Consider embedding for composition.

### 2.4 Naming Conventions

*   **Exported Identifiers**: Use `CamelCase` (starts with an uppercase letter) for package-level names (variables, functions, types, constants) that are intended to be exported (visible outside the package).
*   **Unexported Identifiers**: Use `camelCase` (starts with a lowercase letter) for package-level names that are unexported (private to the package).
*   **Acronyms**: Keep acronyms all uppercase (e.g., `HTTP`, `URL`, `ID`, `API`).
*   **Package Names**: Should be short, all lowercase, and reflect the package's purpose (e.g., `client`, `service`, `config`). Avoid plural names.
*   **Variable Names**: Be descriptive. Short variable names are acceptable for short-lived, local variables (e.g., `i` for loop index, `err` for error).
*   **Interface Names**: Often end with `er` (e.g., `Reader`, `Writer`), or describe what they `Do` (e.g., `Formatter`).

### 2.5 Error Handling

*   **Return Errors**: Functions should return errors as their last return value.
*   **Check Errors Immediately**: Always check for errors immediately after a function call that might return an error.
    ```go
    if err != nil {
        // Handle error
    }
    ```
*   **Error Wrapping**: Use `fmt.Errorf` with `%w` for wrapping errors to preserve the original error chain for debugging.
*   **Sentinel Errors/Custom Error Types**: Define custom error types or sentinel errors for specific error conditions where programmatic handling is required.

### 2.6 Comments

*   Comments should explain *why* something is done, rather than *what* is done (which should be clear from the code itself).
*   Document all exported functions, types, and variables with clear, concise comments.
*   Start comments for exported symbols with the symbol's name (e.g., `// FuncName does X.`).

### 2.7 Code Structure

*   Organize code into logical packages based on functionality.
*   Keep functions concise and focused on a single responsibility.
*   Avoid deeply nested control structures.

### 2.8 Concurrency

*   Use Go's concurrency primitives (goroutines, channels) responsibly.
*   Protect shared data with `sync` package primitives (e.g., `sync.Mutex`, `sync.RWMutex`).
*   Avoid using `time.Sleep` for synchronization; use channels or `sync` primitives instead.

### 2.9 Dependencies

*   New dependencies must be carefully considered and justified.
*   Add dependencies using `go get <package>` to update `go.mod` and `go.sum`.

## 3. Existing Agent Rules

No specific `.cursor/rules/` or `.github/copilot-instructions.md` files were found in this repository. Agents should adhere strictly to the Go language and project-specific guidelines outlined above.

---
*This `AGENTS.md` file was generated to assist agentic coding agents in contributing effectively to the `sopra` project.*
