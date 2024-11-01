// Package metrics constructs the metrics the application will track.
package metrics

import (
	"context"
	"expvar"
	"runtime"
)

// This holds the single instance of metrics value needed for collecting metrics.
// Basing this on the expvar package singleton layout.
var m metrics

// metrics represents the set of metrics gathered.
type metrics struct {
	goroutines *expvar.Int
	requests   *expvar.Int
	errors     *expvar.Int
	panics     *expvar.Int
}

// init function constructs metrics value to be used for capturing metrics.
// This value is stored inside the package level variable, since everything
// inside of expvar is stored as a singleton.
func init() {
	m = metrics{
		goroutines: expvar.NewInt("goroutines"),
		requests:   expvar.NewInt("requests"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
	}
}

type ctxKey int

const key ctxKey = 1

// Set sets the metrics data into the context.
func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, &m)
}

// AddGoroutines refreshes the goroutine metric.
func AddGoroutines(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		g := int64(runtime.NumGoroutine())
		v.goroutines.Set(g)
		return g
	}
	return 0
}

// AddRequests increments the request metric by 1.
func AddRequests(ctx context.Context) int64 {
	v, ok := ctx.Value(key).(*metrics)
	if ok {
		v.requests.Add(1)
		return v.requests.Value()
	}
	return 0
}

// AddErrors increments the errors metric by 1.
func AddErrors(ctx context.Context) int64 {
	v, ok := ctx.Value(key).(*metrics)
	if ok {
		v.errors.Add(1)
		return v.errors.Value()
	}
	return 0
}

// AddPanics increments the panics metric by 1.
func AddPanics(ctx context.Context) int64 {
	v, ok := ctx.Value(key).(*metrics)
	if ok {
		v.panics.Add(1)
		return v.panics.Value()
	}
	return 0
}
