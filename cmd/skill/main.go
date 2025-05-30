// формат запроса https://yandex.ru/dev/dialogs/alice/doc/ru/request
package main

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/VladimirAzanza/alisa_skill/internal/logger"
	"github.com/VladimirAzanza/alisa_skill/internal/store/pg"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func gzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()

			ow.Header().Set("Content-Encoding", "gzip")
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	}
}

// go build -o skill
// ./skill -a :8081
// RUN_ADDR=:8081 ./skill
// RUN_ADDR=:8082 ./skill -a :8081
func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}

	logger.Log.Info("Running server", zap.String("address", flagRunAddr))

	conn, err := sql.Open("pgx", flagDatabaseURI)
	if err != nil {
		return err
	}

	appInstance := newApp(pg.NewStore(conn))
	return http.ListenAndServe(flagRunAddr, logger.RequestLogger(gzipMiddleware(appInstance.webhook)))
}
