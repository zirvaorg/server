package logic

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

var AuthToken string

func CheckAuthFile() bool {
	if _, err := os.Stat(tokenFilePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func readAuthToken() string {
	if AuthToken != "" {
		return AuthToken
	}

	file, err := os.Open(tokenFilePath)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	buf := new(strings.Builder)
	_, err = io.Copy(buf, file)
	if err != nil {
		panic(err)
	}

	AuthToken = strings.TrimSpace(buf.String())
	return AuthToken
}

func CheckAuth(r *http.Request) (bool, error) {
	if !CheckAuthFile() {
		return false, errors.New("auth token file not found")
	}

	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return false, errors.New("authorization header is empty")
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 {
		return false, errors.New("authorization header is invalid")
	}

	if authHeaderParts[1] != readAuthToken() {
		return false, errors.New("authorization header is invalid")
	}

	return true, nil
}
