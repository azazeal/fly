// Package health implements a http.Handler useful for when implementing HTTP
// health checks for fly.io apps.
package health

import (
	"context"
	"net/http"
	"sync"
)

type contextKeyType int

const contextKey contextKeyType = iota + 1

// NewContext derives a Context which carries the given Check from the given
// one.
func NewContext(ctx context.Context, c *Check) context.Context {
	return context.WithValue(ctx, contextKey, c)
}

// FromContext returns the Check the given Context carries.
//
// FromContext panics in case the given Context carries no Check.
func FromContext(ctx context.Context) *Check {
	return ctx.Value(contextKey).(*Check)
}

// Check implements a health check as the logical summation of boolean named
// components.
type Check struct {
	mu         sync.RWMutex
	components map[string]struct{}
}

// Set sets the given components of s.
func (c *Check) Set(components ...string) {
	if len(components) == 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, component := range components {
		delete(c.components, component)
	}
}

// Unset unsets the given components of s.
func (c *Check) Unset(components ...string) {
	if len(components) == 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.components == nil {
		c.components = make(map[string]struct{}, len(components))
	}

	for _, component := range components {
		c.components[component] = struct{}{}
	}
}

// Healthy reports the logical AND of all set/unset components of s.
//
// A Check on which no component has ever been set or unset is always healthy.
func (c *Check) Healthy() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.components) == 0
}

// ServeHTTP implements http.Handler for s.
//
// ServeHTTP responds with http.StatusServiceUnavailable when s is unhealthy and
// with http.StatusNoContent (HEAD) or http.StatusOK (GET)
// when it's not.
//
// In all other cases the ServeHTTP responds with http.StatusMethodNotAllowed.
func (c *Check) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	switch r.Method {
	default:
		respondWith(wr, http.StatusMethodNotAllowed)

		return
	case http.MethodGet, http.MethodHead:
		break
	}

	healthy := c.Healthy()

	if r.Method == http.MethodHead {
		if healthy {
			wr.WriteHeader(http.StatusNoContent)
		} else {
			wr.WriteHeader(http.StatusServiceUnavailable)
		}

		return
	}

	if healthy {
		respondWith(wr, http.StatusOK)

		return
	}

	respondWith(wr, http.StatusServiceUnavailable)
}

func respondWith(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
