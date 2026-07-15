package app

import (
	"net"
	"strings"
	"sync"
	"time"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/labstack/echo/v4"
)

type UpstreamStatus struct {
	Address string `json:"address"`
	Up      bool   `json:"up"`
	Error   string `json:"error,omitempty"`
}

type LivenessStatus struct {
	Host      string           `json:"host"`
	Upstreams []UpstreamStatus `json:"upstreams"`
}

func checkLiveness(addr string) UpstreamStatus {
	dialAddr := addr
	scheme := ""
	if strings.HasPrefix(addr, "http://") {
		dialAddr = strings.TrimPrefix(addr, "http://")
		scheme = "http"
	} else if strings.HasPrefix(addr, "https://") {
		dialAddr = strings.TrimPrefix(addr, "https://")
		scheme = "https"
	} else if strings.HasPrefix(addr, "h2c://") {
		dialAddr = strings.TrimPrefix(addr, "h2c://")
	}

	if !strings.Contains(dialAddr, ":") {
		if scheme == "https" {
			dialAddr += ":443"
		} else {
			dialAddr += ":80"
		}
	}

	conn, err := net.DialTimeout("tcp", dialAddr, 2*time.Second)
	if err != nil {
		return UpstreamStatus{Address: addr, Up: false, Error: err.Error()}
	}
	conn.Close()
	return UpstreamStatus{Address: addr, Up: true}
}

func (a *App) Liveness(c echo.Context) ([]LivenessStatus, error) {
	content, err := a.LastCaddyfile(c)
	if err != nil {
		return nil, err
	}

	blocks, err := caddyfile.Parse("Caddyfile", []byte(content))
	if err != nil {
		return nil, err
	}

	var rawStatuses []LivenessStatus

	for _, b := range blocks {
		var host string
		if len(b.Keys) > 0 {
			host = b.Keys[0].Text
		}

		var upstreams []UpstreamStatus

		for _, s := range b.Segments {
			if len(s) > 0 && s[0].Text == "reverse_proxy" {
				for i := 1; i < len(s); i++ {
					token := s[i].Text
					if token == "{" {
						break
					}
					if token == "*" || strings.HasPrefix(token, "/") || strings.HasPrefix(token, "@") {
						continue
					}
					upstreams = append(upstreams, UpstreamStatus{Address: token})
				}
			}
		}

		if len(upstreams) > 0 {
			rawStatuses = append(rawStatuses, LivenessStatus{Host: host, Upstreams: upstreams})
		}
	}

	// Process concurrently across all hosts
	statuses := make([]LivenessStatus, len(rawStatuses))
	var wg sync.WaitGroup

	for i, status := range rawStatuses {
		wg.Add(1)
		go func(i int, status LivenessStatus) {
			defer wg.Done()
			var innerWg sync.WaitGroup
			var mu sync.Mutex
			checkedUpstreams := make([]UpstreamStatus, 0, len(status.Upstreams))

			for _, u := range status.Upstreams {
				innerWg.Add(1)
				go func(addr string) {
					defer innerWg.Done()
					res := checkLiveness(addr)
					mu.Lock()
					checkedUpstreams = append(checkedUpstreams, res)
					mu.Unlock()
				}(u.Address)
			}
			innerWg.Wait()
			statuses[i] = LivenessStatus{Host: status.Host, Upstreams: checkedUpstreams}
		}(i, status)
	}
	wg.Wait()

	return statuses, nil
}
