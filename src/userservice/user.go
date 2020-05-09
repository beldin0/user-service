package userservice

import "strings"

type User struct {
	Firstname string `db:"first_name"`
	Lastname  string `db:"last_name"`
	Nickname  string `db:"nickname"`
	Password  string `db:"password"`
	Email     string `db:"email"`
	Country   string `db:"country"`
}

func (u User) insert() insertUser {
	return insertUser{
		Firstname:      u.Firstname,
		FirstnameLower: strings.ToLower(u.Firstname),
		Lastname:       u.Lastname,
		LastnameLower:  strings.ToLower(u.Lastname),
		Nickname:       u.Nickname,
		NicknameLower:  strings.ToLower(u.Nickname),
		Password:       u.Password,
		Email:          strings.ToLower(u.Email),
		Country:        strings.ToUpper(u.Country),
	}
}

type insertUser struct {
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
