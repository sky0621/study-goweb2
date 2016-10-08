package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/sky0621/study-goweb2/config"
	"github.com/sky0621/study-goweb2/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
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
		},
	)

	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	// HTTPリクエスト情報を渡したことでHTML上でリクエスト情報が参照可能になる
	t.templ.Execute(w, data) // XXX 本当は戻り値をチェックすべき
}

func main() {
	cfg := config.ParseFlag() // アプリ引数の読み込み

	// サードパーティー認証用の設定（アプリキー（任意）と、とりあえずGoogleDeveloperConsole用）
	gomniauth.SetSecurityKey(cfg.SecKey) // 認証用にアプリのキー（任意）をセット
	gomniauth.WithProviders(
		google.New(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.Domain+cfg.Port+"/auth/callback/google"),
	)

	/*
	 * ルーティング
	 */
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"})) // 要認証
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout/", logoutHandler)
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploadHandler)

	// r := newRoom(UseAuthAvatar)
	r := newRoom(UseGravatar)
	r.tracer = trace.New(os.Stdout) // コンソール出力
	http.Handle("/room", r)

	go r.run() // チャットルーム開始 -> 入退室やメッセージを待ち受ける

	log.Println("Webサーバーを開始します。ポート：", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
