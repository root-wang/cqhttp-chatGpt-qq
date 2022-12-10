package main

import (
	"cqhttp-client/src/log"
	msg "cqhttp-client/src/message"
	"flag"
	"github.com/gorilla/websocket"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var addr = flag.String("addr", "localhost:15733", "http service address")

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Infof("connecting to %s", u.String())

	var c *websocket.Conn
	var err error
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Error("close connection failed")
		}
	}(c)
	for {
		c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Info("connect failed")
			log.Info("after 20s reconnect the server")
			time.Sleep(20 * time.Second)
			continue
		}
		log.Info("connect to cqhttp succeed")

		client := msg.NewClient(c)
		go client.Run()
		go client.ReplyGroupMessage()

		for {
			var receiveMsg *msg.ReceiveMessage
			if err := c.ReadJSON(&receiveMsg); err != nil {
				log.Error(err.Error())
				continue
			}

			err := client.ReceiveMessage(receiveMsg)
			if err != nil {
				log.Error(err.Error())
				continue
			}
		}
	}
}
