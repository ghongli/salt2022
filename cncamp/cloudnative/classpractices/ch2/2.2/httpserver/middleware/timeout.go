package middleware

import (
	"time"
)

const (
	DefaultTimeout = time.Second * 15
	
	// add to give a handler that respects the deadline the opportunity to return an error.
	boost = time.Millisecond * 50
)