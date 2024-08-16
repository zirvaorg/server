package logic

import (
	"os/user"
)

func CheckPrivileges() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}

	if currentUser.Uid != "0" {
		return false
	}

	return true
}
