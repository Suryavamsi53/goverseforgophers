# Go Workspaces and Modules

When your code uses external packages, those packages (distributed as modules) become dependencies. Managing these dependencies effectively is crucial for building robust applications. 

This lesson is structured progressively from **Basic** dependency management to **Advanced** multi-module enterprise workspace configurations, using official patterns from the Go team.

---

## 🟢 Basic Level: Introduction to Modules

A Go module is a collection of related Go packages that are versioned together as a single unit. Modules record precise dependency requirements and create reproducible builds.

### 1. Creating a Module (`go mod init`)
To start a new project, you initialize a module. This creates a `go.mod` file in your current directory.

```bash
$ mkdir myproject
$ cd myproject
$ go mod init github.com/yourusername/myproject
```

### 2. The `go.mod` File
The `go.mod` file defines the module's path and its dependency requirements. It looks like this:

```go
module github.com/yourusername/myproject

go 1.22.0
```

### 3. Adding a Dependency (`go get`)
When you import an external package in your code and run your program, Go automatically downloads the dependency. However, you can explicitly add or update a dependency using `go get`:

```bash
$ go get golang.org/x/text
```

### 4. The `go.sum` File
When you add a dependency, Go creates a `go.sum` file. This file contains cryptographic hashes of the specific module versions you downloaded. This ensures that tomorrow, you (or your CI/CD system) will download the exact same, untampered code.

---

## 🟡 Intermediate Level: Managing Dependencies

As your project grows, your dependency tree will become more complex. Go provides built-in tooling to maintain a clean project state.

### 1. Cleaning up with `go mod tidy`
As you write code, you might add imports you end up not using, or remove code that relied on a third-party library. 
The `go mod tidy` command ensures your `go.mod` matches the source code in your module. It adds any missing modules necessary to build your current packages and removes unused modules that don't provide any relevant packages.

```bash
$ go mod tidy
```
*Best Practice: Always run `go mod tidy` before committing your code.*

### 2. Upgrading Dependencies
You can upgrade a specific dependency to its latest minor or patch release:
```bash
$ go get golang.org/x/text@latest
```
Or target a specific version:
```bash
$ go get golang.org/x/text@v0.14.0
```

### 3. Vendoring (`go mod vendor`)
In highly secure or offline environments (like enterprise CI/CD pipelines), relying on the internet to fetch dependencies during build time is risky. 
Running `go mod vendor` creates a `vendor/` directory containing all your dependencies' source code locally.

```bash
$ go mod vendor
```

---

## 🔴 Advanced Level: Multi-Module Workspaces

In large projects or enterprise environments, you often need to work on multiple interlocking modules simultaneously. Before Go 1.18, you had to use hacky `replace` directives in your `go.mod` file to point to local directories.

**Go Workspaces (`go.work`)** elegantly solve this problem by allowing you to work with multiple modules in local directories concurrently without modifying any `go.mod` files.

### 1. Initializing a Workspace
Imagine a repository where you have a backend server module and a shared library module.
In the root directory of your repository, initialize a workspace:

```bash
$ go work init
```
This generates a `go.work` file.

### 2. Adding Modules to the Workspace
You can add your local modules to the workspace using `go work use`:

```bash
$ go work use ./server
$ go work use ./shared-library
```

Your `go.work` file will now look like this:
```go
go 1.22.0

use (
    ./server
    ./shared-library
)
```

### 3. Why Workspaces Matter
When you run `go build` or `go run` inside the workspace, the Go compiler will automatically resolve dependencies across the local modules listed in `go.work`, rather than trying to fetch them from the internet. 

**Crucially, `go.work` files should NOT be committed to version control.** They are meant for your local, personal development environment, allowing you to fluidly edit multiple interconnected modules at once without breaking the `go.mod` files for the rest of your team.
