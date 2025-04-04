package health

import (
	"cmp"
	"context"
	"slices"
	"sync"
	"time"

	"github.com/webitel/webitel-go-kit/logging/wlog"
)

// CheckRegistry is a registry of health checks from the API and Infra SDKs
// and other parts of the runtime.
type CheckRegistry struct {
	log    *wlog.Logger
	m      sync.Mutex
	checks []Check
}

// NewCheckRegistry creates a new CheckRegistry.
func NewCheckRegistry(log *wlog.Logger) *CheckRegistry {
	return &CheckRegistry{log: log}
}

// Register registers a new health check.
//
// Checks must complete within 5 seconds, otherwise
// they will be terminated and considered failed.
//
// Checks can be called at any time and could have
// multiple goroutines calling them concurrently.
func (c *CheckRegistry) Register(check Check) {
	c.m.Lock()
	defer c.m.Unlock()
	c.checks = append(c.checks, check)
}

// RegisterFunc registers a new health check from a function with a given name
//
// This is a convinced wrapper over [CheckRegistry.Register], see that function
// for more details and expected behavior.
func (c *CheckRegistry) RegisterFunc(name string, check func(ctx context.Context) error) {
	c.Register(&checkFunc{name, check})
}

// GetChecks returns all registered health checks.
func (c *CheckRegistry) GetChecks() []Check {
	c.m.Lock()
	defer c.m.Unlock()
	return c.checks
}

// RunAll runs all health checks and returns the results.
func (c *CheckRegistry) RunAll(ctx context.Context) []CheckResult {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	checks := c.GetChecks()

	// Run all checks in parallel.
	results := make(chan []CheckResult, len(checks))
	var wg sync.WaitGroup
	wg.Add(len(checks))
	for _, check := range checks {
		check := check
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					c.log.Error("health check resulted in a panic", wlog.Any("panic", r))
				}
			}()

			result := check.HealthCheck(ctx)
			for _, v := range result {
				c.log.Info("health check result", wlog.Any("result", v))
			}

			results <- result
		}()
	}

	// Wait for all checks to complete for the context to be canceled.
	var allResults []CheckResult

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait for all checks to complete or the context to be canceled.
	select {
	case <-done:
	case <-ctx.Done():
		allResults = append(allResults, CheckResult{
			Name: "health-checks.run",
			Err:  ctx.Err(),
		})
	}
	close(results) // then close the results channel

	// Collect results.
	for results := range results {
		allResults = append(allResults, results...)
	}

	// Sort results by name.
	slices.SortFunc(allResults, func(a, b CheckResult) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return allResults
}
