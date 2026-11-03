package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/nujabuz/gomusic/internal/entity"
	"gorm.io/gorm"
)

type songModel struct {
	ID        int64 `gorm:"primarykey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type songRepository struct {
	db *gorm.DB
}

func NewSongRepository(db *gorm.DB) entity.SongRepository {
	return &songRepository{
		db: db,
	}
}

func (r *songRepository) Create(ctx context.Context, song *entity.Song) error {
	return r.db.WithContext(ctx).Create(song).Error
}

func (r *songRepository) Update(ctx context.Context, song *entity.Song) error {
	return r.db.WithContext(ctx).Save(song).Error
}

func (r *songRepository) Delete(ctx context.Context, id int64) error {
	return r.db.Where("id = ?", id).Delete(&songModel{}).Error
}

func (r *songRepository) Find(ctx context.Context, id int64) (*entity.Song, error) {
	song := new(songModel)
	err := r.db.Where("id = ?", id).First(song).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, entity.ErrSongNotFound
	}
	return &entity.Song{ID: song.ID, Name: song.Name, CreatedAt: song.CreatedAt, UpdatedAt: song.UpdatedAt}, err
}

func (songModel) TableName() string {
	return "songs"
}
