# Introduction to Go

## What is Go?

The Go programming language is an open source project to make programmers more productive.

Go is expressive, concise, clean, and efficient. Its concurrency mechanisms make it easy to write programs that get the most out of multicore and networked machines, while its novel type system enables flexible and modular program construction. Go compiles quickly to machine code yet has the convenience of garbage collection and the power of run-time reflection. It's a fast, statically typed, compiled language that feels like a dynamically typed, interpreted language.

## Prerequisites

Before getting started, you will need:
* **Programming experience.** The code here assumes that you have read some code and have a general understanding of programming.
* **A tool to edit your code.** Any text editor you have will work fine. Most text editors have good support for Go. The most popular are VSCode (free), GoLand (paid), and Vim (free).
* **A command terminal.** Go works well using any terminal on Linux and Mac, and on PowerShell or cmd in Windows.

## Getting Started: Hello, World

The traditional first program in any language is to print "Hello, World". In Go, it looks like this:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Breaking it down

1. **`package main`**: Packages are Go's way of grouping functions, and it's made up of all the files in the same directory. A package named `main` tells the Go compiler that the package should compile as an executable program instead of a shared library.
2. **`import "fmt"`**: This line imports the `fmt` package, which contains functions for formatting text, including printing to the console. This package is one of the standard library packages you get when you install Go.
3. **`func main()`**: A `main` function executes by default when you run the main package. This is the entry point of your program.
4. **`fmt.Println(...)`**: A call to a function from the `fmt` package to print the text to the screen.

## Running Go Code

Once you have written your code into a file named `hello.go`, you can run it using the `go` tool provided by the Go installation.

Open a terminal and run the following command:

```bash
$ go run hello.go
Hello, World!
```

The `go run` command compiles and runs the Go code in a single step. If you want to compile the code into an executable file without running it immediately, you can use the `go build` command:

```bash
$ go build hello.go
```

This will create an executable file named `hello` (or `hello.exe` on Windows) in the current directory. You can then run the executable directly:

```bash
$ ./hello
Hello, World!
```

## Why Go?

Go was designed at Google in 2007 to improve programming productivity in an era of multicore, networked machines and large codebases. The designers wanted to resolve common criticisms of other languages in use at Google, but keep their useful characteristics:

* Static typing and run-time efficiency (like C)
* Readability and usability (like Python or JavaScript)
* High-performance networking and multiprocessing

This makes Go an excellent choice for modern backend development, cloud infrastructure, and microservices architecture.
