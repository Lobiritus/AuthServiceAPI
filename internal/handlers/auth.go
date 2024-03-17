package handlers

import (
	"auth_service/internal/repository"
	"auth_service/pkg/models"
	"auth_service/pkg/utils"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

// LoginHandler обрабатывает вход пользователя
func LoginHandler(repo *repository.UserRepository, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds models.Credentials // Структура для учетных данных

		// Декодирование учетных данных из структуры запроса
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			logger.Errorw("Error decoding credentials", "error", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Получение пользователя из базы данных
		user, err := repo.GetUserByUsername(creds.Username)
		if err != nil {
			logger.Errorw("User not found", "error", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Проверка пароля
		if !utils.CheckPasswordHash(creds.Password, user.Password) {
			logger.Error("Invalid password")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// генерация JWT токена для пользователя
		token, err := utils.GenerateJWT(user.Username)
		if err != nil {
			logger.Errorw("Error generating JWT", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Отправка токена в ответе
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})

	}
}
