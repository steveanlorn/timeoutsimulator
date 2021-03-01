// Package timeoutsimulator provides an instrument to simulate timeout budgeting
package timeoutsimulator

import (
	"context"
	"time"

	"github.com/steveanlorn/timeoutsimulator/clock"
	"github.com/steveanlorn/timeoutsimulator/deadline"
)

const (
	// RemarkFunctionNilContext denotes incoming context is nil.
	RemarkFunctionNilContext string = "Error nil context"

	// RemarkFunctionTimeout denotes function timeout is exceeded.
	RemarkFunctionTimeout string = "ctx budget %s < latency %s"

	// RemarkFunctionCtxBudget denotes information of context budget.
	RemarkFunctionCtxBudget string = "ctx budget %s"

	// RemarkSimulatorTimeoutExceeded denotes simulator timeout exceeded.
	RemarkSimulatorTimeoutExceeded = "Exceed simulator timeout"
)

// Function denotes accepted function to run in simulator
type Function interface {
	// Run denotes function execution behavior.
	// It returns bool indicates whether function is exceeded the timeout.
	Run(ctx context.Context) bool

	// String returns function name.
	String() string

	// Latency returns function execution time.
	Latency() time.Duration

	// Remark returns a note about the function if any.
	Remark() string
}

// Result denotes data structure for simulator result.
type Result struct {
	// Name is the simulator name.
	Name string

	// Timeout is the simulator timeout.
	TimeoutDuration time.Duration

	// Data is a list contains Function data execution.
	Data []Data
}

// Data denotes Function data execution detail.
type Data struct {
	// Name is function name.
	Name string

	// Latency is function execution time.
	Latency time.Duration

	// IsDeadlineExceeded true if reach time out before function completion.
	IsDeadlineExceeded bool

	// TimeIn is incoming context deadline
	TimeIn time.Duration

	// TimeOut is outcoming context deadline.
	// Subtracted with function execution time.
	TimeOut time.Duration

	// Remark contains a note derived from the function.
	Remark string
}

type function struct {
	Function
	isExecuted bool
}

// Simulator is the timeout budgeting simulator
type Simulator struct {
	name            string
	timeoutDuration time.Duration
	functions       []function
	deadline        deadline.Deadline
}

// NewSimulator initialize new Simulator.
// Set t to zero to run simulator without timeout.
func NewSimulator(n string, t time.Duration, fs []Function) Simulator {
	functions := make([]function, 0, len(fs))
	for _, f := range fs {
		functions = append(functions, function{Function: f})
	}

	return Simulator{
		name:            n,
		timeoutDuration: t,
		functions:       functions,
		deadline:        deadline.NewDeadline(clock.Clock{}),
	}
}

// Start runs the simulator.
// When Start is being called, it will execute each function in order.
func (s Simulator) Start() Result {
	switch s.timeoutDuration {
	case 0:
		return s.run()
	default:
		return s.runWithContext()
	}
}

func (s Simulator) run() Result {
	ctx := context.Background()
	data := make([]Data, 0, len(s.functions))

	for i, f := range s.functions {
		isDeadlineExceeded := f.Run(ctx)

		s.functions[i].isExecuted = true

		data = append(data, Data{
			Name:               f.String(),
			Latency:            f.Latency(),
			IsDeadlineExceeded: isDeadlineExceeded,
			Remark:             f.Remark(),
		})
	}

	return Result{
		Name:            s.name,
		TimeoutDuration: s.timeoutDuration,
		Data:            data,
	}
}

func (s Simulator) runWithContext() Result {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeoutDuration)
	defer cancel()

	data := make([]Data, 0, len(s.functions))

	for i, f := range s.functions {
		timeIn := s.deadline.GetInMillis(ctx)
		isDeadlineExceeded := f.Run(ctx)
		timeOut := s.deadline.GetInMillis(ctx)

		s.functions[i].isExecuted = true

		data = append(data, Data{
			Name:               f.String(),
			Latency:            f.Latency(),
			IsDeadlineExceeded: isDeadlineExceeded,
			TimeIn:             timeIn,
			TimeOut:            timeOut,
			Remark:             f.Remark(),
		})

		if timeOut == 0 {
			break
		}
	}

	for _, f := range s.functions {
		if f.isExecuted {
			continue
		}

		data = append(data, Data{
			Name:               f.String(),
			Latency:            f.Latency(),
			IsDeadlineExceeded: true,
			Remark:             RemarkSimulatorTimeoutExceeded,
		})
	}

	return Result{
		Name:            s.name,
		TimeoutDuration: s.timeoutDuration,
		Data:            data,
	}
}
