package models

import "time"

type User struct {
	ID        int       `db:"id"`
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
