package userservice

import (
	"strings"

	"github.com/beldin0/users/src/user"
)

func toInsert(u *user.User) insertUser {
	return insertUser{
		UserID:         &u.Id,
		Firstname:      u.FirstName,
		FirstnameLower: strings.ToLower(u.FirstName),
		Lastname:       u.LastName,
		LastnameLower:  strings.ToLower(u.LastName),
		Nickname:       u.Nickname,
		NicknameLower:  strings.ToLower(u.Nickname),
		Password:       u.Password,
		Email:          strings.ToLower(u.Email),
		Country:        strings.ToUpper(u.Country),
	}
}

type insertUser struct {
	UserID         *int32 `db:"id"`
	Firstname      string `db:"first_name"`
	FirstnameLower string `db:"first_name_lower"`
	Lastname       string `db:"last_name"`
	LastnameLower  string `db:"last_name_lower"`
	Nickname       string `db:"nickname"`
	NicknameLower  string `db:"nickname_lower"`
	Password       string `db:"password"`
	Email          string `db:"email"`
	Country        string `db:"country"`
}
