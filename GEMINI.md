# Project Overview

`sopra` is a Go service that uses the OpenSky REST API to identify airplanes currently flying over a specific location. The service runs as a background process and authenticates with the OpenSky API headlessly. It returns a JSON response containing a list of aircraft within a configurable radius.

# Building and Running

**TODO:** There are no Go source files in the project yet. Once they are added, this section should be updated with instructions on how to build and run the service.

A typical Go project can be built and run with the following commands:

```sh
# Build the project
go build

# Run the project
./sopra
```

# Development Conventions

*   **Code Style:** Follow standard Go formatting and style guidelines. Use `gofmt` to format your code.
*   **Testing:** Add unit tests for new functionality in `_test.go` files.
*   **Dependencies:** Manage dependencies using Go modules (`go.mod`).
