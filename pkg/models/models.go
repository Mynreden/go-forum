package models

import "time"

type User struct {
	ID       int
	Name     string
	Email    string
	hashedPw []byte
	Created  time.Time
	Active   bool
	Password string
}
