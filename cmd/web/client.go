package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)


type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager *Manager
	// egress используется для предотвращения одновременной записи в websocket-соединении
	egress chan Event
}


func NewClient(conn *websocket.Conn,manager *Manager) *Client{
	return &Client{
		// Коннктион одного пользователя
		connection: conn,
		// Менеджер хранящий всех подключенных пользователей
		manager: manager,
		egress: make(chan Event),
	}
}


// Го рутина запущенная до момента пока сообщение не будет отправленно пользователем
func (c *Client) readMessage(){

	
	// Срабатывает при выходе из вебсокета
	defer func(){
		// Удаляет клиента из нашего списка клиентов
		c.manager.removeClient(c)
	}()

	
	for {
		
		// Это возвращает тип сообщения (обычно текст или двоичные данные), содержимое сообщения и возможную ошибку.
		// ReadMessage является блокирующим вызовом. Это означает, что он будет ожидать до тех пор
		// , пока от клиента не будет получено сообщение, или не обнаружит закрытие соединения. Таким 
		// _ - выведет один в случае если мы отправили текст (все виды можно увидеть в исполнении ReadMessage вверху файла), 
		// payload - выведет сообщение которые мы отправим с фронта 
		_, payload, err := c.connection.ReadMessage()

		if err != nil{
			
			// Если есть ошибка цикл прерывается, и, таким образом, прерывается выполнение горутины.
			if websocket.IsUnexpectedCloseError(err,websocket.CloseGoingAway,websocket.CloseAbnormalClosure){
				log.Printf("error reading message: %v",err)
			}
			// break - может остановить работу данной го рутины и мы переместимся в defer определенный выше
			break
		}
		
		var request Event

		// Преобразуем входящий json в тип Event
		if err := json.Unmarshal(payload,&request); err != nil{
			log.Printf("error marshaling event :%v",err)
			break
		}

		// проверяет есть ли у нас определенный ключ Type в структуре Event по которому мы ищем обработчик пришедших данных
		if err := c.manager.routeEvent(request, c); err != nil{
			log.Println("error handling message: ", err)
		}

	}
}

func (c *Client) writeMessages(){
	
	defer func(){
		c.manager.removeClient(c)	
	}()
	
	
	for {
		select{
		// 	Канал будет выбирать сообщения и фактически один за другим отправлять их в вебсокет
		case message,ok := <- c.egress:
			
			// Если канал egress был закрыт (об этом свидетельствует !ok),
			// то возвращается сообщение закрытия через вебсокеты (websocket.CloseMessage) и функция завершается, 
			// что также приводит к вызову defer.
			if !ok{
				// 
				if err := c.connection.WriteMessage(websocket.CloseMessage,nil); err != nil{
					log.Println("connection closed:", err)
				}
				// Ломает цикл и запускает defer
				return
			}
			

			data,err := json.Marshal(message)

			if err != nil {
				log.Println(err)
				return
			}


			// Если получение сообщения прошло успешно, то оно отправляется через соединение с помощью WriteMessage обратно на фронт, 
			// после чего в консоль выводится сообщение о том, что сообщение было отправлено.
			if err := c.connection.WriteMessage(websocket.TextMessage,data); err != nil{
				log.Printf("failed to send Message %v",err)
			}
			
			log.Println("message was sent ")
		}
	}
}