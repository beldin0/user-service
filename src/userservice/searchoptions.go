package userservice

import (
	"errors"
	"fmt"
	"strings"
)

// SearchOptions provides the means of searching for one or many users
type SearchOptions struct {
	options map[string]string
}

// Search begins a new search
func Search() *SearchOptions {
	return &SearchOptions{
		options: make(map[string]string),
	}
}

// Country adds the specified country to the search parameters
func (o *SearchOptions) Country(country string) *SearchOptions {
	o.options["country"] = country
	return o
}

// Email adds the specified email address to the search parameters
func (o *SearchOptions) Email(email string) *SearchOptions {
	o.options["email"] = strings.ToLower(email)
	return o
}

// Nickname adds the specified nickname to the search parameters
func (o *SearchOptions) Nickname(nickname string) *SearchOptions {
	o.options["nickname_lower"] = strings.ToLower(nickname)
	return o
}

// Name adds the specified first name and last name to the search parameters
func (o *SearchOptions) Name(first, last string) *SearchOptions {
	o.options["first_name_lower"] = strings.ToLower(first)
	o.options["last_name_lower"] = strings.ToLower(last)
	return o
}

func (o *SearchOptions) where() string {
	if o == nil || len(o.options) == 0 {
		return ""
	}
	options := []string{}
	for field, value := range o.options {
		options = append(options, `"`+field+`" LIKE '%`+value+`%'`)
	}
	return " WHERE " + strings.Join(options, " AND ")
}

func (o *SearchOptions) whereExact() string {
	if o == nil || len(o.options) == 0 {
		return ""
	}
	options := []string{}
	for field, value := range o.options {
		options = append(options, `"`+field+`"='`+value+`'`)
	}
	return " WHERE " + strings.Join(options, " AND ")
}

func (o *SearchOptions) modify() (string, error) {
	if o == nil || len(o.options) == 0 {
		return "", errors.New("required searchoptions not provided for modify")
	}
	_, email := o.options["email"]
	_, nick := o.options["nickname_lower"]
	_, country := o.options["country"]
	switch {
	case email:
	case nick && country:
	default:
		return "", fmt.Errorf("required searchoptions not provided for modify: provided with %v", o.options)
	}
	return o.whereExact(), nil
}
