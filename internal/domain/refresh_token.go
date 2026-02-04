package domain

import (
	"time"
)

// RefreshToken represents a refresh token stored in the database
type RefreshToken struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    string    `gorm:"not null;index;size:16" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;not null;type:text" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsRevoked bool      `gorm:"default:false" json:"is_revoked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for RefreshToken
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsExpired checks if the refresh token has expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsValid checks if the refresh token is valid (not expired and not revoked)
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.IsRevoked
}
