package metrics

import (
	"net/http"

	"order-service/internal/pkg/metrics"

	"github.com/felixge/httpsnoop"
	"github.com/go-chi/chi/v5"
)

func Middleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			// Захватываем метрики
			m := httpsnoop.CaptureMetrics(h, writer, request)

			// Получаем путь ПОСЛЕ выполнения
			routeContext := chi.RouteContext(request.Context())
			path := routeContext.RoutePattern()

			if routeContext == nil || (routeContext.RoutePattern() == "/*" && m.Code == http.StatusNotFound) {
				path = "unknown"
			}

			metrics.IncHttpRequestsTotal(request.Method, path, m.Code)
			metrics.ObserveHttpRequestDuration(request.Method, path, m.Code, m.Duration)
		})
	}
}
