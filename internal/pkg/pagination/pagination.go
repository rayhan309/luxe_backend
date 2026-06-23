package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

// Params holds pagination query params
type Params struct {
	Page  int
	Limit int
}

// FromContext extracts pagination params from the Gin context
func FromContext(c *gin.Context) Params {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = DefaultPage
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	return Params{Page: page, Limit: limit}
}

// Offset returns the MongoDB skip value
func (p Params) Offset() int64 {
	return int64((p.Page - 1) * p.Limit)
}
