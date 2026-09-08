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
	if r.Method == "POST" {
		name := r.FormValue("learnerName")

		fmt.Fprintf(w, "Hello %s, welcome to ThinkForge!", name)

	} else if r.Method == "GET" {

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
}

func main() {
	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/learn", learnHandler)

	fmt.Println("ThinkForge is running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
