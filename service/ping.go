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
	IP           string  `json:"ip"`
	SuccessCount int     `json:"success_count"`
	RequestCount int     `json:"request_request"`
	MinRTT       float64 `json:"min_rtt"`
	AvgRTT       float64 `json:"avg_rtt"`
	MaxRTT       float64 `json:"max_rtt"`
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

	rtts := make([]time.Duration, 0, count)
	successCount := 0
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
			return PingResult{}, ctx.Err()
		default:
			reply := make([]byte, 1500)
			_, _, err := c.ReadFrom(reply)
			if err != nil {
				continue
			}
			rtts = append(rtts, time.Since(start))
			successCount++
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

	avgRTT := float64(totalRTT.Milliseconds()) / float64(successCount)
	if successCount == 0 {
		avgRTT = 0
	}

	return PingResult{
		IP:           resolvedIP,
		RequestCount: count,
		SuccessCount: successCount,
		MinRTT:       float64(minRTT.Milliseconds()),
		AvgRTT:       avgRTT,
		MaxRTT:       float64(maxRTT.Milliseconds()),
	}, nil
}
