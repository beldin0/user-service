package userservice

import (
	"github.com/jmoiron/sqlx"
)

// New returns a Service instance utilising the provided database
func New(db *sqlx.DB) *Service {
	return &Service{
		db: db,
	}
}

type Service struct {
	db *sqlx.DB
}

func (s *Service) Add(u User) error {
	_, err := s.db.NamedExec(sqlInsert, u.insert())
	return err
}

func (s *Service) Get(o *SearchOptions) ([]User, error) {
	results := []User{}
	err := s.db.Select(&results, sqlGet+o.where())
	return results, err
}

func (s *Service) Modify(o *SearchOptions, u User) error {
	where, err := o.modify()
	if err != nil {
		return err
	}
	_, err = s.db.NamedExec(sqlModify+where, u.insert())
	return err
}
