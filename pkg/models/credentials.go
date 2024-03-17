package models

// Credentials используется для декодирования учетных данных пользователя из запроса логина
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
