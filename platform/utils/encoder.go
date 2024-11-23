package utils

import (
	"encoding/gob"
	"time"
)

func init() {
	// Register necessary types with gob
	gob.Register(map[string]interface{}{})
	gob.Register(time.Time{}) // Register time.Time
}
