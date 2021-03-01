package context

import (
	"context"
	"time"

	"github.com/steveanlorn/timeoutsimulator/clock"
	"github.com/steveanlorn/timeoutsimulator/deadline"
	"github.com/steveanlorn/timeoutsimulator/function"
)

// WeightedContext denotes a context with weight.
// Weight value will be multiplied by incoming context deadline,
// producing a new context time out with that multiplication value.
//
// If priority is true and the new timeout is lower than timeout threshold,
// it will use all incoming context deadline time.
type WeightedContext struct {
	timeoutThreshold time.Duration
	weight           float64
	isPriority       bool
	budget           time.Duration
	deadline         deadline.Deadline
}

var _ function.Context = (*WeightedContext)(nil)

// NewWeightedContext initializes new WeightedContext.
func NewWeightedContext(timeoutThreshold time.Duration, weight float64, isPriority bool) *WeightedContext {
	return &WeightedContext{
		timeoutThreshold: timeoutThreshold,
		weight:           weight,
		isPriority:       isPriority,
		deadline:         deadline.NewDeadline(clock.Clock{}),
	}
}

// Get returns new weighted context from incoming context.
func (w *WeightedContext) Get(ctx context.Context) (context.Context, context.CancelFunc) {
	switch ctx {
	case nil:
		return ctx, func() {}
	default:
		timeout := w.deadline.GetInMillis(ctx)

		newTimeout := time.Duration(float64(timeout) * w.weight)
		if newTimeout < w.timeoutThreshold && w.isPriority {
			newTimeout = timeout
		}

		w.budget = newTimeout

		return context.WithTimeout(ctx, newTimeout)
	}
}

// Budget returns context budget duration in millisecond.
func (w *WeightedContext) Budget() time.Duration {
	return w.budget
}

// TimeOutContext denotes context with timeout.
type TimeOutContext struct {
	timeOut time.Duration
	budget  time.Duration
}

var _ function.Context = (*TimeOutContext)(nil)

// NewTimeOutContext initializes new NewTimeOutContext.
func NewTimeOutContext(timeOut time.Duration) *TimeOutContext {
	return &TimeOutContext{
		timeOut: timeOut,
	}
}

// Get returns new context with timeout.
func (t *TimeOutContext) Get(ctx context.Context) (context.Context, context.CancelFunc) {
	switch ctx {
	case nil:
		return ctx, func() {}
	default:
		t.budget = t.timeOut
		return context.WithTimeout(ctx, t.timeOut)
	}
}

// Budget returns context budget duration in millisecond.
func (t *TimeOutContext) Budget() time.Duration {
	return t.budget
}
