package mux

import (
	"net"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/luraproject/lura/v2/proxy"
	"github.com/luraproject/lura/v2/router/mux"
	"github.com/luraproject/lura/v2/transport/http/server"
)

const requestParamsAsterisk string = "*"

// NewRequestBuilder gets a RequestBuilder with the received ParamExtractor as a query param
// extraction mechanism
func NewRequestBuilder(paramExtractor mux.ParamExtractor) mux.RequestBuilder {
	return func(r *http.Request, queryString, headersToSend []string) *proxy.Request {
		params := paramExtractor(r)
		headers := make(map[string][]string, 3+len(headersToSend))

		for _, k := range headersToSend {
			if k == requestParamsAsterisk {
				headers = r.Header

				break
			}

			if h, ok := r.Header[textproto.CanonicalMIMEHeaderKey(k)]; ok {
				headers[k] = h
			}
		}

		headers["X-Forwarded-For"] = []string{clientIP(r)}
		headers["X-Forwarded-Host"] = []string{r.Host}
		// if User-Agent is not forwarded using headersToSend, we set
		// the KrakenD router User Agent value
		if _, ok := headers["User-Agent"]; !ok {
			headers["User-Agent"] = server.UserAgentHeaderValue
		} else {
			headers["X-Forwarded-Via"] = server.UserAgentHeaderValue
		}

		query := make(map[string][]string, len(queryString))
		queryValues := r.URL.Query()
		for i := range queryString {
			if queryString[i] == requestParamsAsterisk {
				query = queryValues

				break
			}

			if v, ok := queryValues[queryString[i]]; ok && len(v) > 0 {
				query[queryString[i]] = v
			}
		}

		return &proxy.Request{
			// keep original request path for wildcard routing
			Path:    r.URL.Path,
			Method:  r.Method,
			Query:   query,
			Body:    r.Body,
			Params:  params,
			Headers: headers,
		}
	}
}

// clientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
// Use X-Forwarded-For before X-Real-Ip as nginx uses X-Real-Ip with the proxy's IP.
func clientIP(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}

	if addr := r.Header.Get("X-Appengine-Remote-Addr"); addr != "" {
		return addr
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
