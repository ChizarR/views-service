package apperror

import (
	"errors"
	"net/http"

	"github.com/ChizarR/stats-service/pkg/logging"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	logger := logging.GetLogger()
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppError
		err := h(w, r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			if errors.As(err, &appErr) {
				err = err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(UndefinedError.Marshal())
			}

			logger.Error(err)
			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err).Marshal())
		}
	}
}
