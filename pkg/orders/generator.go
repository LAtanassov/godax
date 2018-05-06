package orders

import (
	"strconv"
	"time"
)

// Generator generates a string
type Generator interface {
	Generate() string
}

// IDGenerator generates IDs
type IDGenerator struct {
}

// NewIDGenerator return an IDGenerator
func NewIDGenerator() Generator {
	return &IDGenerator{}
}

// Generate returns an ID (not unique)
func (i *IDGenerator) Generate() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
