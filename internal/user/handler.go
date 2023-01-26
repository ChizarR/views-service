package user

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
	getUserViewsForTodayURL = "/users_stats/views/today"
	updateUserViewsURL      = "/users_stats/update_views"
	getAllUsersURL          = "/users_stats/all"
)

type handler struct {
	service UserStatService
	logger  *logging.Logger
}

func NewHandler(service UserStatService, logger *logging.Logger) handlers.Handler {
	return &handler{service: service, logger: logger}
}

func (h *handler) Register(router *http.ServeMux) {
	router.HandleFunc(getUserViewsForTodayURL, apperror.Middleware(h.GetUserViewsForToday))
	router.HandleFunc(updateUserViewsURL, apperror.Middleware(h.UpdateUserViews))
	router.HandleFunc(getAllUsersURL, apperror.Middleware(h.GetAllUsers))
}

func (h *handler) GetUserViewsForToday(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		bytes, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			return err
		}

		var getStatUserDTO GetStatUserDTO
		if err = json.Unmarshal(bytes, &getStatUserDTO); err != nil {
			return err
		}

		userStats, err := h.service.GetUserViewsForToday(context.Background(), getStatUserDTO.TgId)
		if err != nil {
			return err
		}

		resp := rest.NewResponse(true, "Get user stats", userStats)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp.Marshal())
	default:
		w.Write([]byte("Only GET available"))
	}
	return nil
}

func (h *handler) UpdateUserViews(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodPost:
		bytes, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			return err
		}

		var uDTO UserDTO
		if err = json.Unmarshal(bytes, &uDTO); err != nil {
			return err
		}

		err = h.service.UpdateUserViews(context.Background(), uDTO)
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

func (h *handler) GetAllUsers(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		allUsers, err := h.service.GetAllUsersStats(context.Background())
		if err != nil {
			return err
		}

		resp := rest.NewResponse(true, "All users", allUsers)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp.Marshal())
		return nil
	default:
		w.Write([]byte("Only GET available"))
		return nil
	}
}
