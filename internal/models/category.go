package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `jsonb:"id" db:"id"`
	Name        string    `jsonb:"name" db:"name"`
	Description *string   `jsonb:"description" db:"description"`
	Color       *string   `jsonb:"color" db:"color"`
	IsActive    bool      `jsonb:"isActive" db:"isActive"`
	CreatedAt   time.Time `jsonb:"createdAt" db:"createdAt"`
	UpdatedAt   time.Time `jsonb:"updatedAt" db:"updatedAt"`
}

type CreateCategoryRequest struct {
	Name        string  `jsonb:"name" validate:"required,max=100"`
	Description *string `jsonb:"description"`
	Color       *string `jsonb:"color" validate:"omitempty,len=7"`
}

type UpdateCategoryRequest struct {
	Name        *string `jsonb:"name" validate:"omitempty,max=100"`
	Description *string `jsonb:"description"`
	Color       *string `jsonb:"color" validate:"omitempty,len=7"`
	IsActive    *bool   `jsonb:"isActive"`
}

// JSONB is a custom type for handling PostgreSQL JSONB fields
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface for database/sql
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for database/sql
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}