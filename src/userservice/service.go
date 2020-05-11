package userservice

import (
	"strings"

	"github.com/beldin0/users/src/logging"
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
func (s *Service) Add(u User) error {
	_, err := s.db.NamedExec(sqlInsert, u.insert())
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
func (s *Service) Get(o *SearchOptions) ([]User, error) {
	results := []User{}
	err := s.db.Select(&results, sqlGet+o.where())
	if err != nil {
		logging.NewLogger().Sugar().With("error", err).Warn("error executing query")
	}
	return results, err
}

// Modify updates a users details based on the provided SearchOptions
// The searchoptions must include either an email address or a nickname and country
// Search terms must match exactly the entries in the existing user row.
func (s *Service) Modify(o *SearchOptions, u User) error {
	where, err := o.modify()
	if err != nil {
		logging.NewLogger().Sugar().With("error", err).Warn("error executing query")
		return err
	}
	_, err = s.db.NamedExec(sqlModify+where, u.insert())
	if err != nil {
		logging.NewLogger().Sugar().With("error", err).Warn("error executing query")
	}
	return err
}

// Delete deletes a users details based on the provided SearchOptions
// The searchoptions must include either an email address or a nickname and country
// Search terms must match exactly the entries in the existing user row.
func (s *Service) Delete(o *SearchOptions) error {
	where, err := o.modify()
	if err != nil {
		logging.NewLogger().Sugar().With("error", err).Warn("error executing query")
		return err
	}
	_, err = s.db.Exec(sqlDelete+where, nil)
	if err != nil {
		logging.NewLogger().Sugar().With("error", err).Warn("error executing query")
	}
	return err
}
