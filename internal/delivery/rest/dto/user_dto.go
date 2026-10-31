package dto

import (
	"time"
)

type RequestRegister struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type RequestLogin struct {
	Username string `json:"username" validate:"required,max=50"`
	Password string `json:"password" validate:"required,max=72"`
}

type ResponseUser struct {
	ID        int64     `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
}
