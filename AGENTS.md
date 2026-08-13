# Agent Guidelines for scootgo

This document provides guidelines for agents working within the `scootgo` repository. Adhering to these guidelines ensures consistency, maintainability, and quality across the codebase.

## 1. Build, Lint, and Test Commands

### General Commands:

* **Build:** `make build` or `go build -tags mysql -o build/scootgo .`
  * Builds the main application binary, ensuring all dependencies are resolved and the code is compiled.
* **Lint:** `make lint` or `golangci-lint run --new-from-rev=origin/main --build-tags=mysql`
  * Runs the `golangci-lint` tool to check for code quality and style issues, including potential bugs, dead code, and unhandled errors. This is crucial for maintaining a high code standard.
* **Format:** `make fmt` or `go fmt ./pkg/...`
  * Automatically formats Go source code according to Go's standard style. Always run this before committing to ensure consistent formatting across the project.
* **Test All:** `make test` or `go test -tags=mysql ./pkg/... -v -skip "^TestFunctional"`
  * Runs all unit tests within the `pkg/` directory, skipping functional tests which often require external services. This is for quick verification of isolated component logic.
* **Clean:** `make clean`
  * Removes build artifacts and temporary files, helping to maintain a clean working directory.

### Running a Single Test:

To run a specific test, use the `go test` command with the `-run` flag. The `-run` flag takes a regular expression that matches the name of the test function you want to execute.

Example:

```bash
go test -tags=mysql ./pkg/path/to/your/package -v -run "TestSpecificFunctionName"
```

Replace `./pkg/path/to/your/package` with the actual path to the package containing the test, and `"TestSpecificFunctionName"` with the exact name or a regex pattern for your test function. This allows for focused testing during development.

## 2. Code Style Guidelines

### General Principles:

* **Consistency:** Always prioritize consistency with existing code. Observe the patterns and styles already present in the files you are modifying. This ensures a uniform codebase that is easy to navigate and understand.
* **Readability:** Write code that is easy to understand and follow. Use clear variable names, concise logic, and appropriate comments where necessary.
* **Simplicity:** Favor simple, straightforward solutions over complex ones. Avoid over-engineering; a simpler solution is often more maintainable and less prone to bugs.
* **Performance:** While simplicity is key, be mindful of performance in critical paths. Profile and optimize where necessary, but don't sacrifice readability for minor performance gains unless absolutely essential.

### Language-Specific (Go) Guidelines:

* **Imports:**
  * Group imports into standard library, then third-party, then internal project packages.
  * Each group should be separated by a blank line.
  * Place each import in the correct group and in sorted order before committing. Mis-grouped or unsorted imports will be reordered by `go fmt`, producing noisy diffs that obscure the real change.
  * Example:

    ```go
    package main

    import (
        "fmt"
        "net/http"

        "github.com/some/thirdparty"

        "scootgo/pkg/internal/mymodule"
    )
    ```

* **Guard Clauses and Nil Checks:**
  * Do not add nil checks for values that cannot be nil in normal usage. A pointer type used to avoid copying (e.g. `*runner.Runner`) is not the same as an optional value — guarding it adds noise and misleading test coverage.
  * Before adding a defensive nil check, confirm that the value can actually be absent at runtime. If it cannot, the check should not exist.
* **Formatting:**
  * Adhere strictly to `go fmt` output. The `make fmt` command should be run before committing any changes. Consistent formatting is a non-negotiable aspect of code quality.
* **Naming Conventions:**
  * Follow Go's standard naming conventions:
    * `CamelCase` for exported identifiers (functions, types, variables accessible outside the package).
    * `camelCase` or `snake_case` for unexported identifiers (functions, types, variables within the package).
    * Acronyms (like `API`, `HTTP`, `URL`) should be all uppercase when multi-letter (`APIKey`), but `Id` becomes `ID`. This improves clarity and adherence to Go's idiomatic style.
    * Test files should end with `_test.go` (e.g., `my_module_test.go`).
* **Error Handling:**
  * Handle errors explicitly. Do not ignore errors. Every function returning an error must have its error return value checked.
  * Return errors as the last return value of a function.
  * Use `fmt.Errorf` for creating new errors and `errors.Wrap` or similar for adding context to existing errors (if a wrapping library like `pkg/errors` is used, check `go.mod` for common choices). This provides valuable debugging information.
  * Prefer sentinel errors for specific, anticipated error conditions, and custom error types for more complex error scenarios that carry additional data.
* **Types:**
  * Use meaningful type names that convey their purpose.
  * Define custom types for domain concepts when appropriate, rather than relying solely on built-in types. This improves type safety and code clarity.
* **Comments:**
  * Add comments sparingly. Focus on _why_ a particular piece of code exists or _what_ complex logic is trying to achieve, rather than just _how_ it works (which should be evident from the code itself).
  * Exported functions, variables, and types should have clear doc comments explaining their purpose, parameters, and return values.
* **Structs:**
  * Keep structs lean and focused on a single responsibility.
  * Tag fields with `json:`, `xml:`, `db:`, `yaml:`, etc., when serializing/deserializing or interacting with databases. This enables proper data mapping and external integration.
* **Functions:**
  * Keep functions small and focused on a single task. This improves testability and readability.
  * Avoid excessive parameters; consider grouping related parameters into a struct for better organization and maintainability.
* **Concurrency:**
  * Use Go's concurrency primitives (`goroutines`, `channels`, `sync` package) responsibly. Understand the implications of concurrent operations.
  * Avoid data races by using appropriate synchronization mechanisms (e.g., `sync.Mutex`, `sync.RWMutex`, channels).
  * Common patterns like worker pools, fan-out/fan-in, and context cancellation should be used where applicable.

### Project-Specific Considerations:

* **`pkg/` and `internal/`:**
  * `pkg/`: Contains reusable library code that can be imported by external applications. This includes shared domain models, utility functions, and interfaces.
  * `internal/`: Contains private application code that cannot be imported by other repositories. This typically houses application-specific logic, service implementations, and handlers.
  * Respect the architectural boundaries implied by these directories.
  * Within `internal/`, common subdirectories include `api` (for HTTP handlers and request/response models), `service` (for business logic), and `repository` (for database interactions).

### Dependency Management:

* **`go.mod` and `go.sum`:** These files manage the project's dependencies. Always use Go Modules for dependency management.
* **Updating Dependencies:** Use `go get -u` to update dependencies to their latest compatible versions. Ensure `go mod tidy` is run after making changes to dependencies to clean up unused modules.
* **Vendoring:** If vendoring is required (e.g., for air-gapped environments), use `go mod vendor`. Otherwise, rely on Go Modules to download dependencies during build.

## 3. Testing Conventions

* **Table-based unit tests:** Always write unit tests as table-driven tests. Each public function gets one `TestX` function containing a `tests` slice of anonymous structs. Each row has a descriptive `name` field that reads as a sentence (e.g., `"session create fails returns error"`). Private functions are not tested directly; their behaviour is covered through the public functions that call them.

* **Mocked dependencies:** Tests always use mocked dependencies (gomock) for every interface the unit under test depends on. Never use a real implementation (e.g., `session.InMemoryService()`) where an interface mock is available. Use `make mock` to generate mocks; add new mockgen commands to the `mock` target when introducing new mockable interfaces.

* **Explicit argument matchers:** When setting up mock expectations, always specify exact expected argument values. Do not use `gomock.Any()` except for `ctx`. Compute expected values independently from raw inputs — never call the production function under test to build them. If the same function appears on both sides of an assertion, any logical bug in it produces the same wrong output on both sides and the test passes silently. For example, populate a context struct directly in the test row rather than calling the helper that builds it.

* **Extract interface for non-mockable concrete dependencies:** When a struct field is a concrete external type (not an interface), extract a minimal interface covering only the methods actually called. Place the interface in its own file (e.g., `runner_port.go`) so mockgen can target it without regenerating unrelated mocks. The public constructor retains the concrete type as its parameter.

## 4. Design Principles

* **Immutability:** Construct objects fully at creation time. Pass all required values into the constructor rather than mutating fields after construction. This makes the object's state predictable at every point in the code and eliminates a class of bugs caused by partially initialised structs.

* **Layer Ownership:** Each layer owns its own concerns and validates them itself. Do not duplicate validation logic across layers, and do not write tests for behaviour that belongs to a different layer. When a handler delegates to a usecase, the usecase is responsible for domain validation — the handler should not replicate it just to make stub-based tests pass.
  
