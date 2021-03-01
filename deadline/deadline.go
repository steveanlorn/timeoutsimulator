package deadline

import (
	"context"
	"time"

	"github.com/steveanlorn/timeoutsimulator/clock"
)

// Deadline denotes deadline struct.
type Deadline struct {
	clock clock.Time
}

// NewDeadline initialize new Deadline.
func NewDeadline(clock clock.Time) Deadline {
	return Deadline{
		clock: clock,
	}
}

// GetInMillis returns context deadline in millisecond unit.
func (d Deadline) GetInMillis(ctx context.Context) time.Duration {
	deadline, _ := ctx.Deadline()

	unixTime := deadline.UnixNano()
	diffTime := unixTime - d.clock.Now().UnixNano()
	diffTime = diffTime / 1e6

	return time.Duration(diffTime) * time.Millisecond
}
