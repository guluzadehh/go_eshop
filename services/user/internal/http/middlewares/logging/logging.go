package loggingmdw

import (
	"bufio"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	requestmdw "github.com/guluzadehh/go_eshop/services/user/internal/http/middlewares/request"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func newResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	if statusCode != http.StatusOK {
		w.statusCode = statusCode
		w.ResponseWriter.WriteHeader(statusCode)
	}
}

func (w *responseWriterWrapper) Write(data []byte) (int, error) {
	written, err := w.ResponseWriter.Write(data)
	w.bytesWritten = written
	return written, err
}

func (w *responseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hijacker.Hijack()
}

func LogRequests(log *slog.Logger) mux.MiddlewareFunc {
	log = log.With(slog.String("component", "middleware/logging"))

	return func(next http.Handler) http.Handler {
		log.Info("logging middleware is enabled")

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", requestmdw.GetRequestId(r)),
			)

			ww := newResponseWriterWrapper(w)

			t1 := time.Now()
			defer func() {
				log.Info("request completed",
					slog.Int("status", ww.statusCode),
					slog.Int("bytes", ww.bytesWritten),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
