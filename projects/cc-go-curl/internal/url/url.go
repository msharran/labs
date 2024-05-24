package url

import (
	"errors"
	"strings"
)

var (
	schemePorts = map[string]string{
		"http":  "80",
		"https": "443",
	}
)

type URL struct {
	Scheme string
	Host   string
	Port   string
	Path   string
}

func (u *URL) PortOrDefault() string {
	if u.Port == "" {
		return schemePorts[u.Scheme]
	}

	return u.Port
}

func Parse(url string) (*URL, error) {
	scheme, rest, ok := strings.Cut(url, "://")
	if !ok {
		return nil, errors.New("invalid URL format")
	}

	if _, ok := schemePorts[scheme]; !ok {
		return nil, errors.New("unsupported protocol scheme")
	}

	u := new(URL)
	u.Scheme = scheme

	addr, rest, _ := strings.Cut(rest, "/")
	u.Path = "/" + rest

	host, port, ok := strings.Cut(addr, ":")
	if !ok || port == "" {
		port = schemePorts[scheme]
	}
	u.Host = host
	u.Port = port

	return u, nil
}
