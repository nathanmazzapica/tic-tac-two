package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/nathanmazzapica/tic-tac-two/internal/lobby"
	"github.com/nathanmazzapica/tic-tac-two/internal/ws"
	"html/template"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	HandshakeTimeout: 10 * time.Second,
}

type Server struct {
	Store   *lobby.Store
	Context context.Context
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server := Server{
		Store:   lobby.NewStore(),
		Context: ctx,
	}

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

	http.HandleFunc("/ws", server.serveWebsocket)

	log.Println("listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	l := s.Store.New()

	fmt.Println(l.Sink())
	client := ws.NewClient("1", conn, ws.NewSink(l.Sink()), l.Subscribe("1"))
	go l.Run(s.Context)
	go client.Listen(s.Context)
	go client.Send(s.Context)
}
