package main

import (
	"errors"
	"fmt"
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
	handlers map[string]EventHandler
	// без sync.RWMutex не будут работать методы Lock и Unlock в addClient и RemoveClient
	// sync.RWMutex в структуре Manager используется для безопасного состояния clients в многопоточной среде.
	// Это служит для исключения условий гонки (race condition), которые могут возникнуть,
	// когда два или более потока пытаются одновременно прочитать и записать данные в clients.
	// RWMutex разрешает множественные потоки считывать данные одновременно (что может ускорить производительность за счет распараллеливания чтения),
	// но только один поток может записывать данные (что исключает возможность конфликта данных).
}


func NewManager() *Manager {
	m := &Manager{
		clients: make(ClientList),
		handlers: make(map[string]EventHandler),
	}

	// Установка все возможных обработчиков в менеджере (обработчики срабатывают взависимости от входящих данных)
	m.setupEventHandlers()

	// возвращает manager с установленными обработчиками событий взависимости от Event.type
	return m
}

func (m *Manager) routeEvent(event Event,c *Client) error{
	
	// Если значение по ключу типа в Event и если мы нашли event среди указанных нами, то он исходя от него будет выполненна 
	// определенная логика 
	if handler,ok := m.handlers[event.Type]; ok{
		// Если функция EventHandler возвращает ошибку при передаваемых event и client
		if err := handler(event,c); err != nil {
			return err 
		}
		return nil
	}else{
		return errors.New("there is no such event type")
	}
}

// Добавление обработчика в наш менеджер 
func (m *Manager) setupEventHandlers(){
	// под обработчик с ключом send_message вызывается функция SendMessage,
	// EventSendMessage в данном случае константа в файле events.go
	m.handlers[EventSendMessage] = SendMessage
}

func SendMessage(event Event,c *Client) error {
	fmt.Println(event)
	return nil
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

	// Запускает горутину которая обрабатывает форму с сообщением
	go client.readMessage()

	// Запускаем горутину которая выводит сообщение
	go client.writeMessages()
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

	// Если ok true то есть если найден в мапе, закрывает соединение и удаляет клиента из мапы
	if _,ok := m.clients[client]; ok {
		
		client.connection.Close()
		delete(m.clients,client)
	}
}


