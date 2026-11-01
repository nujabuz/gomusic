package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/nujabuz/gomusic/internal/entity"
	"gorm.io/gorm"
)

type userModel struct {
	ID        int64     `gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) entity.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	m := r.toModel(user)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return entity.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *userRepository) Find(ctx context.Context, user *entity.User) (*entity.User, error) {
	m := r.toModel(user)
	if err := r.db.WithContext(ctx).Where(m).First(m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}
	return r.toEntity(m), nil
}

// helper
func (r *userRepository) toModel(user *entity.User) *userModel {
	return &userModel{ID: user.ID, Username: user.Username, Password: user.Password, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}
}

func (r *userRepository) toEntity(user *userModel) *entity.User {
	return &entity.User{ID: user.ID, Username: user.Username, Password: user.Password, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}
}

func (userModel) TableName() string { return "users" }
