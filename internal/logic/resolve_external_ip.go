package logic

import (
	"io"
	"net"
	"net/http"
)

var serviceUrls = []string{
	"https://api.ipify.org",
	"https://api.seeip.org",
	"https://ifconfig.co/ip",
	"https://ifconfig.io/ip",
	"https://icanhazip.com",
	"https://ident.me",
	"https://myexternalip.com/raw",
	"https://ipinfo.io/ip",
	"https://icanhazip.com",
	"https://ip.seeip.org",
	"https://wtfismyip.com/text",
	"https://ipecho.net/plain",
	"https://icanhazip.com",
	"https://ifconfig.me/ip",
	"https://ip.tyk.nu",
	"https://ipecho.net/plain",
	"https://l2.io/ip",
}

func ResolveExternalIP() string {
	for i := 0; i < len(serviceUrls)-1; i += 2 {
		ip1 := fetchIP(serviceUrls[i])
		ip2 := fetchIP(serviceUrls[i+1])

		if ip1 != "" && ip1 == ip2 {
			return ip1
		}
	}

	return "[EXTERNAL_IP]"
}

func fetchIP(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	ip := string(body)
	if net.ParseIP(ip) != nil {
		return ip
	}

	return ""
}
