package main

import "errors"

var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません。")

// ユーザのプロフィール画像を表すインタフェース
type Avatar interface {
	// 渡されたクライアントのアバターURLを取得（※取得できなかったらErrNoAvatarURLを返す）
	AvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) AvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}
