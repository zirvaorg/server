package utils

import "strings"

func FormatURL(ipOrDomain string) string {
	if !strings.HasPrefix(ipOrDomain, "http://") && !strings.HasPrefix(ipOrDomain, "https://") {
		ipOrDomain = "http://" + ipOrDomain
	}
	return ipOrDomain
}
