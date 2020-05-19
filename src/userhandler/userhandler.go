package userhandler

import (
	"context"
	"errors"

	"github.com/beldin0/users/src/logging"
	pb "github.com/beldin0/users/src/user"
	"github.com/beldin0/users/src/userservice"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
)

type userHandler struct {
	service *userservice.Service
}

// New returns a userHandler instance
func New(db *sqlx.DB) pb.UserServiceServer {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50),
		first_name_lower VARCHAR(50),
		last_name VARCHAR(50),
		last_name_lower VARCHAR(50),
		nickname VARCHAR(30),
		nickname_lower VARCHAR(30) UNIQUE,
		password VARCHAR(32),
		email VARCHAR(50) UNIQUE,
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

func (h *userHandler) Add(ctx context.Context, user *pb.User) (*pb.User, error) {
	err := h.service.Add(user)
	if errors.Is(err, userservice.ErrDuplicate) {
		logging.NewLogger().Sugar().
			With("error", err).
			Info("duplicate user add request")
		return nil, err
	}
	if err != nil {
		logging.NewLogger().Sugar().
			With("error", err).
			Warn("database error")
		return nil, err
	}
	logging.NewLogger().Sugar().
		With("user", user).
		Info("user added")
	return user, nil
}

func (h *userHandler) Search(ctx context.Context, user *pb.User) (*pb.UsersResponse, error) {
	search, err := buildSearch(user)
	if err != nil {
		logging.NewLogger().Sugar().
			With("user", user).
			With("error", err).
			Warn("error building search")
		return nil, err
	}
	users, err := h.service.Get(search)
	if err != nil {
		logging.NewLogger().Sugar().
			With("user", user).
			With("error", err).
			Warn("error executing search")
		return nil, err
	}
	return &pb.UsersResponse{Users: users}, err
}

func (h *userHandler) Delete(ctx context.Context, id *pb.UserId) (*empty.Empty, error) {
	err := h.service.Delete(id.Id)
	if err != nil {
		logging.NewLogger().Sugar().
			With("error", err).
			Warn("database error")
	}
	return &empty.Empty{}, err
}

func (h *userHandler) Get(ctx context.Context, id *pb.UserId) (*pb.User, error) {
	user, err := h.service.Get(userservice.Get(id.Id))
	if err != nil {
		logging.NewLogger().Sugar().
			With("id", id.Id).
			With("error", err).
			Warn("server error")
		return nil, err
	}
	return user[0], err
}

func (h *userHandler) Modify(ctx context.Context, user *pb.User) (*pb.User, error) {
	err := h.service.Modify(user.Id, user)
	if err != nil {
		logging.NewLogger().Sugar().
			With("error", err).
			Warn("database error")
		return nil, err
	}
	return user, err
}

func buildSearch(user *pb.User) (*userservice.SearchOptions, error) {
	search := userservice.Search()
	var filters bool
	if user.Country != "" {
		search.Country(user.Country)
		filters = true
	}
	if user.Email != "" {
		search.Email(user.Email)
		filters = true
	}
	if user.Nickname != "" {
		search.Nickname(user.Nickname)
		filters = true
	}
	if user.FirstName != "" {
		search.FirstName(user.FirstName)
		filters = true
	}
	if user.LastName != "" {
		search.LastName(user.LastName)
		filters = true
	}
	if !filters {
		return nil, errors.New("no search parameters provided")
	}
	return search, nil
}
