package userservice

import "strings"

// User represents a user
type User struct {
	UserID    int    `db:"id" json:"userId"`
	Firstname string `db:"first_name" json:"firstName"`
	Lastname  string `db:"last_name" json:"lastName"`
	Nickname  string `db:"nickname" json:"nickname"`
	Password  string `db:"password" json:"password"`
	Email     string `db:"email" json:"email"`
	Country   string `db:"country" json:"country"`
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
