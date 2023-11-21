package main

import "encoding/json"


// Преобразование будет выполнено только тогда, когда вы явно декодируете RawMessage. Это полезно, 
// когда структура JSON-документа заранее неизвестна или может измениться в зависимости от контекста или других полей структуры.
type Event struct {
	Type    string `json:"type"`
	// поскольку входящие данные с форм, мы не знаем какой формат будет определенн мы вводим этот тип данных 
	// например входящие данные от сообщения могут отличаться от данных изменения комнаты
	Payload json.RawMessage `json:"payload"` 
}

// В такой манере будут описаны все функции для взаимодействия с нашими вебсокетами
// Данная функция может принимать разные входящие данные и их форматы поэтому мы не указываем здесь конкретную реализацию
// Мы заворачиваем данную логику в тип EventHandler, чтобы представить функцию в качестве значения
// , которое можно передавать и использовать в качестве аргумента или возвращаемого значения в других функциях.
// по сути по аналогии похоже на интерфейс
type EventHandler func (event Event, c *Client) error

const (
	EventSendMessage = "send_message"
)

// Мы будем ожидать что всякие раз когда срабатывает EventSendMessage
// Мы ожидаем входящее сообщение и хотим знать кто его отправил
type SendMessageEvent struct{
	Message string `json:"message"`
	From    string `json:"from"`
}

