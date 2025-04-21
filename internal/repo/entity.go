package repo

import "time"

type User struct {
	ID         string    `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	Password   string    `json:"password" db:"password"`
	Created_at time.Time `json:"created_at" db:"created_at"`
	Updated_at time.Time `json:"updated_at" db:"updated_at"`
}
