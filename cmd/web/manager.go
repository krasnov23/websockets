package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	// websocketUpgrader используется для преобразования входящих HTTP-запросов в постоянное соединение websocket
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")

	// Обновление http соединения в websocket
	conn, err := websocketUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	conn.Close()
}
