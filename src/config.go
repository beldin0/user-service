package main

import "fmt"

type config struct {
	Host     string `envconfig:"POSTGRES_HOST" default:"db"`
	Port     int    `envconfig:"POSTGRES_PORT" default:"5432"`
	User     string `envconfig:"POSTGRES_USER" default:"postgres"`
	Password string `envconfig:"POSTGRES_PASSWORD" default:"password"`
	DBName   string `envconfig:"POSTGRES_DB_NAME" default:"postgres"`
	SSLMode  bool   `envconfig:"POSTGRES_SSLMODE"`
}

func (c config) ConnString() string {
	return fmt.Sprintf(
		`host=%s port=%d user=%s password=%s dbname=%s sslmode=%s`,
		c.Host, c.Port, c.User, c.Password, c.DBName, func() string {
			if c.SSLMode {
				return "enable"
			}
			return "disable"
		}())
}
