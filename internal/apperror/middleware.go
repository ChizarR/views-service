package apperror

import (
	"errors"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
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

			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err).Marshal())
		}
	}
}