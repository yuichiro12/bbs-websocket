package main

import (
	"log"
	"net/http"
)

func main() {
	r := newRoom()
	http.Handle("/ws", r)
	go r.run()
	go r.listen()
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
