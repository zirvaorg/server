package logic

import (
	"errors"
	"net"
	"os/user"
	"server/internal/msg"
)

func CheckEnvironment(p string) error {
	if !userIsRoot() {
		return errors.New(msg.PrivilegesErr)
	}

	if !portIsAvailable(p) {
		return errors.New(msg.ServerPortInUse)
	}

	return nil
}

func portIsAvailable(port string) bool {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}

	listener.Close()

	return true
}

func userIsRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}

	if currentUser.Uid != "0" {
		return false
	}

	return true
}
