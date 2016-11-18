package main

import "fmt"

// Room is a room of clients
type Room struct {
	name      string
	clients   map[*Client]bool
	join      chan *Client
	leave     chan *Client
	broadcast chan []byte
}

// Rooms is a map of all rooms
var Rooms map[string]*Room

func init() {
	Rooms = make(map[string]*Room)
}

// Returns pointer to room with name as roomName
// If room was created it returns true
func getRoom(roomName string) (*Room, bool) {
	if room, ok := Rooms[roomName]; ok {
		fmt.Println("Using existing room", roomName)
		return room, false
	}

	room := &Room{
		name:      "attano",
		clients:   make(map[*Client]bool),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan []byte),
	}
	Rooms[roomName] = room
	fmt.Println("Creating new room", roomName)
	return room, true
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
