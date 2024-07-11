package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type GameState struct {
	players map[string]*Player
	conns   map[*websocket.Conn]string // conns key to player key
	ChannelSet
}

func NewGameState() *GameState {
	return &GameState{
		players:    make(map[string]*Player),
		conns:      make(map[*websocket.Conn]string),
		ChannelSet: *NewChannelSet(),
	}
}

func (g *GameState) Run(tickRate int) {
	ticker := time.NewTicker(time.Second / time.Duration(tickRate))
	defer ticker.Stop()

	for {
		select {
		case pc := <-g.addPlayer:
			g.players[pc.player.ID] = &pc.player
			g.conns[pc.conn] = pc.player.ID
			g.send <- Send{pc.conn, &Message{"your_id", pc.player.ID}}
			log.Printf("Player added: %v", pc.player) // Debug: Print added player

		case conn := <-g.removePlayer:
			if playerID, ok := g.conns[conn]; ok {
				delete(g.players, playerID)
				delete(g.conns, conn)
				g.broadcast <- Message{"player_left", playerID}
				log.Printf("Player removed: %v", playerID) // Debug: Print removed player
			}

		case msg := <-g.broadcast:
			for conn := range g.conns {
				if err := conn.WriteJSON(&Message{
					Type:    msg.Type,
					Payload: msg.Payload,
				}); err != nil {
					continue
				}
			}
			log.Printf("Broadcast message: %v", msg) // Debug: Print broadcast message

		case playerMove := <-g.movePlayer:
			if player, ok := g.players[playerMove.playerID]; ok {
				switch playerMove.direction {
				case "w":
					player.Y -= 5
				case "a":
					player.X -= 5
				case "s":
					player.Y += 5
				case "d":
					player.X += 5
				}
				log.Printf("Player moved: %v", player) // Debug: Print player movement
			}

		case data := <-g.send:
			data.target.WriteJSON(data.Msg)
			log.Printf("Sent message to player: %v", data.Msg) // Debug: Print sent message

		case <-ticker.C:
			g.broadcast <- Message{"update", g.players}
		}
	}

}
