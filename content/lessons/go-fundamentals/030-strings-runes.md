# Strings and Runes (UTF-8)

To process text correctly in Go, you must understand the critical difference between a byte, a character, and a rune. 

Go handles text differently than languages like C or early Java. Go source code is strictly **UTF-8** encoded, and strings in Go are actually just read-only slices of bytes.

## 1. Strings are Byte Slices

A `string` is literally just an immutable `[]byte`. 
When you check the length of a string using `len()`, it returns the **number of bytes**, NOT the number of human-readable characters.

```go
func main() {
    s1 := "Hello"
    fmt.Println(len(s1)) // Prints 5 (5 characters, 5 bytes)

    s2 := "こんにちは" // "Hello" in Japanese
    fmt.Println(len(s2)) // Prints 15! (5 characters, but each is 3 bytes in UTF-8)
}
```

## 2. Enter the `rune`

Because UTF-8 characters can be anywhere from 1 to 4 bytes long, you cannot simply index a string like `s2[0]` to get the first character. `s2[0]` will only give you the first *byte* of the 3-byte Japanese character.

To represent a true Unicode character (a code point), Go provides the `rune` type (which is simply an alias for `int32`).

```mermaid
block-beta
  columns 3
  space:3
  A["String: '世' (World)"]:3
  block:bytes:3
    B["Byte 0\n(e4)"]
    C["Byte 1\n(b8)"]
    D["Byte 2\n(96)"]
  end
  space:3
  E["Rune: 19990"]:::runeBlock
  style E fill:#4CAF50,stroke:#333,stroke-width:2px;
```

## 3. Iterating Safely over Strings

If you use a standard `for i := 0` loop over a string, you will iterate over the raw bytes, which will corrupt multi-byte characters.

**❌ Bad (Iterating Bytes):**
```go
s := "Hi 世界"
for i := 0; i < len(s); i++ {
    fmt.Printf("%x ", s[i]) // Prints raw hex bytes, corrupts the characters
}
```

**✅ Good (Iterating Runes via `range`):**
When you use a `for...range` loop on a string, Go performs magic. It automatically decodes the UTF-8 bytes on the fly and returns valid `rune` values and their starting byte index!

```go
s := "Hi 世界"
for index, char := range s {
    // char is a rune!
    fmt.Printf("Index: %d, Char: %c\n", index, char)
}
// Output:
// Index: 0, Char: H
// Index: 1, Char: i
// Index: 2, Char:  
// Index: 3, Char: 世  (Notice index jumped from 2 to 3)
// Index: 6, Char: 界  (Notice index jumped from 3 to 6 because 世 is 3 bytes!)
```

To count the actual human-readable characters in a string, you must use the `utf8` package:
```go
import "unicode/utf8"

count := utf8.RuneCountInString("こんにちは") // Returns 5
```
