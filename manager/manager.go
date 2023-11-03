package manager

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	// websocketUpgrader используется для преобразования входящих HTTP-запросов в постоянное соединение websocket
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
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
		fmt.Println(err)
		return
	}

	conn.Close()
}
