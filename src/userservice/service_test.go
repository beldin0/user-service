package userservice_test

import (
	"crypto/sha256"
	"testing"

	"github.com/beldin0/users/src/userservice"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestService_Add(t *testing.T) {
	tests := []struct {
		name    string
		u       userservice.User
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
			user := userservice.User{}
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
		want    []userservice.User
		wantErr bool
	}{
		{
			name:    "nil searchoptions",
			want:    []userservice.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by country",
			opts:    userservice.Search().Country("UK"),
			want:    []userservice.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by nickname",
			opts:    userservice.Search().Nickname("johnny123"),
			want:    []userservice.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by partial nickname",
			opts:    userservice.Search().Nickname("johnny"),
			want:    []userservice.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by name",
			opts:    userservice.Search().Name("john", "smith"),
			want:    []userservice.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by email",
			opts:    userservice.Search().Email("john.smith@faceit.com"),
			want:    []userservice.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by partial email",
			opts:    userservice.Search().Email("john.smith"),
			want:    []userservice.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by country",
			opts:    userservice.Search().Country("UK"),
			want:    []userservice.User{testUser()},
			wantErr: false,
		},
		{
			name:    "search by country: country not found",
			opts:    userservice.Search().Country("FR"),
			want:    []userservice.User{},
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
		u userservice.User
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
				u: userservice.User{
					Firstname: "John",
					Lastname:  "Smith",
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
				u: userservice.User{
					Firstname: "John",
					Lastname:  "Smith",
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
				u: userservice.User{
					Firstname: "John",
					Lastname:  "Smith",
					Nickname:  "Johnny123",
					Password:  encrypted("password"),
					Email:     "john.smith@faceit.com",
					Country:   "FR",
				},
			},
			wantErr: true,
		},
		{
			name: "modify by email",
			args: args{
				o: userservice.Search().Email("john.smith@faceit.com"),
				u: userservice.User{
					Firstname: "John",
					Lastname:  "Smith",
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
			err = s.Modify(tt.args.o, tt.args.u)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
			user := userservice.User{}
			require.NoError(t, db.Get(&user, `SELECT first_name, last_name, nickname, password, email, country FROM users`))
			require.Equal(t, tt.args.u, user)
		})
	}
}

func setup(db *sqlx.DB) error {
	_, err := db.Exec(`CREATE TABLE users (
		first_name VARCHAR(50),
		first_name_lower VARCHAR(50),
		last_name VARCHAR(50),
		last_name_lower VARCHAR(50),
		nickname VARCHAR(30),
		nickname_lower VARCHAR(30) UNIQUE,
		password VARCHAR(32),
		email VARCHAR(50) UNIQUE PRIMARY KEY,
		country VARCHAR(3)
	)`)
	return err
}

func testUser() userservice.User {
	return userservice.User{
		Firstname: "John",
		Lastname:  "Smith",
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
