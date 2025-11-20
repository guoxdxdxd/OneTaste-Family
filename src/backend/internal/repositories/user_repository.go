package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"onetaste-family/backend/internal/models"
	"onetaste-family/backend/pkg/database"
)

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")
	// ErrUserExists 用户已存在
	ErrUserExists = errors.New("user already exists")
)

// UserRepository 用户数据访问层
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository 创建用户Repository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.GetDB(),
	}
}

// Create 创建用户
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, phone, password, nickname, avatar, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.ID,
		user.Phone,
		user.Password,
		user.Nickname,
		user.Avatar,
		user.Status,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// 检查是否是唯一约束冲突
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_phone_key\"" {
			return ErrUserExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByPhone 根据手机号获取用户
func (r *UserRepository) GetByPhone(phone string) (*models.User, error) {
	query := `
		SELECT id, phone, password, nickname, avatar, status, created_at, updated_at
		FROM users
		WHERE phone = $1 AND status = 1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, phone).Scan(
		&user.ID,
		&user.Phone,
		&user.Password,
		&user.Nickname,
		&user.Avatar,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return user, nil
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, phone, password, nickname, avatar, status, created_at, updated_at
		FROM users
		WHERE id = $1 AND status = 1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Phone,
		&user.Password,
		&user.Nickname,
		&user.Avatar,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

// ExistsByPhone 检查手机号是否已存在
func (r *UserRepository) ExistsByPhone(phone string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)`
	
	var exists bool
	err := r.db.QueryRow(query, phone).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user exists: %w", err)
	}

	return exists, nil
}
