package service

import (
	"context"
	"errors"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"math/rand"
	"net"
	"time"
)

type PingResult struct {
	IP     string  `json:"ip"`
	MinRTT float64 `json:"min_rtt"`
	AvgRTT float64 `json:"avg_rtt"`
	MaxRTT float64 `json:"max_rtt"`
}

func Ping(ip string, count int) (PingResult, error) {
	if count <= 0 {
		return PingResult{}, errors.New("count should be greater than 0")
	}

	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return PingResult{}, errors.New("ICMP listen error")
	}
	defer c.Close()

	dst, err := net.ResolveIPAddr("ip4", ip)
	if err != nil {
		return PingResult{}, errors.New("IP resolve error")
	}

	rtts := make([]time.Duration, count)
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()

	for i := 0; i < count; i++ {
		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 1,
			Body: &icmp.Echo{
				ID:   rand.Intn(65000),
				Seq:  i,
				Data: []byte(EchoData),
			},
		}
		msgBytes, _ := msg.Marshal(nil)

		start := time.Now()
		if _, err := c.WriteTo(msgBytes, dst); err != nil {
			return PingResult{}, errors.New("ICMP write error")
		}

		c.SetReadDeadline(time.Now().Add(2 * time.Second))
		select {
		case <-ctx.Done():
			return PingResult{}, errors.New("ICMP timeout")
		default:
			_, _, err := c.ReadFrom(make([]byte, 1500))
			if err != nil {
				return PingResult{}, errors.New("ICMP read error")
			}
			rtts[i] = time.Since(start)
		}
	}

	minRTT, maxRTT, totalRTT := rtts[0], rtts[0], rtts[0]
	for _, rtt := range rtts[1:] {
		if rtt < minRTT {
			minRTT = rtt
		}
		if rtt > maxRTT {
			maxRTT = rtt
		}
		totalRTT += rtt
	}

	return PingResult{
		IP:     ip,
		MinRTT: float64(minRTT.Milliseconds()),
		AvgRTT: float64(totalRTT.Milliseconds()) / float64(count),
		MaxRTT: float64(maxRTT.Milliseconds()),
	}, nil
}
