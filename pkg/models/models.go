package models

import "time"

type User struct {
	ID       int
	Name     string
	Email    string
	HashedPw []byte
	Created  time.Time
	Active   bool
}
