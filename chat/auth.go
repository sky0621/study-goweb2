package main

import "net/http"

// 次のハンドラーを要素として持つ
type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 未認証 -> /loginにリダイレクト
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

// ヘルパー関数（次に実行したいハンドラーを渡すと、認証ハンドルした後で、渡したハンドラーを呼ぶ、デコレータ―パターン）
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
