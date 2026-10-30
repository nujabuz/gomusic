package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/nujabuz/music-cloud/internal/domain"
	"gorm.io/gorm"
)

type userModel struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Username  string    `gorm:"column:username;size:50;unique;not null"`
	Password  string    `gorm:"column:password;size:72;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (userModel) TableName() string {
	return "users"
}

type userRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	m := r.toModel(user)

	err := r.db.WithContext(ctx).Create(m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, domain.ErrUserAlreadyExists
		}
		return nil, err
	}
	d := r.toDomain(m)

	return d, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	m := r.toModel(user)

	res := r.db.WithContext(ctx).Model(m).Where("username = ?", m.Username).Updates(m)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domain.ErrUserNotFound
	}
	d := r.toDomain(m)

	return d, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	user := new(userModel)
	err := r.db.WithContext(ctx).Where("username = ?", username).Take(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	d := r.toDomain(user)

	return d, nil
}

func (r *userRepository) toModel(user *domain.User) *userModel {
	if user == nil {
		return nil
	}

	m := new(userModel)
	m.ID = user.ID
	m.Username = user.Username
	m.Password = user.Password
	m.CreatedAt = user.CreatedAt
	m.UpdatedAt = user.UpdatedAt

	return m
}
func (r *userRepository) toDomain(user *userModel) *domain.User {
	if user == nil {
		return nil
	}
	d := new(domain.User)
	d.ID = user.ID
	d.Username = user.Username
	d.Password = user.Password
	d.CreatedAt = user.CreatedAt
	d.UpdatedAt = user.UpdatedAt

	return d
}
