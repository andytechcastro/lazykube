# Agent Guidelines for lazykube

This document provides instructions for agentic coding agents working on the **lazykube** project.

## Project Overview
**lazykube** is a Terminal User Interface (TUI) for Kubernetes, written in Go. It follows **Clean Architecture** principles to separate concerns, ensure testability, and manage complexity across multiple Kubernetes clusters.

## 🛠 Commands

### Build & Run
- **Build**: `go build -o lazykube cmd/lazykube/main.go`
- **Run**: `./lazykube`
- **Tidy dependencies**: `go mod tidy`

### Testing
- **Run all tests**: `go test ./...`
- **Run a single test**: `go test -v -run TestName ./path/to/pkg` (e.g., `go test -v -run TestCreatePod ./internal/infrastructure/k8s`)
- **Run tests with race detection**: `go test -race ./...` (Highly recommended for Interactors)

### Linting & Formatting
- **Format code**: `go fmt ./...`
- **Lint**: `golangci-lint run` (if installed)

## 🏗 Architecture & Style Guidelines

### Clean Architecture Layers
The project is strictly organized into layers following the Standard Go Project Layout:

1. **Domain** (`internal/domain/`): Core business logic and data structures (formerly entities). No dependencies on other layers.
2. **Usecases** (`internal/usecase/`): Application-specific business rules.
   - **Interactors**: Implement the business logic (found in `internal/usecase/`).
   - **Ports (Interfaces)**: Define interfaces for external data access (found in `internal/usecase/port/`).
3. **Adapters** (`internal/adapter/`):
   - **Controllers** (`internal/adapter/controller/`): Orchestrate the flow of data between the TUI and Usecases.
4. **Infrastructure** (`internal/infrastructure/`): Implementation details.
   - `k8s/`: Kubernetes API implementations of the ports.
   - `datastore/`: Kubernetes client initialization.
   - `tui/`: TView-based UI components and layout.
   - `config/`: Configuration handling.
5. **Registry** (`internal/registry/`): Dependency injection and wiring.

### Code Style
- **Formatting**: Always run `go fmt`. Use tabs for indentation.
- **Naming Conventions**:
  - Interfaces: `SomethingInteractor`, `SomethingGateway`.
  - Structs: `somethingInteractor` (private) or `SomethingInteractor` (public).
  - Receivers: Use short, descriptive names (e.g., `pi *podInteractor`, `pg *podGateway`).
  - Filenames: Use `snake_case.go`.
- **Imports**: Group imports as follows:
  1. Standard library
  2. Local project imports (`lazykube/internal/...`)
  3. External dependencies (`github.com/...`, `k8s.io/...`)
- **Error Handling**:
  - **NEVER** use `panic()` in production code. TUIs will crash and leave the terminal in a broken state.
  - Return `error` and wrap it using `fmt.Errorf("context: %w", err)`.
  - In TUI, use `resourceDict.ErrorModal.ShowError(err)` to display errors to the user instead of printing to stdout.

### Concurrency Patterns
- **CRITICAL**: Go maps are NOT thread-safe. When aggregating data from multiple clusters concurrently (common in Interactors), use `sync.Mutex` or channels.
- **Context Management**: Always pass and respect `context.Context`. Gateways should use the provided context for API calls to ensure requests can be cancelled (e.g., when the user switches views).
- **Goroutines**: Ensure goroutines used for streaming (logs, exec) are properly terminated using context cancellation or dedicated quit channels.

### Kubernetes Specifics
- **Mocking**: Use `k8s.io/client-go/kubernetes/fake` for testing Gateways.
- **Performance**: Avoid unnecessary API calls. Prefer `SharedInformers` for long-running watches to reduce load on the Kubernetes API server and improve UI responsiveness.
- **Client Casting**: Be careful when casting `kubernetes.Interface` to `*kubernetes.Clientset`. Only do it when necessary (e.g., for certain SPDY/exec operations) and check the `ok` value.

## 🖥 TUI Development (`internal/infrastructure/tui/`)
- **Framework**: Built with `github.com/rivo/tview` and `github.com/gdamore/tcell/v2`.
- **Global State**: Managed via `resourceDict` (defined in `dict_resources.go`). This object provides access to the `App`, `Pages`, and other UI components.
- **Navigation**:
  - Use `resourceDict.SetFocus(primitive)` to switch between components.
  - Modals are managed via `tview.Pages`. Use `pages.ShowPage(name)` and `pages.HidePage(name)`.
- **Component Design**: Each UI component (Table, Menu, View) should be encapsulated in its own file/struct and initialized via a `New...` function.

## ⚠️ Common Pitfalls & Anti-patterns to Avoid
1. **Inefficient Transformation**: Avoid using `json.Marshal`/`Unmarshal` to convert entities to maps (e.g., in `ResourceToData`). Use explicit DTOs or typed conversions.
2. **Eager Loading**: Don't initialize all Kubernetes clients/gateways at startup if the user has dozens of contexts. Use lazy loading where appropriate.
3. **Leaky Abstractions**: Keep Infrastructure details (like Kubernetes-specific types) out of the Entities and Usecases as much as possible.
4. **God Objects**: Avoid putting too much orchestration logic in `resourceDict`. Orchestration belongs in the **Controllers**.

## 📖 Recommended reading
- Clean Architecture (Robert C. Martin)
- Standard Go Project Layout
- Kubernetes client-go documentation
- tview documentation and examples
