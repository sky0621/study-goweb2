package main

type room struct {
	forward chan []byte      // 他の全てのクライアントに転送するためのメッセージを保持するチャネル
	join    chan *client     // チャットルームに参加しようとしているクライアントのためのチャネル
	leave   chan *client     // チャットルームから退室しようとしているクライアントのためのチャネル
	clients map[*client]bool // 在室中の全てのクライアントを保持
}

// バックグラウンドでゴルーチン実行させる用
func (r *room) run() {
	for {
		select {
		case client := <-r.join: // 【入室】
			r.clients[client] = true
		case client := <-r.leave: // 【退室】
			delete(r.clients, client) // 在室状態から消す
			close(client.send)        // 消したクライアントの送信用チャネルを閉じる
		case msg := <-r.forward: // 【全クライアントにメッセージ転送】
			for client := range r.clients {
				select {
				case client.send <- msg: // １人１人のクライアントのチャネルにメッセージを流し込む
				default: // 【転送失敗】
					delete(r.clients, client) // 在室状態から消す
					close(client.send)        // 消したクライアントの送信用チャネルを閉じる
				}
			}
		}
	}
}
