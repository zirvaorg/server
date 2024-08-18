package service

import (
	"errors"
	"io"
	"net/http"
	"server/internal/utils"
	"time"
)

const (
	UserAgent    = "zirva.org/1.0"
	MaxRedirects = 10
	TimeOut      = 10 * time.Second
)

type HttpResult struct {
	URI          string   `json:"uri"`
	StatusCode   int      `json:"status_code"`
	ResponseTime float64  `json:"response_time"`
	Error        string   `json:"error,omitempty"`
	Redirects    []string `json:"redirects,omitempty"`
	StatusCodes  []int    `json:"status_codes,omitempty"`
	ResolvedIPs  []string `json:"resolved_ips,omitempty"`
}

type redirectResult struct {
	Redirects   []string
	StatusCodes []int
	ResolvedIPs []string
	Error       error
}

func Http(ipOrDomain string) (HttpResult, error) {
	ipOrDomain = utils.FormatURL(ipOrDomain)

	resolvedIP, err := utils.ResolveIP(ipOrDomain)
	if err != nil {
		return HttpResult{URI: ipOrDomain, Error: err.Error(), ResolvedIPs: []string{resolvedIP}}, err
	}

	start := time.Now()
	client := &http.Client{
		Timeout: TimeOut,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", ipOrDomain, nil)
	if err != nil {
		return HttpResult{URI: ipOrDomain, Error: err.Error(), ResolvedIPs: []string{resolvedIP}}, err
	}
	req.Header.Set("User-Agent", UserAgent)

	result, err := handleRedirects(client, req)
	if err != nil {
		return HttpResult{URI: ipOrDomain, Error: err.Error(), ResolvedIPs: result.ResolvedIPs}, err
	}

	responseTime := time.Since(start).Seconds()
	return HttpResult{
		URI:          ipOrDomain,
		StatusCode:   result.StatusCodes[len(result.StatusCodes)-1],
		ResponseTime: responseTime,
		Redirects:    result.Redirects,
		StatusCodes:  result.StatusCodes,
		ResolvedIPs:  result.ResolvedIPs,
	}, nil
}

func handleRedirects(client *http.Client, req *http.Request) (redirectResult, error) {
	var result redirectResult

	for redirectCount := 0; redirectCount < MaxRedirects; redirectCount++ {
		resp, err := client.Do(req)
		if err != nil {
			result.Error = err
			return result, err
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)

		result.Redirects = append(result.Redirects, resp.Request.URL.String())
		result.StatusCodes = append(result.StatusCodes, resp.StatusCode)

		resolvedIP, err := utils.ResolveIP(resp.Request.URL.String())
		if err != nil {
			result.Error = err
			return result, err
		}
		result.ResolvedIPs = append(result.ResolvedIPs, resolvedIP)

		location, err := resp.Location()
		if err != nil {
			if errors.Is(err, http.ErrNoLocation) {
				break
			}
			result.Error = err
			return result, err
		}

		req, err = http.NewRequest("GET", location.String(), nil)
		if err != nil {
			result.Error = err
			return result, err
		}
		req.Header.Set("User-Agent", UserAgent)
	}

	return result, nil
}
