package model

import "github.com/google/uuid"

type Supplier struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	ContactInfo string    `json:"contact_info"`
}
