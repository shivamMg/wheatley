package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
)

const (
	corvo = "abcdefghijklmnopqrstuvwxyz"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var tmplts = template.Must(template.ParseGlob("templates/*"))

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := tmplts.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	bind := ":3000"
	http.HandleFunc("/", indexPageHandler)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Connection made...")

		roomName := "attano"
		room, created := getRoom(roomName)
		if created {
			go room.run()
		}

		rand.Seed(time.Now().Unix())
		randName := corvo[:rand.Intn(26)]
		client := &Client{name: randName, send: make(chan []byte, 256), room: room, conn: conn}
		fmt.Println("Client created:", client.name)
		client.room.join <- client

		go client.writePump()
		client.readPump()

		//msg := Message{"shivam", "Test msg", time.Now().Format(time.UnixDate)}
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(bind, nil)
}
