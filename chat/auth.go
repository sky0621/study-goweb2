package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	gomniauthcommon "github.com/stretchr/gomniauth/common"

	"github.com/stretchr/objx"
)

// ChatUser ...
type ChatUser interface {
	UniqueID() string
	AvatarURL() string // gomniauth/common.Userインタフェースでも定義されているため、chatUser構造体のUserフィールドに適切な値がセットされていれば実装したことになる
}

type chatUser struct {
	gomniauthcommon.User // 型の埋め込み（gomniauth/common.Userインタフェースを実装したことになる）
	uniqueID             string
}

// ChatUserインタフェースの１メソッドを実装
func (u *chatUser) UniqueID() string {
	return u.uniqueID
}

// 次のハンドラーを要素として持つ
type authHandler struct {
	next http.Handler
}

// 認証を上被せする用
func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		// 未認証 -> /loginにリダイレクト
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

// MustAuth ヘルパー関数（次に実行したいハンドラーを渡すと、認証ハンドルした後で、渡したハンドラーを呼ぶ、デコレータ―パターン）
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// サードパーティーへのログイン待ち受け用
// /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("認証プロバイダ(%s)の取得に失敗しました： %s", provider, err), http.StatusBadRequest)
			return
		}

		loginURL, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("認証プロバイダ(%s)におけるGetBeginAuthURLの呼び出し中にエラーが発生しました： %s", provider, err), http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "callback":
		// サードパーティプロバイダに応じた認証プロバイダを取得し、認証の完了からユーザ情報取得まで行う
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("認証プロバイダ(%s)の取得に失敗しました： %s", provider, err), http.StatusBadRequest)
			return
		}

		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(w, fmt.Sprintf("認証プロバイダ(%s)における認証を完了できませんでした： %s", provider, err), http.StatusBadRequest)
			return
		}

		user, err := provider.GetUser(creds)
		if err != nil {
			http.Error(w, fmt.Sprintf("認証プロバイダ(%s)におけるユーザの取得に失敗しました： %s", provider, err), http.StatusBadRequest)
			return
		}

		chatUser := &chatUser{User: user}
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Name()))
		chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))
		avatarURL, err := avatars.AvatarURL(chatUser)
		if err != nil {
			http.Error(w, fmt.Sprintf("認証プロバイダ(%s)におけるAvatarURL取得に失敗しました： %s", provider, err), http.StatusBadRequest)
			return
		}

		authCookieValue := objx.New(map[string]interface{}{
			"userid":     chatUser.uniqueID,
			"name":       user.Name(),
			"avatar_url": avatarURL,
		}).MustBase64()
		log.Println("[authCookieValue] " + authCookieValue)
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}
