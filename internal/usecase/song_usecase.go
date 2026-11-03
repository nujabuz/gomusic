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
	song := new(entity.Song)
	song.Name = in.Name

	if err := u.songRepository.Create(ctx, song); err != nil {
		return nil, err
	}

	return song, nil
}

func (u *songUseCase) Remove(ctx context.Context, id int64) error {
	return u.songRepository.Delete(ctx, id)
}

func (u *songUseCase) Update(ctx context.Context, in *entity.InputSong) (*entity.Song, error) {
	song := new(entity.Song)
	song.Name = in.Name
	song.ID = in.ID

	if _, err := u.songRepository.Find(ctx, in.ID); err != nil {
		return nil, err
	}

	if err := u.songRepository.Update(ctx, song); err != nil {
		return nil, nil
	}

	return song, nil
}

func (u *songUseCase) GetByID(ctx context.Context, id int64) (*entity.Song, error) {
	song, err := u.songRepository.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	return song, nil
}
