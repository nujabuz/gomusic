package rest

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/nujabuz/gomusic/internal/entity"
)

type RequestAddSong struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

type RequestUpdateSong struct {
	ID   int64  `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,min=2,max=255"`
}

type RequestDeleteSong struct {
	ID int64 `json:"id" validate:"required"`
}

type ResponseSong struct {
	ID        int64     `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
}

type SongHandler struct {
	SongUseCase     entity.SongUseCase
	TokenMiddleware fiber.Handler
}

func NewSongHandler(songUseCase entity.SongUseCase, tokenMiddleware fiber.Handler) *SongHandler {
	return &SongHandler{
		SongUseCase:     songUseCase,
		TokenMiddleware: tokenMiddleware,
	}
}

func (h *SongHandler) Routes(api fiber.Router) {
	song := api.Group("/songs")
	song.Use(h.TokenMiddleware)

	song.Post("/", h.add)
	song.Patch("/", h.update)
	song.Delete("/", h.delete)
	song.Get("/:id", h.getByID)
}

func (h *SongHandler) add(ctx fiber.Ctx) error {
	request := new(RequestAddSong)
	if err := ctx.Bind().Body(request); err != nil {
		return ResponseFailedParseRequest(ctx, err)
	}

	out, err := h.SongUseCase.Add(ctx, &entity.InputSong{Name: request.Name})
	if err != nil {
		return EntityErrorTranslator(ctx, err)
	}

	res := new(ResponseSong)
	res.Name = out.Name
	res.ID = out.ID
	res.CreatedAt = out.CreatedAt
	res.UpdatedAt = out.UpdatedAt

	return ResponseSuccess(ctx, fiber.StatusCreated, "successfully added song", res)
}

func (h *SongHandler) update(ctx fiber.Ctx) error {
	request := new(RequestUpdateSong)
	if err := ctx.Bind().Body(request); err != nil {
		return ResponseFailedParseRequest(ctx, err)
	}

	out, err := h.SongUseCase.Update(ctx, &entity.InputSong{ID: request.ID, Name: request.Name})
	if err != nil {
		return EntityErrorTranslator(ctx, err)
	}

	res := &ResponseSong{
		ID:        out.ID,
		Name:      out.Name,
		UpdatedAt: out.UpdatedAt,
		CreatedAt: out.CreatedAt,
	}

	return ResponseSuccess(ctx, fiber.StatusOK, "successfully update song", res)
}

func (h *SongHandler) delete(ctx fiber.Ctx) error {
	request := new(RequestDeleteSong)
	if err := ctx.Bind().Body(request); err != nil {
		return ResponseFailedParseRequest(ctx, err)
	}

	err := h.SongUseCase.Remove(ctx, request.ID)
	if err != nil {
		return EntityErrorTranslator(ctx, err)
	}

	return ResponseSuccess(ctx, fiber.StatusNoContent, "successfully delete song", nil)
}

func (h *SongHandler) getByID(ctx fiber.Ctx) error {
	songID, err := strconv.ParseInt(ctx.Params("id"), 0, 64)
	if err != nil {
		return ResponseFailed(ctx, fiber.StatusNotFound, entity.ErrSongNotFound.Error(), fiber.ErrNotFound)
	}

	out, err := h.SongUseCase.GetByID(ctx, songID)
	if err != nil {
		return EntityErrorTranslator(ctx, err)
	}

	res := &ResponseSong{
		ID:        out.ID,
		Name:      out.Name,
		CreatedAt: out.CreatedAt,
		UpdatedAt: out.UpdatedAt,
	}

	return ResponseSuccess(ctx, fiber.StatusOK, "successfully get song by id", res)
}
