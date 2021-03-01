package function

import (
	"context"
	"fmt"
	"time"

	"github.com/steveanlorn/timeoutsimulator"
)

// Context denotes context interface.
type Context interface {
	// Get returns context and its cancel func.
	Get(ctx context.Context) (context.Context, context.CancelFunc)

	// Budget returns context budget time duration
	Budget() time.Duration
}

// SleepFunction a function that will sleep
// and returns when it awakes or when context is canceled.
//
// Set context to nil to denote a function without timeout handler.
type SleepFunction struct {
	name    string
	sleep   time.Duration
	context Context
	remark  string
}

// NewSleepFunction initialize SleepFunction.
func NewSleepFunction(name string, sleep time.Duration, context Context) *SleepFunction {
	return &SleepFunction{
		name:    name,
		sleep:   sleep,
		context: context,
	}
}

// Run runs the SleepFunction.
func (s *SleepFunction) Run(ctx context.Context) bool {
	switch s.context {
	case nil:
		return s.run()
	default:
		return s.runWithContext(ctx)
	}
}

func (s *SleepFunction) run() bool {
	time.Sleep(s.sleep)
	return false
}

func (s *SleepFunction) runWithContext(ctx context.Context) bool {
	if ctx == nil {
		s.remark = timeoutsimulator.RemarkFunctionNilContext
		return false
	}

	ctx, cancel := s.context.Get(ctx)

	done := make(chan struct{}, 1)
	go func() {
		time.Sleep(s.sleep)
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		s.remark = fmt.Sprintf(timeoutsimulator.RemarkFunctionTimeout, s.context.Budget().String(), s.Latency().String())
		cancel()
		return true
	case <-done:
		s.remark = fmt.Sprintf(timeoutsimulator.RemarkFunctionCtxBudget, s.context.Budget().String())
		cancel()
		return false
	}
}

// String returns function name
func (s *SleepFunction) String() string {
	return s.name
}

// Latency returns sleep time duration.
func (s *SleepFunction) Latency() time.Duration {
	return s.sleep
}

// Remark returns note about the function.
func (s *SleepFunction) Remark() string {
	return s.remark
}
