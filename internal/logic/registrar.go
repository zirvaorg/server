package logic

import (
	"crypto/rand"
	"fmt"
	"os"
)

const tokenFilePath = "/usr/bin/zirva-app/.zirva_auth"

var TempRegistrarToken string

func GenerateRegistrarToken() string {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		panic(err)
	}

	TempRegistrarToken = fmt.Sprintf("%s-%x", "zirva", token)

	return TempRegistrarToken
}

func Registrar(token string) (bool, error) {
	file, err := os.Create(tokenFilePath)
	if err != nil {
		return false, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	_, err = file.WriteString(token)
	if err != nil {
		return false, err
	}

	TempRegistrarToken = ""

	return true, nil
}
