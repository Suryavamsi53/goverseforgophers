# The `path/filepath` Package

When working with files and directories, you might be tempted to concatenate strings manually: `folder + "/" + filename`. 

**Never do this.** 

If you compile your program for Windows, Windows uses backslashes `\` for paths, not forward slashes `/`. Your code will instantly break.

The `path/filepath` package handles OS-specific path separators and parsing automatically.

## 1. Path Construction (`Join`)

Always use `filepath.Join` to construct paths. It automatically adds the correct separator (`/` on Linux/Mac, `\` on Windows) and cleans up any messy double-slashes.

```go
import (
    "fmt"
    "path/filepath"
)

func main() {
    dir := "/var/log/"
    file := "app.log"
    
    // ❌ Bad: "/var/log//app.log"
    fmt.Println(dir + "/" + file) 
    
    // ✅ Good: "/var/log/app.log"
    fmt.Println(filepath.Join(dir, file)) 
}
```

## 2. Extracting Information

If you have a long absolute path, you can easily parse it:

```go
path := "/home/user/documents/report.pdf"

fmt.Println(filepath.Base(path)) // "report.pdf" (The final element)
fmt.Println(filepath.Dir(path))  // "/home/user/documents" (The directory)
fmt.Println(filepath.Ext(path))  // ".pdf" (The extension)
```

## 3. Walking Directories (`WalkDir`)

If you need to recursively scan an entire directory (e.g., to find all `.png` images inside a massive folder structure), you use `filepath.WalkDir`.

You pass it a root directory and an anonymous callback function. The package will recursively dive into every subfolder and execute your function for every single file it finds.

```go
func main() {
    root := "./images"

    err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err // Handle permission errors
        }
        
        // Skip directories
        if d.IsDir() {
            return nil 
        }
        
        // Process only PNG files
        if filepath.Ext(path) == ".png" {
            fmt.Println("Found image:", path)
        }
        
        return nil
    })

    if err != nil {
        fmt.Println("Error walking directory:", err)
    }
}
```

### Performance Insight: `Walk` vs `WalkDir`
In older Go codebases, you will see `filepath.Walk`. **This is deprecated for performance reasons.** 
`filepath.Walk` executed an `os.Stat` system call (which requires reading the hard drive) on every single file to gather metadata, making it incredibly slow for large directories. 
`filepath.WalkDir` (introduced in Go 1.16) uses highly optimized directory reads that skip the heavy `Stat` calls entirely, making it massively faster!
