package main

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[int]*client
	ids     map[*client]int
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		// idとclientの全単射を持っておく
		clients: make(map[int]*client),
		ids:     make(map[*client]int),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			id := client.id
			r.clients[id] = client
			r.ids[client] = id
		case client := <-r.leave:
			delete(r.clients, r.ids[client])
			delete(r.ids, client)
			close(client.send)
		case data := <-r.forward:
			d := jsonDecode(data)
			for _, id := range d.Ids {
				if r.clients[id] != nil {
					r.clients[id].send <- jsonEncode(d)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (r *room) listen() {
	file := getSockUrl()

	defer os.Remove(file)
	listener, err := net.Listen("unix", file)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}
	err = os.Chmod(file, 0777)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		go func() {
			defer conn.Close()
			data := make([]byte, 0)
			for {
				buf := make([]byte, 1024)
				nr, err := conn.Read(buf)
				if err != nil {
					if err != io.EOF {
						log.Printf("error: %v", err)
					}
					break
				}
				// 実際に読み込んだバイト以降のデータを除去したデータに変換
				buf = buf[:nr]
				// slice同士の結合は二つ目のsliceの後ろに...をつける
				data = append(data, buf...)
			}

			r.forward <- data
		}()
	}
}

func getSockUrl() string {
	data, err := ioutil.ReadFile(".sockUrl")
	if err != nil {
		log.Printf("ReadFile: %v", err)
	}
	return strings.TrimSpace(string(data))
}

// 各clientからのwss://のconnection要求に対して1度だけ発火
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("ServeHTTP: %v", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	client.getId()
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
