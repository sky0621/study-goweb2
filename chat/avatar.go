package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

// ErrNoAvatarURL ...
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません。")

// Avatar ユーザのプロフィール画像を表すインタフェース
type Avatar interface {
	AvatarURL(u *chatUser) (string, error)
}

// TryAvatars ...
type TryAvatars []Avatar

func (a TryAvatars) AvatarURL(u *chatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.AvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

// AuthAvatar ...
type AuthAvatar struct{}

// UseAuthAvatar ...
var UseAuthAvatar AuthAvatar

// AvatarURL ...
func (_ AuthAvatar) AvatarURL(u *chatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

// GravatarAvatar ...
type GravatarAvatar struct{}

// UseGravatar ...
var UseGravatar GravatarAvatar

// AvatarURL ...
func (_ GravatarAvatar) AvatarURL(u *chatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

// FileSystemAvatar ...
type FileSystemAvatar struct{}

// UseFileSystemAvatar ...
var UseFileSystemAvatar FileSystemAvatar

// AvatarURL ...
func (_ FileSystemAvatar) AvatarURL(u *chatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
