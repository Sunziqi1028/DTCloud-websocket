package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var max = 10000

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	//	Path: "/echo"
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8001"}
	log.Println("connecting to ", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial err:", err)
		return
	}
	defer conn.Close()
	go timeWriter(conn)
	i := 0
	for {
		// 第一个包
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("read err:", err)
			return
		}
		log.Println("main ReadMessage read server mt:", mt, " message:", string(msg[:]))
		i++
	}
	// 阻塞
	// select {}
	time.Sleep(60 * time.Second)
	log.Println("client exit")
}

func timeWriter(conn *websocket.Conn) {
	var i = 0
	for {
		// 发第一个消息
		msg := &Account{Name: "第一个包 hello,张三", Age: i, Passwd: "123456"}
		jsonData, err := json.Marshal(msg)
		if err != nil {
			log.Println("client timeWriter Marshal err:", err, " msg:", msg)
			break
		}
		conn.WriteMessage(websocket.TextMessage, jsonData)
		// cpu阻塞下，等待读取完
		// time.Sleep(5 * time.Second)
		i++
		if i > max {
			break
		}
	}
}

type Account struct {
	Name   string
	Age    int
	Passwd string
}
