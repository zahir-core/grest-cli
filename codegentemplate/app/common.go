package app

import (
	"strings"

	"github.com/google/uuid"
)

func NewUUID() string {
	return uuid.NewString()
}

func NewToken() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
