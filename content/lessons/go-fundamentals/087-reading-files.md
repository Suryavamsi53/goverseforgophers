# Reading Files

We've established that the `os` and `io` packages are the foundation of file interaction. Let's look at the three distinct architectural patterns for reading files, ranging from simple scripts to enterprise-grade memory management.

## Pattern 1: Slurping the Whole File (For Small Files)

If you are reading a configuration file (`config.json`), a simple SSH key, or anything under a few Megabytes, you can simply load the entire file into memory at once.

```go
import "os"

func main() {
    // os.ReadFile handles opening, reading, and closing the file automatically
    data, err := os.ReadFile("config.json")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(string(data))
}
```

## Pattern 2: Line-by-Line Streaming (For Text Logs)

If you have a 10GB server log file and you need to search for the word "Error", you cannot use `os.ReadFile` (your RAM will explode).

Instead, you use the `bufio.Scanner`. It reads the file chunk-by-chunk and yields one line of text at a time. Memory usage stays flat at a few kilobytes!

```go
import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    file, err := os.Open("server.log")
    if err != nil {
        panic(err)
    }
    defer file.Close() // ALWAYS defer close!

    scanner := bufio.NewScanner(file)
    
    // scanner.Scan() moves forward one line and returns true. 
    // It returns false when it reaches the End of File (EOF).
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "Error") {
            fmt.Println("Found Error:", line)
        }
    }
    
    // Check if the scanner crashed due to a bad read
    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading file:", err)
    }
}
```

## Pattern 3: Fixed Chunk Buffers (For Binary/Video Files)

If you are reading a binary file (like an `.mp4` video), there are no "lines" of text. A `bufio.Scanner` will crash because it searches for newline `\n` characters that don't exist.

To process huge binary files, you use a fixed-size byte buffer.

```go
func main() {
    file, _ := os.Open("video.mp4")
    defer file.Close()

    // Create a fixed 4KB buffer in memory
    buffer := make([]byte, 4096) 

    for {
        // file.Read fills our buffer and returns the number of bytes read
        bytesRead, err := file.Read(buffer)
        
        // If we hit the end of the file, break the loop
        if err == io.EOF {
            break
        }
        
        // Process the 4KB chunk (only process up to bytesRead, 
        // because the final chunk might be smaller than 4KB!)
        processBinaryData(buffer[:bytesRead])
    }
}
```
