package utils

import (
	"errors"
	"net"
	"net/url"
)

func ResolveIP(domainOrIP string) (string, error) {
	parsedURL, err := url.Parse(domainOrIP)
	if err == nil && parsedURL.Host != "" {
		host := parsedURL.Hostname()
		domainOrIP = host
	}

	if net.ParseIP(domainOrIP) != nil {
		return domainOrIP, nil
	}

	ips, err := net.LookupIP(domainOrIP)
	if err != nil || len(ips) == 0 {
		return "", errors.New("unable to resolve domain")
	}

	return ips[0].String(), nil
}
