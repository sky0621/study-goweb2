package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sky0621/study-goweb2/trace"
)

type room struct {
	forward chan []byte      // 他の全てのクライアントに転送するためのメッセージを保持するチャネル
	join    chan *client     // チャットルームに参加しようとしているクライアントのためのチャネル
	leave   chan *client     // チャットルームから退室しようとしているクライアントのためのチャネル
	clients map[*client]bool // 在室中の全てのクライアントを保持
	tracer  trace.Tracer     // チャットログを受け取るインタフェース「Tracer」
}

// 構造体「チャットルーム」の初期化用メソッド
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(), // デフォルトでは出力なしのトレーサーを設定
	}
}

// バックグラウンドでゴルーチン実行させる用
func (r *room) run() {
	for {
		select {
		case client := <-r.join: // 【入室】
			r.clients[client] = true
			r.tracer.Trace("新しいクライアントが参加しました")
		case client := <-r.leave: // 【退室】
			delete(r.clients, client) // 在室状態から消す
			close(client.send)        // 消したクライアントの送信用チャネルを閉じる
			r.tracer.Trace("クライアントが退室しました")
		case msg := <-r.forward: // 【全クライアントにメッセージ転送】
			r.tracer.Trace("メッセージを受信しました： ", string(msg))
			for client := range r.clients {
				select {
				case client.send <- msg: // １人１人のクライアントのチャネルにメッセージを流し込む
					r.tracer.Trace(" -- クライアントに送信されました")
				default: // 【転送失敗】
					delete(r.clients, client) // 在室状態から消す
					close(client.send)        // 消したクライアントの送信用チャネルを閉じる
					r.tracer.Trace(" -- 送信に失敗しました。クライアントをクリーンナップします。")
				}
			}
		}
	}
}

// room の参照(*room)を http.Handler 型に適合させる。（※同じシグネチャを持つ ServeHTTP メソッドを追加するだけ）

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// WebSocketを使うには、HTTP接続をアップグレードする必要がある
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

// HTTPリクエストが来るたびに呼ばれるメソッド
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil) // HTTP接続をアップグレードしてソケット生成
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	// 構造体「クライアント」を初期化
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize), // 指定の型のチャネルを指定バッファサイズで生成
		room:   r,
	}

	r.join <- client // 生成したクライアントをチャットルームの入室用チャネル（join）に投入！

	defer func() { r.leave <- client }() // クライアントの終了時に退室の処理を行う。ユーザがいなくなった際のクリーンナップ。

	go client.write() // ゴルーチン実行

	client.read() // メインスレッド上で実行。接続が保持されたまま、終了指示が出るまで他の処理をブロックする。
}
