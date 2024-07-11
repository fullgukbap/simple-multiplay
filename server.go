package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	gameState = NewGameState()
)

const (
	tickRate = 160
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./views/index.html")
	})

	r.Get("/ws", handleConnection)

	go http.ListenAndServe(":8080", r)
	gameState.Run(tickRate)
}
