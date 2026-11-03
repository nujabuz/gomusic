package usecase

import (
	"context"

	"github.com/nujabuz/gomusic/internal/entity"
)

type songUseCase struct {
	songRepository entity.SongRepository
}

func NewSongUseCase(songrepository entity.SongRepository) entity.SongUseCase {
	return &songUseCase{
		songRepository: songrepository,
	}
}

func (u *songUseCase) Add(ctx context.Context, in *entity.InputSong) (*entity.Song, error) {
	if err := u.songRepository.Create(ctx, &entity.Song{Name: in.Name}); err != nil {
		return nil, err
	}

	return &entity.Song{Name: in.Name}, nil
}

func (u *songUseCase) Remove(ctx context.Context, id int64) error {
	return u.songRepository.Delete(ctx, id)
}

func (u *songUseCase) GetByID(ctx context.Context, id int64) (*entity.Song, error) {
	song, err := u.songRepository.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	return song, nil
}
