package userservice_test

import (
	"crypto/sha256"
	"strings"
	"testing"

	"github.com/beldin0/users/src/user"
	"github.com/beldin0/users/src/userservice"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestService_Add(t *testing.T) {
	tests := []struct {
		name    string
		u       *user.User
		wantErr bool
	}{
		{
			name:    "user test",
			u:       testUser(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, err := sqlx.Connect("sqlite3", ":memory:")
			require.NoError(t, err)
			require.NoError(t, setup(db))
			s := userservice.New(db)
			err = s.Add(tt.u)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
			user := user.User{}
			require.NoError(t, db.Get(&user, `SELECT first_name, last_name, nickname, password, email, country FROM users`))
			require.Equal(t, tt.u, user)
		})
	}
}

func TestService_Get(t *testing.T) {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	require.NoError(t, err)
	require.NoError(t, setup(db))
	err = userservice.New(db).Add(testUser())
	require.NoError(t, err)
	tests := []struct {
		name    string
		opts    *userservice.SearchOptions
		want    []*user.User
		wantErr bool
	}{
		{
			name:    "nil searchoptions",
			want:    []*user.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by country",
			opts:    userservice.Search().Country("UK"),
			want:    []*user.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by nickname",
			opts:    userservice.Search().Nickname("johnny123"),
			want:    []*user.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by partial nickname",
			opts:    userservice.Search().Nickname("johnny"),
			want:    []*user.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by name",
			opts:    userservice.Search().Name("john", "smith"),
			want:    []*user.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by email",
			opts:    userservice.Search().Email("john.smith@faceit.com"),
			want:    []*user.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by partial email",
			opts:    userservice.Search().Email("john.smith"),
			want:    []*user.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by country",
			opts:    userservice.Search().Country("UK"),
			want:    []*user.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by country: country not found",
			opts:    userservice.Search().Country("FR"),
			want:    []*user.User{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := userservice.New(db)
			got, err := s.Get(tt.opts)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestService_Modify(t *testing.T) {
	type args struct {
		o *userservice.SearchOptions
		u *user.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "modify nil searchoptions",
			args: args{
				o: nil,
				u: &user.User{
					FirstName: "John",
					LastName:  "Smith",
					Nickname:  "Johnny123",
					Password:  encrypted("password"),
					Email:     "john.smith@faceit.com",
					Country:   "FR",
				},
			},
			wantErr: true,
		},
		{
			name: "modify nickname only",
			args: args{
				o: userservice.Search().Nickname("johnny123"),
				u: &user.User{
					FirstName: "John",
					LastName:  "Smith",
					Nickname:  "Johnny123",
					Password:  encrypted("password"),
					Email:     "john.smith@faceit.com",
					Country:   "FR",
				},
			},
			wantErr: true,
		},
		{
			name: "modify nickname & country",
			args: args{
				o: userservice.Search().Nickname("johnny123").Country("UK"),
				u: &user.User{
					FirstName: "John",
					LastName:  "Smith",
					Nickname:  "Johnny123",
					Password:  encrypted("password"),
					Email:     "john.smith@faceit.com",
					Country:   "FR",
				},
			},
			wantErr: false,
		},
		{
			name: "modify by email",
			args: args{
				o: userservice.Search().Email("john.smith@faceit.com"),
				u: &user.User{
					FirstName: "John",
					LastName:  "Smith",
					Nickname:  "Johnny123",
					Password:  encrypted("password"),
					Email:     "john.smith@faceit.com",
					Country:   "FR",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, err := sqlx.Connect("sqlite3", ":memory:")
			require.NoError(t, err)
			require.NoError(t, setup(db))
			err = userservice.New(db).Add(testUser())
			require.NoError(t, err)
			s := userservice.New(db)
			err = s.Modify(1, tt.args.u)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
			user := user.User{}
			require.NoError(t, db.Get(&user, `SELECT first_name, last_name, nickname, password, email, country FROM users`))
			require.Equal(t, tt.args.u, user)
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name     string
		o        *userservice.SearchOptions
		wantErr  bool
		wantUser user.User
	}{
		{
			name:    "delete nil searchoptions",
			o:       nil,
			wantErr: true,
			wantUser: user.User{
				FirstName: "John",
				LastName:  "Smith",
				Nickname:  "Johnny123",
				Password:  encrypted("password"),
				Email:     "john.smith@faceit.com",
				Country:   "UK",
			},
		},
		{
			name:    "delete nickname only",
			o:       userservice.Search().Nickname("johnny123"),
			wantErr: true,
			wantUser: user.User{
				FirstName: "John",
				LastName:  "Smith",
				Nickname:  "Johnny123",
				Password:  encrypted("password"),
				Email:     "john.smith@faceit.com",
				Country:   "UK",
			},
		},
		{
			name:     "delete nickname & country",
			o:        userservice.Search().Nickname("johnny123").Country("UK"),
			wantUser: user.User{},
			wantErr:  false,
		},
		{
			name:     "delete by email",
			o:        userservice.Search().Email("john.smith@faceit.com"),
			wantUser: user.User{},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, err := sqlx.Connect("sqlite3", ":memory:")
			require.NoError(t, err)
			require.NoError(t, setup(db))
			err = userservice.New(db).Add(testUser())
			require.NoError(t, err)

			s := userservice.New(db)
			err = s.Delete(1)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
			user := user.User{}
			err = db.Get(&user, `SELECT first_name, last_name, nickname, password, email, country FROM users`)
			if err != nil && strings.Contains(err.Error(), "no rows in result set") {
				err = nil
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantUser, user)
		})
	}
}

func setup(db *sqlx.DB) error {
	_, err := db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY,
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
	return err
}

func testUser() *user.User {
	return &user.User{
		FirstName: "John",
		LastName:  "Smith",
		Nickname:  "Johnny123",
		Password:  encrypted("password"),
		Email:     "john.smith@faceit.com",
		Country:   "UK",
	}
}

func encrypted(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	return string(hash.Sum(nil))
}
