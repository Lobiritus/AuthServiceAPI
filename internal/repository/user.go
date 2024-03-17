package repository

import (
	"auth_service/pkg/models"
	"auth_service/pkg/utils"
	"database/sql"
	"go.uber.org/zap"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// SQL запрос для вставки новой записи в таблицу пользователей
// Возвращает ошибку, если не удалось добавить пользователя
func (r *UserRepository) Create(user *models.User, logger *zap.SugaredLogger) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (username,password,email) VALUES ($1, $2, $3) RETURNING ID`
	err = r.db.QueryRow(query, user.Username, hashedPassword, user.Email).Scan(&user.ID)

	if err != nil {
		// В последствии можно добавить логирование
		logger.Errorw("Failed to insert data")
		return err
	}
	return nil
}

func (r *UserRepository) Exists(username, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)`
	err := r.db.QueryRow(query, username, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username,password,email FROM users WHERE username =$1`
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
