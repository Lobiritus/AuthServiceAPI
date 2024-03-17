package handlers

import (
	"auth_service/internal/repository"
	"auth_service/pkg/models"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

// Функция регистрации пользователей
func RegisterUser(repo *repository.UserRepository, logger *zap.SugaredLogger) http.HandlerFunc {
	logger.Infow("Called Register user")
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			logger.Errorw("Failed to decode user", "error", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		logger.Infow("Check username and email", "username", "email", user.Username, user.Email)

		if len(user.Username) == 0 || len(user.Password) == 0 || len(user.Email) == 0 {
			logger.Errorw("User data is missing")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка в запросе: отсутствуют необходимые данные пользователя"})
			return
		}

		exists, err := repo.Exists(user.Username, user.Email)
		if err != nil {
			logger.Errorw("Failed to check if user exists", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if exists {
			logger.Infow("User already exists", "username", user.Username, "email", user.Email)
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		if err := repo.Create(&user, logger); err != nil {
			logger.Errorw("Failed to create user", "error", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		logger.Infow("User created successfully", "username", user.Username, "userID", user.ID)

		user.Password = ""
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь успешно зарегистрирован"})
	}

}
