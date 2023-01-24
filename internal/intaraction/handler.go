package intaraction

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ChizarR/stats-service/internal/apperror"
	"github.com/ChizarR/stats-service/internal/handlers"
	"github.com/ChizarR/stats-service/pkg/logging"
	"github.com/ChizarR/stats-service/pkg/rest"
)

var _ handlers.Handler = &handler{}

const (
	totlaIntaractionsURL          = "/intaraction/total"
	todayInraractionsURL          = "/intaraction/today"
	updateCategoryIntaractionsURL = "/intaraction/update"
	getAllCategoryIntaractionsURL = "/intaraction/all"
)

type handler struct {
	service IntaractionService
	logger  *logging.Logger
}

func NewHandler(service IntaractionService, logger *logging.Logger) handlers.Handler {
	return &handler{service: service, logger: logger}
}

func (h *handler) Register(router *http.ServeMux) {
	router.HandleFunc(totlaIntaractionsURL, apperror.Middleware(h.GetTotalIntaractions))
	router.HandleFunc(todayInraractionsURL, apperror.Middleware(h.GetTodayIntaractions))
	router.HandleFunc(updateCategoryIntaractionsURL, apperror.Middleware(h.UpdateIntaractions))
	router.HandleFunc(getAllCategoryIntaractionsURL, apperror.Middleware(h.GetAllCategoryIntaractons))
}

func (h *handler) GetTotalIntaractions(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		allIntrs, err := h.service.GetAllIntaractions(context.Background())
		if err != nil {
			return err
		}

		resp := rest.NewResponse(true, "All intaractions", allIntrs)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp.Marshal())
	default:
		w.Write([]byte("Only GET available"))
	}
	return nil
}

func (h *handler) GetTodayIntaractions(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		intr, err := h.service.GetTodayIntaractions(context.Background())
		if err != nil {
			return err
		}

		resp := rest.NewResponse(true, "Today intaraction", intr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp.Marshal())
	default:
		w.Write([]byte("Only GET available"))
	}
	return nil
}

func (h *handler) UpdateIntaractions(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodPost:
		bytes, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			return err
		}

		var intaractionsDTO IntaractionDTO
		if err = json.Unmarshal(bytes, &intaractionsDTO); err != nil {
			return err
		}

		err = h.service.AddNewIntaractions(context.Background(), intaractionsDTO)
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusCreated)
		return nil
	default:
		w.Write([]byte("Only POST available"))
	}
	return nil
}

func (h *handler) GetAllCategoryIntaractons(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		allIntaractions, err := h.service.GetAllIntaractions(context.Background())
		if err != nil {
			return err
		}

		resp := rest.NewResponse(true, "All intaractions", allIntaractions)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp.Marshal())
		return nil
	default:
		w.Write([]byte("Only GET available"))
		return nil
	}
}
