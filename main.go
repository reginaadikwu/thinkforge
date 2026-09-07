package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, "Unable to load page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func learnHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/learn.html")

	if err != nil {
		http.Error(w, "Internal Server Error: Could not load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)

	if err != nil {
		http.Error(w, "Internal Server Error: Could not render template", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/learn", learnHandler)

	fmt.Println("ThinkForge is running on http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
