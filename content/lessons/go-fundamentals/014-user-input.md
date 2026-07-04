# User Input

Building interactive CLI (Command Line Interface) applications requires reading input from the user. Go provides several ways to capture data from standard input (`os.Stdin`).

## 1. Using `fmt.Scan`

The simplest way to read user input is using the `fmt` package. 
* `fmt.Scan`: Reads space-separated values into variables.
* `fmt.Scanln`: Similar to `Scan`, but stops scanning at a newline.
* `fmt.Scanf`: Reads formatted text.

To use them, you pass the memory address (using the `&` pointer operator) of the variable where you want to store the input.

```go
package main

import "fmt"

func main() {
    var name string
    var age int

    fmt.Print("Enter your name and age: ")
    // Expects input like: Alice 25
    fmt.Scan(&name, &age) 
    
    fmt.Printf("Hello %s, you are %d years old.\n", name, age)
}
```
*Limitation: `fmt.Scan` stops at spaces. If a user types "John Doe", `name` will only capture "John".*

## 2. Reading Full Lines with `bufio.Reader`

If you need to read an entire sentence (including spaces), use a buffered reader from the `bufio` package attached to standard input (`os.Stdin`).

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    
    fmt.Print("Enter your full name: ")
    
    // ReadString reads until it encounters the specified byte (newline)
    fullName, _ := reader.ReadString('\n')
    
    // Trim the trailing newline character
    fullName = strings.TrimSpace(fullName)
    
    fmt.Printf("Welcome, %s!\n", fullName)
}
```

## 3. Reading Streams with `bufio.Scanner`

For continuous reading (like building a REPL or processing a text file line-by-line), `bufio.Scanner` is the most idiomatic and efficient tool.

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println("Type something (or 'exit' to quit):")
    
    // scanner.Scan() waits for the next token (default is line)
    for scanner.Scan() {
        text := scanner.Text()
        if text == "exit" {
            break
        }
        fmt.Println("You typed:", text)
    }
}
```
