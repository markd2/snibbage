package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record or 8-track found")
