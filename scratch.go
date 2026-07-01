package main

import (
	"html/template"
	"os"
)

func main() {
	tmpl := template.New("").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
	})
	
	// Create dummy files
	os.WriteFile("test1.html", []byte(`{{define "base"}}{{end}}`), 0644)
	os.WriteFile("test2.html", []byte(`{{define "leaderboard"}}{{add 1 2}}{{end}}`), 0644)
	defer os.Remove("test1.html")
	defer os.Remove("test2.html")

	var err error
	tmpl, err = tmpl.ParseGlob("test1.html")
	if err != nil {
		panic(err)
	}
	tmpl, err = tmpl.ParseGlob("test2.html")
	if err != nil {
		panic(err)
	}

	tmpl.ExecuteTemplate(os.Stdout, "test1.html", nil)
}
