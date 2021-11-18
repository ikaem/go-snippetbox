// pkg/models/models.go

package models

import (
	"errors"
	"time"
)

// this is an error type
var ErrNoRecord = errors.New("models: no matching record found")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
