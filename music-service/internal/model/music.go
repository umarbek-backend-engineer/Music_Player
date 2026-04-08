package model

import (
	"time"

	"github.com/google/uuid"
)

type Music struct {
	ID          uuid.UUID `json:"id,omitempty"`
	FileName    string    `json:"filename"`
	FilePath    string    `json:"filepath"`
	Uploadad_at time.Time `json:"uploaded_at"`
}
