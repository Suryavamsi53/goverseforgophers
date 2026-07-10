package main

import (
	"html/template"
	"log"
)

func main() {
	_, err := template.ParseFiles("ui/templates/layouts/base.html", "ui/templates/pages/practice_questions.html")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Templates parsed successfully!")
}
