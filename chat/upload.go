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
	log.Println("r.FormFile(avatarFile)")
	file, header, err := r.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer file.Close()

	log.Println("ioutil.ReadAll(file)")
	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	log.Println("filepath.Join(avatars, userID+filepath.Ext(header.Filename))")
	filename := filepath.Join("avatars", userID+filepath.Ext(header.Filename))
	log.Println(filename)
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, "成功")
}
