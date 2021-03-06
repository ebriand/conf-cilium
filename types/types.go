package types

import "github.com/google/uuid"

type Hero struct {
	Name             string    `json:"name"`
	SecretIdentityID uuid.UUID `json:"secretIdentityID"`
	ImageURL         string    `json:"image_url"`
}

type Identity struct {
	ID       uuid.UUID `json:"id"`
	RealName string    `json:"realName"`
}
