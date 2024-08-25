package utils

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

func parseAndLookupIP(input string) (string, string, error) {
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		input = "http://" + input
	}

	parsedURL, err := url.Parse(input)
	if err != nil {
		return "", "", err
	}

	ips, err := net.LookupIP(parsedURL.Hostname())
	if err != nil {
		return "", "", err
	}

	if len(ips) > 0 {
		return ips[0].String(), parsedURL.Port(), nil
	}

	return "", "", errors.New("unable to resolve IP address")
}

func ResolveIPWithOutPort(input string) (string, error) {
	if net.ParseIP(input) != nil {
		return input, nil
	}

	host, _, err := net.SplitHostPort(input)
	if err == nil && net.ParseIP(host) != nil {
		return host, nil
	}

	ip, _, err := parseAndLookupIP(input)
	if err != nil {
		return "", err
	}

	return ip, nil
}

func ResolveIPWithPort(input string) (string, error) {
	if net.ParseIP(input) != nil {
		return net.JoinHostPort(input, "80"), nil
	}

	host, port, err := net.SplitHostPort(input)
	if err == nil && net.ParseIP(host) != nil {
		if port == "" {
			port = "80"
		}
		return net.JoinHostPort(host, port), nil
	}

	ip, port, err := parseAndLookupIP(input)
	if err != nil {
		return "", err
	}

	if port == "" {
		port = "80"
	}

	return net.JoinHostPort(ip, port), nil
}

func ResolveHostname(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		return ip
	}

	return names[0]
}
