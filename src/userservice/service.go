package userservice

import (
	"strings"

	"github.com/beldin0/users/src/logging"
	"github.com/beldin0/users/src/user"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// New returns a Service instance utilising the provided database
func New(db *sqlx.DB) *Service {
	return &Service{
		db: db,
	}
}

// Service is a User Service, providing the methods to interact with the database
type Service struct {
	db *sqlx.DB
}

// Add adds a new User to the database
func (s *Service) Add(u *user.User) error {
	_, err := s.db.NamedExec(sqlInsert, toInsert(u))
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			err = errors.Wrap(ErrDuplicate, err.Error())
		}
		logging.NewLogger().Sugar().With("error", err).Warn("error executing query")
	}
	return err
}

// Get searches for users based on the provided SearchOptions
// a nil SearchOptions returns a list of all users
// matches are made using LIKE so can be partial search terms
func (s *Service) Get(o *SearchOptions) ([]user.User, error) {
	results := []user.User{}
	err := s.db.Select(&results, sqlGet+o.where())
	if err != nil {
		logging.NewLogger().Sugar().With("error", err).Warn("error executing query")
	}
	return results, err
}

// Modify updates a users details based on the provided SearchOptions
// The searchoptions must include either an email address or a nickname and country
// Search terms must match exactly the entries in the existing user row.
func (s *Service) Modify(userID int32, u *user.User) error {
	u.Id = userID
	_, err := s.db.NamedExec(sqlModify, toInsert(u))
	if err != nil {
		logging.NewLogger().Sugar().With("error", err).Warn("error executing query")
	}
	return err
}

// Delete deletes a users details based on the provided SearchOptions
// The searchoptions must include either an email address or a nickname and country
// Search terms must match exactly the entries in the existing user row.
func (s *Service) Delete(userID int) error {
	_, err := s.db.Exec(sqlDelete, userID)
	if err != nil {
		logging.NewLogger().Sugar().With("error", err).Warn("error executing query")
	}
	return err
}
