package app

import (
	"gorm.io/gorm"
	"grest.dev/grest"
)

type mock struct{}

func Mock() mock {
	return mock{}
}

func (mock) DB() (*gorm.DB, error) {
	mockDB, _, mockErr := grest.NewMockDB()
	return mockDB, mockErr
}
