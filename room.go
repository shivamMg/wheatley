package main

import "fmt"

// Room is a room of clients
type Room struct {
	clients   map[*Client]bool
	join      chan *Client
	leave     chan *Client
	broadcast chan []byte
}

func newRoom() *Room {
	return &Room{
		clients:   make(map[*Client]bool),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan []byte),
	}
}

func (room *Room) run() {
	for {
		select {
		case client := <-room.join:
			fmt.Println("Client joined room:", client.name)
			room.clients[client] = true
		case client := <-room.leave:
			if _, ok := room.clients[client]; ok {
				fmt.Println("Client left room:", client.name)
				delete(room.clients, client)
				close(client.send)
			}
		case strMsg := <-room.broadcast:
			for client := range room.clients {
				select {
				case client.send <- strMsg:
				default:
					close(client.send)
					delete(room.clients, client)
				}
			}
		}
	}
}
