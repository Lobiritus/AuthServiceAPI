package models

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"` // так как pet project, храним не в хеш виде
	Email    string `json:"email"`
}
