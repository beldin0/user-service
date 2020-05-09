package userservice

const sqlInsert = `INSERT INTO users
(
	first_name,
	first_name_lower,
	last_name,
	last_name_lower,
	nickname,
	nickname_lower,
	password,
	email,
	country
)
VALUES
(
	:first_name,
	:first_name_lower,
	:last_name,
	:last_name_lower,
	:nickname,
	:nickname_lower,
	:password,
	:email,
	:country
);`

const sqlGet = `SELECT first_name, last_name, nickname, password, email, country FROM users`

const sqlModify = `UPDATE users SET
	first_name=:first_name,
	first_name_lower=:first_name_lower,
	last_name=:last_name,
	last_name_lower=:last_name_lower,
	nickname=:nickname,
	nickname_lower=:nickname_lower,
	password=:password,
	email=:email,
	country=:country`
