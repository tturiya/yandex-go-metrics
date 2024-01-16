package logger

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type rwWrapper struct {
	http.ResponseWriter
	status       int
	bytesWritten int
	wroteHeader  bool
	response     string
}

func (rw *rwWrapper) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.bytesWritten = 0
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
	return
}

func (rw *rwWrapper) Write(b []byte) (int, error) {
	rw.wroteHeader = true
	written, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten = written
	rw.response = string(b)
	return written, err
}

func Logger(next http.Handler) http.Handler {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	fn := func(w http.ResponseWriter, r *http.Request) {
		reqTiming := time.Now()
		wrapped := rwWrapper{w, 200, 0, false, ""}
		next.ServeHTTP(&wrapped, r)

		log.Info().
			Int("status", wrapped.status).
			Str("method", r.Method).
			Str("URI", r.URL.String()).
			Str("latency_human", time.Since(reqTiming).String()).
			Int("response_bytes", wrapped.bytesWritten).
			Str("response", wrapped.response).
			Msg("")
	}

	return http.HandlerFunc(fn)
}
