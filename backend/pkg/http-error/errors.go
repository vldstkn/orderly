package http_error

import (
	"log/slog"
	"net/http"
)

func BadRequest(w http.ResponseWriter, log *slog.Logger, opt string, err error) {
	log.Error(opt, slog.String("err", err.Error()))
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
