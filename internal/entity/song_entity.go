package entity

import (
	"context"
	"time"
)

type Song struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SongRepository Port
type SongRepository interface {
	Create(ctx context.Context, song *Song) error
	Update(ctx context.Context, song *Song) error
	Delete(ctx context.Context, id int64) error
	Find(ctx context.Context, id int64) (*Song, error)
}

// SongUseCase Port
type SongUseCase interface {
	Add(ctx context.Context, in *InputSong) (*Song, error)
	Update(ctx context.Context, in *InputSong) (*Song, error)
	Remove(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*Song, error)
}

// InputSong used in SongUseCase
type InputSong struct {
	ID   int64
	Name string
}
