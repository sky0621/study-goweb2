package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	file, header, err := r.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	filename := filepath.Join("avatars", userID+filepath.Ext(header.Filename))
	log.Println(filename)
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, "成功")
}
