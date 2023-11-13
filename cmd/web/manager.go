package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
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
	clients ClientList
	sync.RWMutex 
}



func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),

	}
}

// Подключение к вебсокету 
func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")

	// Обновление http соединения в websocket
	conn, err := websocketUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	// Добавление нового клиента при соединении
	client := NewClient(conn,m)
	
	// добавление новвого клиента в слайс который является свойством типа Manager
	m.addClient(client)
}

// Добавление клиента в менеджер
func (m *Manager) addClient(client *Client){
	// Если идет одновременно два соединения , блокирует одно из них 
	m.Lock()
	
	// Пропускает его и разблокирует обратно 
	defer m.Unlock()

	// добавление зашедшего клиента в мапу
	m.clients[client] = true
}

// Удаление клиента 
func (m *Manager) removeClient(client *Client){
	m.Lock()
	defer m.Unlock()

	if _,ok := m.clients[client]; ok {
		
		client.connection.Close()
		delete(m.clients,client)
	}
}

