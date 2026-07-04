# Writing Files

Writing to files follows the same three architectural patterns as reading. 

## Pattern 1: Dumping the Whole File

If you have a small string or byte slice in memory and just want to save it to disk, use `os.WriteFile`. It creates the file (or truncates it if it exists), writes the data, sets the permissions, and closes the file in one shot.

```go
import "os"

func main() {
    data := []byte("Hello, World!")
    // 0644 means standard Unix permissions (Read/Write for owner, Read for others)
    err := os.WriteFile("output.txt", data, 0644) 
    if err != nil {
        panic(err)
    }
}
```

## Pattern 2: Appending and `os.OpenFile`

If you are building a custom logger, you do not want to overwrite the file! You want to open the file in **Append Mode**.

`os.Create` and `os.Open` are just convenience wrappers. For advanced control, use `os.OpenFile`, which allows you to pass specific bitmask flags to the Operating System.

```go
func appendLog(message string) {
    // os.O_APPEND: Add to the end of the file
    // os.O_CREATE: Create the file if it doesn't exist
    // os.O_WRONLY: Open in Write-Only mode
    flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
    
    file, err := os.OpenFile("system.log", flags, 0644)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    file.WriteString(message + "\n")
}
```

## Pattern 3: Buffered Writing (`bufio.Writer`)

Every time you call `file.Write()`, Go executes a "Syscall" to wake up the Operating System kernel and instruct the hard drive to spin up and write data. Syscalls are incredibly slow.

If you are writing 100,000 tiny strings in a loop, hitting the hard drive 100,000 times will bring your application to a crawl.

Instead, wrap your file in a `bufio.Writer`. It collects all your tiny writes in a block of RAM. When the RAM chunk gets full (usually 4KB), it flushes it to the hard drive in a single, massive Syscall.

```go
import (
    "bufio"
    "os"
)

func main() {
    file, _ := os.Create("massive_data.txt")
    defer file.Close()

    // Wrap the file in a buffer
    writer := bufio.NewWriter(file)

    for i := 0; i < 100000; i++ {
        // This writes instantly to RAM, NOT the hard drive
        writer.WriteString("Some data line\n") 
    }

    // CRITICAL: You MUST call Flush() when done! 
    // Otherwise, the final chunk of data remaining in RAM will be lost forever.
    writer.Flush() 
}
```

## Advanced Storage Safety: `fsync`
Even if you call `writer.Flush()`, the Operating System might hold the data in *its own* memory cache before writing to the physical metal platter. If the server loses power, data is corrupted. For ultra-critical database engines, developers call `file.Sync()` (which executes an `fsync` syscall) to force the hardware controller to physically write to disk immediately.
