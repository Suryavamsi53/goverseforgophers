package main

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
)

func main() {
	md := `**1.1** What is the zero value of an int in Go?
a) null b) 0 c) undefined d) Compile error
**Answer: b) 0**`

	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
}

