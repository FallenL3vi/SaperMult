package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Upgrading connection to WebSocket
	c, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}

	defer func() {
		log.Printf("closing connection")
		c.Close()
	}()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error %s when reading message from client", err)
			return
		}
		//Websocket protocol supprots binary and text messages
		//Check if binary if it is than send a error message
		if mt == websocket.BinaryMessage {
			err = c.WriteMessage(websocket.TextMessage, []byte("server doesn't support binary messages"))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
			}
			return
		}

		log.Printf("Receive message %s", string(message))

		//If msg is not Start send a error message
		if strings.Trim(string(message), "\n") != "start" {
			err = c.WriteMessage(websocket.TextMessage, []byte("You did not say the magic word!"))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
				return
			}
			continue
		}

		log.Println("start responding to client...")

		i := 1

		for {
			response := fmt.Sprintf("Notification %d", i)
			err = c.WriteMessage(websocket.TextMessage, []byte(response))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
				return
			}
			i = i + 1
			time.Sleep(2 * time.Second)
		}
	}

}

func main() {
	webSocketHandler := webSocketHandler{
		//Use default upgarder from Gorilla
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				//origin := r.Header.Get("Origin")
				//return origin == "http://localhost:3000"

				//!!! DELETE IT ON PRODUCTION !!!
				return true
			},
		},
	}

	//Register handler
	http.Handle("/", webSocketHandler)
	log.Print("Starting server...")
	//Start server
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
