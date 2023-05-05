package stopwatch

import (
	"context"
	"net/http"
)

type contextKey string

var (
	contextKeyStopWatch = contextKey("stopwatch")
)

func StopWatchMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := NewStopWatch()
		sw.Start()
		r = requestWithStopWatch(r, sw)
		next.ServeHTTP(w, r)

		sw.ShowResults()
	})
}

func StopWatchFromContext(ctx context.Context) StopWatch {
	sw := ctx.Value(contextKeyStopWatch)
	if sw == nil {
		return StartNoopStopWatch()
	}
	return sw.(StopWatch)
}

func contextWithStopWatch(ctx context.Context, sw StopWatch) context.Context {
	return context.WithValue(ctx, contextKeyStopWatch, sw)
}

func requestWithStopWatch(r *http.Request, sw StopWatch) *http.Request {
	return r.WithContext(contextWithStopWatch(r.Context(), sw))
}
