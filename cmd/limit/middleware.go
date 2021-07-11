package main

import (
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/quells/flyio-rate-limit/pkg/limit"
)

func rateLimiter() mux.MiddlewareFunc {
	var counter limit.Counter
	redisURL := os.Getenv("FLY_REDIS_CACHE_URL")
	if redisURL == "" {
		counter = limit.NewMapCounter()
	} else {
		counter = limit.NewRedisCounter(redisURL)
	}

	count := envInt("RATE_COUNT", 3)
	ttl := envInt("RATE_TTL", 10)

	return limit.StatusCodeRate(counter, getIP, count, ttl)
}

var hostPort = regexp.MustCompile(`^(.+):\d+$`)

func getIP(req *http.Request) string {
	flyIP := req.Header.Get("Fly-Client-IP")
	if flyIP != "" {
		return flyIP
	}

	forwardedFor := req.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		parts := strings.Split(forwardedFor, ", ")
		return parts[0]
	}

	matches := hostPort.FindStringSubmatch(req.RemoteAddr)
	switch len(matches) {
	case 2:
		return matches[1]
	default:
		return req.RemoteAddr
	}
}
