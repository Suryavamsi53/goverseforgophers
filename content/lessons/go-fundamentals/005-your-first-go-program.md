# Your First Go Program

The tradition for every new programming language is to start by printing "Hello, World!" to the screen. In this lesson, we will write our first Go program, run it, and build it into a standalone executable.

## The Code

Create a new file named `main.go` and add the following code:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

## Code Breakdown

Let's understand exactly what every line of this program does:

1. `package main`: Every Go program is made up of packages. Programs start running in the package named `main`. This line tells the Go compiler that this file should compile as an executable program rather than a shared library.
2. `import "fmt"`: This tells the Go compiler to include the `fmt` (short for format) package from the Go standard library. This package contains functions for formatting text, including printing to the console.
3. `func main()`: The `main` function is special—it is the entry point of the executable program. When you run your code, execution starts here.
4. `fmt.Println("Hello, World!")`: This calls the `Println` function from the imported `fmt` package, passing it the string `"Hello, World!"`. It prints the string to the terminal, followed by a new line.

## Running the Program

During development, the easiest way to execute your Go code is using the `go run` command. It quickly compiles and runs your code in a single step without leaving a binary file behind.

Open your terminal and run:

```bash
$ go run main.go
Hello, World!
```

## Building an Executable

Unlike interpreted languages (like Python or JavaScript), Go is a **compiled language**. This means your source code is translated directly into machine code before it is executed.

To compile your code into a permanent, standalone binary executable, use the `go build` command:

```bash
$ go build main.go
```

This will generate an executable file in your current directory named `main` (or `main.exe` on Windows).

You can then run this executable directly, without needing the Go toolchain installed:

```bash
$ ./main
Hello, World!
```

Because Go binaries are statically linked by default, you can easily share this single binary file with other users on the same operating system, and they can run it without needing to install Go or any dependencies!
