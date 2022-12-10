package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	msg "mirai-go/src/message"
	"net/url"
	"os"
	"os/signal"
)

var addr = flag.String("addr", "localhost:15733", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	client := msg.NewClient(c)
	go client.Run()
	go client.SendGroupMessage()

	for {
		var receiveMsg *msg.ReceiveMessage
		if err := c.ReadJSON(&receiveMsg); err != nil {
			return
		}

		err := client.ReceiveMessage(receiveMsg)
		if err != nil {
			return
		}
	}

}
