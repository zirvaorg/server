package service

import (
	"net"
	"server/internal/utils"
	"time"
)

type ConnectionResult struct {
	IP           string  `json:"ip,omitempty"`
	ResponseTime float64 `json:"response_time,omitempty"`
}

func TcpOrUdp(p string, op string) (ConnectionResult, error) {
	start := time.Now()

	resolvedAddr, err := utils.ResolveIPWithPort(p)
	if err != nil {
		return ConnectionResult{
			IP:           resolvedAddr,
			ResponseTime: time.Since(start).Seconds(),
		}, err
	}

	var conn net.Conn
	if op == "tcp" {
		conn, err = net.DialTimeout("tcp", resolvedAddr, TimeOut)
	} else if op == "udp" {
		conn, err = net.DialTimeout("udp", resolvedAddr, TimeOut)
	} else {
		return ConnectionResult{
			IP:           resolvedAddr,
			ResponseTime: time.Since(start).Seconds(),
		}, err
	}

	if err != nil {
		return ConnectionResult{
			IP:           resolvedAddr,
			ResponseTime: time.Since(start).Seconds(),
		}, err
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	return ConnectionResult{
		IP:           resolvedAddr,
		ResponseTime: time.Since(start).Seconds(),
	}, nil
}
