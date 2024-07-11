package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	playerID := uuid.New().String()
	playerName := fmt.Sprintf("Player_%s", playerID[:6])
	player := Player{
		ID:   playerID,
		Name: playerName,
		X:    0,
		Y:    0,
	}

	gameState.addPlayer <- PlayerConn{
		player: player,
		conn:   conn,
	}

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		switch msg.Type {
		case "movement":
			gameState.movePlayer <- PlayerMove{playerID, msg.Payload.(string)}
		}
	}

	gameState.removePlayer <- conn

}
