package main

import "time"

// 1メッセージ
type message struct {
	Name    string
	Message string
	When    time.Time
}
