package userhandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/beldin0/users/src/logging"
	"github.com/beldin0/users/src/userservice"
	"github.com/jmoiron/sqlx"
)

type userHandler struct {
	service *userservice.Service
}

// New returns a userHandler instance
func New(db *sqlx.DB) http.Handler {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		first_name VARCHAR(50),
		first_name_lower VARCHAR(50),
		last_name VARCHAR(50),
		last_name_lower VARCHAR(50),
		nickname VARCHAR(30),
		nickname_lower VARCHAR(30) UNIQUE,
		password VARCHAR(32),
		email VARCHAR(50) UNIQUE PRIMARY KEY,
		country VARCHAR(3)
	)`)
	if err != nil {
		logging.NewLogger().Sugar().
			With("error", err).
			Error("problem setting up database")
		panic(err)
	}
	return &userHandler{
		service: userservice.New(db),
	}
}

func (h *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logging.NewLogger().Sugar().
		With("path", r.URL.Path).
		With("params", r.URL.RawQuery).
		Info(r.Method)
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	case http.MethodPost:
		h.Post(w, r)
	case http.MethodPut:
		h.Put(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *userHandler) Get(w http.ResponseWriter, r *http.Request) {
	search, err := buildSearch(r.URL.Query())
	if err != nil {
		logging.NewLogger().Sugar().
			With("error", err).
			Warn("bad request")
		http.Error(w, "Request error", http.StatusBadRequest)
		return
	}
	users, err := h.service.Get(search)
	if err != nil {
		logging.NewLogger().Sugar().
			With("error", err).
			Warn("server error")
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (h *userHandler) Post(w http.ResponseWriter, r *http.Request) {
	var user userservice.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logging.NewLogger().Sugar().
			With("error", err).
			Warn("bad request")
		http.Error(w, "unable to decode request body", http.StatusBadRequest)
		return
	}
	err = h.service.Add(user)
	if errors.As(err, &userservice.ErrDuplicate) {
		logging.NewLogger().Sugar().
			With("error", err).
			Info("duplicate user add request")
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}
	if err != nil {
		logging.NewLogger().Sugar().
			With("error", err).
			Warn("database error")
		http.Error(w, "error inserting into database", http.StatusInternalServerError)
		return
	}
	logging.NewLogger().Sugar().
		With("user", user).
		Info("user added")
	w.WriteHeader(http.StatusCreated)
}

func (h *userHandler) Put(w http.ResponseWriter, r *http.Request) {
	search, err := buildSearch(r.URL.Query())
	if err != nil {
		logging.NewLogger().Sugar().
			With("error", err).
			Warn("bad request")
		http.Error(w, "Request error", http.StatusBadRequest)
		return
	}
	var user userservice.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logging.NewLogger().Sugar().
			With("error", err).
			Warn("bad request")
		http.Error(w, "unable to decode request body", http.StatusBadRequest)
		return
	}
	err = h.service.Modify(search, user)
	if err != nil {
		logging.NewLogger().Sugar().
			With("error", err).
			Warn("database error")
		http.Error(w, "error updating user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
	search, err := buildSearch(r.URL.Query())
	if err != nil {
		logging.NewLogger().Sugar().
			With("error", err).
			Warn("bad request")
		http.Error(w, "Request error", http.StatusBadRequest)
		return
	}
	err = h.service.Delete(search)
	if err != nil {
		logging.NewLogger().Sugar().
			With("error", err).
			Warn("database error")
		http.Error(w, "error deleting user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func buildSearch(params url.Values) (*userservice.SearchOptions, error) {
	search := userservice.Search()
	for key := range params {
		switch key {
		case "country":
			search.Country(params.Get(key))
			fallthrough
		case "email":
			search.Email(params.Get(key))
			fallthrough
		case "nickname":
			search.Nickname(params.Get(key))
			fallthrough
		case "name":
			value, err := url.QueryUnescape(params.Get(key))
			if err != nil {
				return search, err
			}
			names := strings.Split(value, " ")
			if len(names) < 2 {
				return search, err
			}
			search.Name(names[0], names[len(names)-1])
		}
	}
	return search, nil
}
