package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/sky0621/study-goweb2/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
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
	// HTTPリクエスト情報を渡したことでHTML上でリクエスト情報が参照可能になる
	t.templ.Execute(w, r) // XXX 本当は戻り値をチェックすべき
}

const cb = "/auth/callback/"

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() // コマンドラインで指定した -addr=":9999" 文字列から必要な情報を取得して *addr にセット

	var secKey = flag.String("secKey", "dummy", "セキュリティキー")
	log.Println(*secKey)
	gomniauth.SetSecurityKey(*secKey)

	var lhost = flag.String("host", "localhost", "ドメイン")
	baseURL := *lhost + *addr + cb
	log.Println(baseURL)
	gomniauth.WithProviders(
		facebook.New("a", "a", baseURL+"facebook"),
		github.New("a", "a", baseURL+"github"),
		google.New("a", "a", baseURL+"google"),
	)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))

	http.Handle("/login", &templateHandler{filename: "login.html"})

	http.HandleFunc("/auth/", loginHandler)

	r := newRoom()
	r.tracer = trace.New(os.Stdout) // コンソール出力
	http.Handle("/room", r)

	go r.run() // チャットルーム開始 -> 入退室やメッセージを待ち受ける

	log.Println("Webサーバーを開始します。ポート：", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
