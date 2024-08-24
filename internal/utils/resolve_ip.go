package utils

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

func ResolveIPWithPort(domainOrIP string) (string, error) {
	var host, port string

	if strings.HasPrefix(domainOrIP, "http://") || strings.HasPrefix(domainOrIP, "https://") {
		parsedURL, err := url.Parse(domainOrIP)
		if err != nil {
			return "", errors.New("invalid URL format")
		}
		host = parsedURL.Hostname()
		port = parsedURL.Port()
		if port == "" {
			port = "80"
		}
	} else {
		if strings.Contains(domainOrIP, ":") {
			parts := strings.Split(domainOrIP, ":")
			host = parts[0]
			port = parts[1]
		} else {
			host = domainOrIP
			port = "80"
		}
	}

	if net.ParseIP(host) == nil {
		ips, err := net.LookupIP(host)
		if err != nil || len(ips) == 0 {
			return "", errors.New("unable to resolve domain")
		}
		host = ips[0].String()
	}

	return net.JoinHostPort(host, port), nil
}

func ResolveIP(domainOrIP string) (string, error) {
	if net.ParseIP(domainOrIP) != nil {
		return domainOrIP, nil
	}

	parsedURL, err := url.Parse(domainOrIP)
	if err == nil && parsedURL.Host != "" {
		host := parsedURL.Hostname()
		domainOrIP = host
	}

	ips, err := net.LookupIP(domainOrIP)
	if err != nil || len(ips) == 0 {
		return "", errors.New("unable to resolve domain")
	}

	return ips[0].String(), nil
}
