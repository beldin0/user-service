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
	rows, err := s.db.NamedQuery(sqlInsert, toInsert(u))
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			err = errors.Wrap(ErrDuplicate, err.Error())
		}
		logging.NewLogger().Sugar().With("error", err).Warn("error executing query")
		return err
	}
	defer rows.Close()
	if rows.Next() {
		var i int32
		if err := rows.Scan(&i); err != nil {
			return err
		}
		u.Id = i
	}
	logging.NewLogger().Sugar().
		With("function", "add").
		With("user", u).
		Info("new user added")
	return nil
}

// Get searches for users based on the provided SearchOptions
// a nil SearchOptions returns a list of all users
// matches are made using LIKE so can be partial search terms
func (s *Service) Get(o *SearchOptions) ([]*user.User, error) {
	query := sqlGet + o.where()
	rows, err := s.db.Query(query)
	if err != nil {
		logging.NewLogger().Sugar().
			With("query", query).
			With("error", err).
			Warn("error executing query")
		return nil, err
	}
	defer rows.Close()
	results := []*user.User{}
	for rows.Next() {
		u := user.User{}
		if err := rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Nickname, &u.Password, &u.Email, &u.Country); err != nil {
			logging.NewLogger().Sugar().
				With("query", sqlGet+o.where()).
				With("error", err).
				Warn("error processing rows query")
		}
		results = append(results, &u)
	}
	logging.NewLogger().Sugar().
		With("function", "get").
		With("search", o.options).
		With("results", len(results)).
		Info("returning results")
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
		return err
	}
	logging.NewLogger().Sugar().
		With("function", "modify").
		With("user", u).
		Info("user updated")
	return nil
}

// Delete deletes a users details based on the provided SearchOptions
// The searchoptions must include either an email address or a nickname and country
// Search terms must match exactly the entries in the existing user row.
func (s *Service) Delete(userID int32) error {
	_, err := s.db.Exec(sqlDelete, userID)
	if err != nil {
		logging.NewLogger().Sugar().With("error", err).Warn("error executing query")
		return err
	}
	logging.NewLogger().Sugar().
		With("function", "delete").
		With("userID", userID).
		Info("user deleted")
	return nil
}
