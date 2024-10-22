package service

import (
	"context"
	"math/rand"
	"net"
	"server/internal/utils"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	maxHops    = 64
	timeOut    = 30 * time.Second
	timeOutHop = 3 * time.Second
	numProbes  = 3
)

type TracerouteHop struct {
	TTL      int             `json:"ttl"`
	Addr     string          `json:"addr"`
	Hostname string          `json:"hostname"`
	RTTs     []time.Duration `json:"rtts"`
}

type TracerouteResult struct {
	Hops []TracerouteHop `json:"hops"`
}

func Traceroute(address string) (TracerouteResult, error) {
	resolvedIP, err := utils.ResolveIPWithOutPort(address)
	if err != nil {
		return TracerouteResult{}, err
	}

	ipAddr, err := net.ResolveIPAddr("ip4", resolvedIP)
	if err != nil {
		return TracerouteResult{}, err
	}

	c, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return TracerouteResult{}, err
	}
	defer c.Close()

	p := ipv4.NewPacketConn(c)

	var result TracerouteResult
	targetReached := false

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	for ttl := 1; ttl <= maxHops && !targetReached; ttl++ {
		var RTTs []time.Duration
		var addr net.Addr
		for probe := 0; probe < numProbes; probe++ {
			select {
			case <-ctx.Done():
				return result, ctx.Err()
			default:
				start := time.Now()
				if err := p.SetTTL(ttl); err != nil {
					RTTs = append(RTTs, -1)
					continue
				}

				msg := icmp.Message{
					Type: ipv4.ICMPTypeEcho,
					Code: 0,
					Body: &icmp.Echo{
						ID:   rand.Intn(65000),
						Seq:  probe,
						Data: []byte(EchoData),
					},
				}
				msgBytes, err := msg.Marshal(nil)
				if err != nil {
					RTTs = append(RTTs, -1)
					continue
				}

				if _, err := p.WriteTo(msgBytes, nil, ipAddr); err != nil {
					RTTs = append(RTTs, -1)
					continue
				}

				rb := make([]byte, 1500)
				if err := c.SetReadDeadline(time.Now().Add(timeOutHop)); err != nil {
					RTTs = append(RTTs, -1)
					continue
				}
				n, a, err := c.ReadFrom(rb)
				if err != nil {
					RTTs = append(RTTs, -1)
					continue
				}
				addr = a

				rm, err := icmp.ParseMessage(1, rb[:n])
				if err != nil {
					RTTs = append(RTTs, -1)
					continue
				}

				duration := time.Since(start)
				switch rm.Type {
				case ipv4.ICMPTypeTimeExceeded, ipv4.ICMPTypeEchoReply:
					RTTs = append(RTTs, duration)
					if rm.Type == ipv4.ICMPTypeEchoReply {
						targetReached = true
					}
				default:
					RTTs = append(RTTs, -1)
				}
			}
		}

		for len(RTTs) < numProbes {
			RTTs = append(RTTs, -1)
		}

		hopAddr := "*"
		hostname := "*"
		if addr != nil {
			hopAddr = addr.String()
			hop := TracerouteHop{TTL: ttl, Addr: hopAddr, RTTs: RTTs}
			hop.Hostname = utils.ResolveHostname(hopAddr)
			result.Hops = append(result.Hops, hop)
		} else {
			result.Hops = append(result.Hops, TracerouteHop{
				TTL:      ttl,
				Addr:     hopAddr,
				Hostname: hostname,
				RTTs:     RTTs,
			})
		}
	}

	return result, nil
}
