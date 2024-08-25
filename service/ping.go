package service

import (
	"context"
	"errors"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"math/rand"
	"net"
	"server/internal/utils"
	"time"
)

type PingResult struct {
	IP     string  `json:"ip"`
	MinRTT float64 `json:"min_rtt"`
	AvgRTT float64 `json:"avg_rtt"`
	MaxRTT float64 `json:"max_rtt"`
}

func Ping(domain string, count int) (PingResult, error) {
	if count <= 0 {
		return PingResult{}, errors.New("count should be greater than 0")
	}

	resolvedIP, err := utils.ResolveIPWithOutPort(domain)
	if err != nil {
		return PingResult{}, err
	}

	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return PingResult{}, err
	}
	defer c.Close()

	dst, err := net.ResolveIPAddr("ip4", resolvedIP)
	if err != nil {
		return PingResult{}, err
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
			return PingResult{}, err
		}

		c.SetReadDeadline(time.Now().Add(2 * time.Second))
		select {
		case <-ctx.Done():
			return PingResult{}, err
		default:
			_, _, err := c.ReadFrom(make([]byte, 1500))
			if err != nil {
				return PingResult{}, err
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
		IP:     resolvedIP,
		MinRTT: float64(minRTT.Milliseconds()),
		AvgRTT: float64(totalRTT.Milliseconds()) / float64(count),
		MaxRTT: float64(maxRTT.Milliseconds()),
	}, nil
}
