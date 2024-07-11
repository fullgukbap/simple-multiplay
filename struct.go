package main

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

type PlayerConn struct {
	player Player
	conn   *websocket.Conn
}

type Message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type PlayerMove struct {
	playerID  string
	direction string
}

type Send struct {
	target *websocket.Conn
	Msg    *Message
}

type ChannelSet struct {
	addPlayer    chan PlayerConn
	removePlayer chan *websocket.Conn
	movePlayer   chan PlayerMove
	broadcast    chan Message
	send         chan Send
}

func NewChannelSet() *ChannelSet {
	return &ChannelSet{
		addPlayer:    make(chan PlayerConn, 100),
		removePlayer: make(chan *websocket.Conn, 100),
		movePlayer:   make(chan PlayerMove, 100),
		broadcast:    make(chan Message, 100),
		send:         make(chan Send, 100),
	}
}
