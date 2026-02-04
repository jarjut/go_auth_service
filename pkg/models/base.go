package models

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

// BaseModelNanoID provides a base model with NanoID primary key
// Uses 16-character unambiguous alphanumeric alphabet (removes 0, 1, I, O, l, o)
type BaseModelNanoID struct {
	ID string `gorm:"primaryKey;size:16" json:"id"`
}

// BeforeCreate generates a NanoID before creating the entity
func (m *BaseModelNanoID) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		// Use unambiguous alphabet: no 0, 1, I, O, l, o
		m.ID, _ = gonanoid.Generate("23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz", 16)
	}
	return
}
