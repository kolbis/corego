package utils

import "github.com/google/uuid"

// NewUUID returns new uuid as string
func NewUUID() string {
	return uuid.New().String()
}
