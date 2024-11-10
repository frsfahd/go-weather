package server

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	bucket_size, _ = strconv.Atoi(os.Getenv("BUCKET_SIZE"))
	token_rate, _  = strconv.Atoi(os.Getenv("TOKEN_RATE"))
)

func Limiter() gin.HandlerFunc {
	type client struct {
		limiter   rate.Limiter
		last_seen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			// Lock the mutex to protect this section from race conditions.
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.last_seen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	// limit := rate.NewLimiter(token_rate, bucket_size)
	return func(ctx *gin.Context) {
		// net.SplitHostPort(ctx.RemoteIP())
	}
}
