package main

import (
	"fmt"
	"github.com/nathanmazzapica/tic-tac-two/internal/lobby"
	"html/template"
	"log"
	"net/http"
)

func main() {
	lobbyStore := lobby.NewStore()
	l := lobbyStore.New()
	fmt.Println(l)

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve HTML template
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
