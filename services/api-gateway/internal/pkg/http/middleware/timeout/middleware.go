package timeout

import (
	"context"
	"net/http"
	"time"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		timeoutStr := request.Header.Get("X-Request-Timeout")

		if d, err := time.ParseDuration(timeoutStr); err == nil && d > 0 {
			ctx, cancel := context.WithTimeout(request.Context(), d)
			defer cancel()

			next.ServeHTTP(writer, request.WithContext(ctx))
			return
		}

		next.ServeHTTP(writer, request)
	})
}
