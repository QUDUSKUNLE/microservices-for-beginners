package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type client struct {
	count int
	time  time.Time
}

var (
	store = make(map[string]*client)
	mu    = sync.Mutex{}
)

func init() {
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			mu.Lock()
			for ip, c := range store {
				if time.Since(c.time) > time.Minute {
					delete(store, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		mu.Lock()
		c, ok := store[ip]
		if !ok {
			store[ip] = &client{count: 1, time: time.Now()}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if time.Since(c.time) > time.Minute {
			c.count = 0
			c.time = time.Now()
		}

		c.count++
		if c.count > 10 {
			mu.Unlock()
			http.Error(w, "too many requests", 429)
			return
		}

		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
