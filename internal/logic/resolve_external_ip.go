package logic

import (
	"io"
	"net"
	"net/http"
)

func ResolveExternalIP() string {
	urls := []string{
		"https://api.ipify.org",
		"https://api.seeip.org",
		"https://ifconfig.co/ip",
		"https://ifconfig.io/ip",
		"https://icanhazip.com",
		"https://ident.me",
		"https://myexternalip.com/raw",
		"https://wtfismyip.com/text",
		"https://ipinfo.io/ip",
		"https://ipecho.net/plain",
		"https://icanhazip.com",
		"https://ifconfig.me/ip",
		"https://ip.tyk.nu",
		"https://ipecho.net/plain",
		"https://l2.io/ip",
		"https://wtfismyip.com/text",
		"https://icanhazip.com",
		"https://ifconfig.me/ip",
		"https://ip.seeip.org",
		"https://ip.tyk.nu",
		"https://ipecho.net/plain",
		"https://l2.io/ip",
		"https://wtfismyip.com/text",
		"https://icanhazip.com",
		"https://ifconfig.me/ip",
		"https://ip.tyk.nu",
		"https://ipecho.net/plain",
	}

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		ip := string(body)
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	return "[EXTERNAL_IP]"
}
