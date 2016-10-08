package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// チャットルームにつながった１人１人のクライアントを表す
type client struct {
	socket   *websocket.Conn
	send     chan *message // メッセージが送られるチャネル
	room     *room         // このクライアントが参加しているチャットルーム
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		msg.AvatarURL, _ = c.room.avatar.AvatarURL(c)
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	// 1byteずつ処理
	for msg := range c.send {
		// 1byteずつ WebScocket に流し込む
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}
