package timeoutsimulator_test

import (
	"os"
	"time"

	"github.com/steveanlorn/timeoutsimulator"
	"github.com/steveanlorn/timeoutsimulator/context"
	"github.com/steveanlorn/timeoutsimulator/function"
	"github.com/steveanlorn/timeoutsimulator/servetable"
)

func ExampleSimulator_Start() {
	timeoutThreshold := 30 * time.Millisecond

	functions := []timeoutsimulator.Function{
		function.NewSleepFunction("func 1", 10*time.Millisecond, context.NewWeightedContext(timeoutThreshold, 0.2, false)),
		function.NewSleepFunction("func 2", 80*time.Millisecond, context.NewWeightedContext(timeoutThreshold, 0.5, false)),
		function.NewSleepFunction("func 3", 20*time.Millisecond, context.NewWeightedContext(timeoutThreshold, 0.3, false)),
		function.NewSleepFunction("func 4", 50*time.Millisecond, context.NewWeightedContext(timeoutThreshold, 0.6, true)),
		function.NewSleepFunction("func 5", 30*time.Millisecond, context.NewTimeOutContext(30*time.Millisecond)),
	}

	sim := timeoutsimulator.NewSimulator("my simulator", 100*time.Millisecond, functions)
	result := sim.Start()

	_ = servetable.Generate(result, os.Stdout)
	// Output:
	// ================================
	// SIMULATOR:my simulator
	// TIMEOUT DURATION:100ms
	// --------------------------------
	// NAME   |LATENCY |IS TIMEOUT |TIME IN |TIME OUT |NOTE                             |
	// func 1 |10ms    |false      |99ms    |89ms     |ctx budget 19.8ms                |
	// func 2 |80ms    |true       |89ms    |44ms     |ctx budget 44.5ms < latency 80ms |
	// func 3 |20ms    |true       |44ms    |31ms     |ctx budget 13.2ms < latency 20ms |
	// func 4 |50ms    |true       |31ms    |0s       |ctx budget 31ms < latency 50ms   |
	// func 5 |30ms    |true       |0s      |0s       |Exceed simulator timeout         |
	// ================================
}
