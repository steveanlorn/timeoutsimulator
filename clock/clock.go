// Package clock responsible to provide time data
package clock

import "time"

// Time is time interface
type Time interface {
	Now() time.Time
}

// Clock implements Time interface.
type Clock struct{}

// Now returns time.Now()
func (u Clock) Now() time.Time {
	return time.Now()
}
