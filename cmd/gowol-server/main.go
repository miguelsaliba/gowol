package main

import (
	"log"
	"net/http"
	"text/template"
)

func handler(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("tmpl/form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
