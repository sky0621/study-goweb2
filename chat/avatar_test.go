package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	gomniauthtest "github.com/stretchr/gomniauth/test"

	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar

	// gomniauthcommon.Userのモックを生成して、chatUser構造体に渡す
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL) // エラーを返すモック
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.AvatarURL(testChatUser)

	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、AuthAvatar.AvatarURLはErrNoAvatarURLを返すべきです")
	}

	testURL := "http://url-to-avatar/"

	// gomniauthcommon.Userのモックを再生成して、chatUser構造体に渡す
	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testURL, nil) // テスト用のURLを返すモック
	url, err = authAvatar.AvatarURL(testChatUser)

	if err != nil {
		t.Error("値が存在する場合、AuthAvatar.AvatarURLはエラーを返すべきではありません")
	} else {
		if url != testURL {
			t.Error("AuthAvatar.AvatarURLは正しいURLを返すべきです")
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar

	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.AvatarURL(user)
	if err != nil {
		t.Error("GravatarAvatar.AvatarURLはエラーを返すべきではありません")
	}
	if url != "//www.gravatar.com/avatar/abc" {
		t.Errorf("GravatarAvatar.AvatarURLが%sという誤った値を返しました", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {

	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := fileSystemAvatar.AvatarURL(user)
	if err != nil {
		t.Error("FileSystemAvatar.AvatarURLはエラーを返すべきではありません")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.AvatarURLが%sという誤った値を返しました", url)
	}
}
