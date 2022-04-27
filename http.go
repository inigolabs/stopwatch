package stopwatch

import (
	"context"
	"net/http"
)

type contextKey string

var (
	contextKeyTimer = contextKey("timer")
)

func TimerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := NewTimer()
		timer.Start()
		r = requestWithTimer(r, timer)
		next.ServeHTTP(w, r)
		timer.Stop()

		timer.ShowResults()
	})
}

func TimerFromContext(ctx context.Context) Timer {
	timer := ctx.Value(contextKeyTimer)
	if timer == nil {
		return StartNoopTimer()
	}
	return timer.(Timer)
}

func contextWithTimer(ctx context.Context, timer Timer) context.Context {
	return context.WithValue(ctx, contextKeyTimer, timer)
}

func requestWithTimer(r *http.Request, timer Timer) *http.Request {
	return r.WithContext(contextWithTimer(r.Context(), timer))
}
