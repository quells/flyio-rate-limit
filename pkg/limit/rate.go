package limit

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Counter interface {
	GetCount(ctx context.Context, ip string) (int, error)
	Increment(ctx context.Context, ip string, ttl int) error
}

type IPExtractor func(*http.Request) string

func StatusCodeRate(counter Counter, getIP IPExtractor, limit, ttl int) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return statusCodeCounter{
			counter: counter,
			getIP:   getIP,
			next:    next,
			limit:   limit,
			ttl:     ttl,
		}
	}
}

type statusCodeCounter struct {
	counter Counter
	getIP   IPExtractor
	next    http.Handler
	limit   int
	ttl     int
}

func (rl statusCodeCounter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ip := rl.getIP(req)
	logger := log.WithField("ip", ip)

	count, err := rl.counter.GetCount(req.Context(), ip)
	if err != nil {
		logger.WithError(err).Error("could not get count")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if count >= rl.limit {
		rw.Header().Set("Retry-After", strconv.Itoa(rl.ttl))
		rw.WriteHeader(http.StatusTooManyRequests)
		return
	}

	recorder := httptest.NewRecorder()
	rl.next.ServeHTTP(recorder, req)
	if recorder.Code == http.StatusUnauthorized || recorder.Code == http.StatusForbidden {
		err = rl.counter.Increment(req.Context(), ip, rl.ttl)
		if err != nil {
			logger.WithError(err).Error("could not get count")
		}
	}

	copyResponse(logger, rw, recorder.Result())
}

func copyResponse(logger *log.Entry, dst http.ResponseWriter, src *http.Response) {
	dst.WriteHeader(src.StatusCode)
	for k, vs := range src.Header {
		for _, v := range vs {
			dst.Header().Add(k, v)
		}
	}
	if src.Body != nil {
		_, err := io.Copy(dst, src.Body)
		if err != nil {
			logger.WithError(err).Error("could not copy body")
		}
	}
}
