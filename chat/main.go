package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

// テンプレート管理用の構造体
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// テンプレートハンドラーをレシーバとするHTTPリクエスト処理関数
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(
		func() {
			t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
		})
	t.templ.Execute(w, nil) // XXX 本当は戻り値をチェックすべき
}

func main() {
	http.Handle("/", &templateHandler{filename: "chat.html"})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
