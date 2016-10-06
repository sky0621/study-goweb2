package main

import "github.com/gorilla/websocket"

// チャットルームにつながった１人１人のクライアントを表す
type client struct {
	socket *websocket.Conn
	send   chan []byte // メッセージが送られるチャネル
	room   *room       // このクライアントが参加しているチャットルーム
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg // メッセージは即座にフォワードチャネルに送る
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	// 1byteずつ処理
	for msg := range c.send {
		// 1byteずつ WebScocket に流し込む
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break // エラーあったら、その時点で処理ストップ
		}
	}
	c.socket.Close()
}
