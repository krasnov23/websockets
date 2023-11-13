package main

import (
	"log"
	"net/http"
)

func main() {
	setUpApi()

	// Сервер будет прослушиваться на данном порте
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setUpApi() {

	// Создание нового соединения
	newManager := NewManager()

	/*Внутри http.Handle указывается путь к корню сервера (/), который будет обслуживаться HTTP-файловым сервером.
	Файловый сервер будет использовать директорию ./frontend для расположения файлов, которые будет предоставлять.*/
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	// Соединяем с front end частью вызывая по этому адресу websocket
	http.HandleFunc("/ws", newManager.ServeWS)

}
