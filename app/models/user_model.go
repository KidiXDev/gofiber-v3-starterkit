package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            uuid.UUID  `bun:"id,type:uuid,default:gen_random_uuid(),pk" json:"id"`
	Email         string     `bun:"email,unique,notnull" json:"email"`
	PasswordHash  *string    `bun:"password_hash" json:"-"`
	Username      string     `bun:"username,unique,notnull" json:"username"`
	Avatar        *string    `bun:"avatar" json:"avatar"`
	Bio           *string    `bun:"bio" json:"bio"`
	DeletedAt     *time.Time `bun:"deleted_at,soft_delete" json:"deletedAt"`
	CreatedAt     time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updatedAt"`
}
