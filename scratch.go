package main

import (
	"fmt"
	"os"

	"github.com/russross/blackfriday/v2"
)

func main() {
	md := `**1.1** What is the zero value of an int in Go?
a) null b) 0 c) undefined d) Compile error
**Answer: b) 0**`
	
	html := blackfriday.Run([]byte(md))
	fmt.Println(string(html))
}
