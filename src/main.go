package main

import (
	client2 "cqhttp-client/src/client"
	"cqhttp-client/src/constant"
	"cqhttp-client/src/log"
	"cqhttp-client/src/log/file"
	msg "cqhttp-client/src/message"
	"cqhttp-client/src/module/gpt/API"
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
			panic("close connection failed")
		}
	}(c)

	for {
		c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Errorf("connect failed :%v", err)
			log.Info("after 20s reconnect the server")
			time.Sleep(20 * time.Second)
			continue
		}

		log.Info("connect to cqhttp succeed")

		// token := ""

		client := client2.NewClient(c)
		// client.AddModule("gpt", gpt.NewGpt(token))
		api := API.InitApi()
		client.AddModule(constant.GPTTEXT, api.APIByName(API.TextCompletion))
		client.AddModule(constant.GPTIMAGE, api.APIByName(API.ImageGeneration))
		client.AddModule(constant.GPTCHAT, api.APIByName(API.ChatCompletion))
		go client.Run()
		go client.ReplyGroupMessage()

	receiveLoop:
		for {
			var receiveMsg *msg.ReceiveMessage
			if err := c.ReadJSON(&receiveMsg); err != nil {
				log.Errorf("read receive message 2 JSON failed: %v", err)
				break receiveLoop
			}

			go func(receiveMsg *msg.ReceiveMessage) {
				err := client.ReceiveMessage(receiveMsg)
				if err != nil {
					log.Errorf("handle receive message failed: %v", err)
				}
			}(receiveMsg)
		}
	}
}

func init() {
	println(file.Prefix)
}
