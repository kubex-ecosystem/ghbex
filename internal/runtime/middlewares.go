package runtime

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"
)

type Middleware func(Operator) Operator

func Chain(op Operator, mws ...Middleware) Operator {
	for i := len(mws) - 1; i >= 0; i-- {
		op = mws[i](op)
	}
	return op
}

// Recorder de métricas.
type Recorder func(fields map[string]any)

// WithMeter registra métricas de duração, erro e cache_hit.
func WithMeter(rec Recorder) Middleware {
	return func(next Operator) Operator {
		return opFunc{
			name:    next.Name(),
			version: next.Version(),
			run: func(ctx context.Context, in OpInput) (OpOutput, error) {
				start := time.Now()
				out, err := next.Run(ctx, in)
				if rec != nil {
					rec(map[string]any{
						"op":        next.Name(),
						"ver":       next.Version(),
						"dur_ms":    time.Since(start).Milliseconds(),
						"err":       err,
						"repo":      in.Repo,
						"cache_hit": lookupMetric(out.Metrics, "cache_hit"),
					})
				}
				return out, err
			},
		}
	}
}

func lookupMetric(ms []Metric, name string) float64 {
	for _, m := range ms {
		if m.Name == name {
			return m.Value
		}
	}
	return 0
}

// WithTimeout limite por operador.
func WithTimeout(d time.Duration) Middleware {
	return func(next Operator) Operator {
		return opFunc{
			name:    next.Name(),
			version: next.Version(),
			run: func(ctx context.Context, in OpInput) (OpOutput, error) {
				if d <= 0 {
					return next.Run(ctx, in)
				}
				c, cancel := context.WithTimeout(ctx, d)
				defer cancel()
				return next.Run(c, in)
			},
		}
	}
}

// WithRetry backoff exponencial com jitter e política de re-tentativa.
func WithRetry(max int, base time.Duration, shouldRetry func(error) bool) Middleware {
	if max < 1 {
		max = 1
	}
	if base <= 0 {
		base = 200 * time.Millisecond
	}
	if shouldRetry == nil {
		shouldRetry = func(err error) bool {
			if err == nil {
				return false
			}
			return !errors.Is(err, ErrBudgetExceeded) && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded)
		}
	}
	return func(next Operator) Operator {
		return opFunc{
			name:    next.Name(),
			version: next.Version(),
			run: func(ctx context.Context, in OpInput) (OpOutput, error) {
				var out OpOutput
				var err error
				for attempt := 0; attempt < max; attempt++ {
					out, err = next.Run(ctx, in)
					if err == nil || !shouldRetry(err) {
						return out, err
					}
					// backoff
					d := base * (1 << attempt)
					d = d + time.Duration(rand.Int63n(int64(d/3)+1)) // jitter
					select {
					case <-ctx.Done():
						return out, ctx.Err()
					case <-time.After(d):
					}
				}
				return out, err
			},
		}
	}
}

// ===== Budget =====

type Budget struct {
	mu   sync.Mutex
	max  float64
	used float64
}

func NewBudget(max float64) *Budget { return &Budget{max: max} }
func (b *Budget) Left() float64     { b.mu.Lock(); defer b.mu.Unlock(); return b.max - b.used }
func (b *Budget) charge(x float64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.used+x > b.max {
		return false
	}
	b.used += x
	return true
}

// WithBudget aplica hard-stop por custo.
// costFn extrai custo do OpOutput (ex.: somar métricas com name=="cost_usd").
func WithBudget(b *Budget, costFn func(OpOutput) float64) Middleware {
	if b == nil || costFn == nil {
		return func(next Operator) Operator { return next }
	}
	return func(next Operator) Operator {
		return opFunc{
			name:    next.Name(),
			version: next.Version(),
			run: func(ctx context.Context, in OpInput) (OpOutput, error) {
				if b.Left() <= 0 {
					return OpOutput{}, ErrBudgetExceeded
				}
				out, err := next.Run(ctx, in)
				cost := costFn(out)
				if cost > 0 && !b.charge(cost) {
					return OpOutput{}, ErrBudgetExceeded
				}
				return out, err
			},
		}
	}
}

// ===== Cache =====

type CacheStore interface {
	Get(key string) (OpOutput, bool)
	Set(key string, val OpOutput)
}

type memoryCache struct{ m sync.Map }

func NewMemoryCache() CacheStore { return &memoryCache{} }
func (c *memoryCache) Get(k string) (OpOutput, bool) {
	v, ok := c.m.Load(k)
	if !ok {
		return OpOutput{}, false
	}
	return v.(OpOutput), true
}
func (c *memoryCache) Set(k string, v OpOutput) { c.m.Store(k, v) }

// WithCache curto-circuito com chave determinística.
func WithCache(store CacheStore, keyFn func(Operator, OpInput) string) Middleware {
	if store == nil {
		return func(next Operator) Operator { return next }
	}
	if keyFn == nil {
		keyFn = MakeCacheKey
	}
	return func(next Operator) Operator {
		return opFunc{
			name:    next.Name(),
			version: next.Version(),
			run: func(ctx context.Context, in OpInput) (OpOutput, error) {
				key := keyFn(next, in)
				if out, ok := store.Get(key); ok {
					// Marcar cache_hit
					out.Metrics = append(out.Metrics, Metric{Name: "cache_hit", Value: 1, Unit: "bool"})
					return out, nil
				}
				out, err := next.Run(ctx, in)
				if err == nil {
					store.Set(key, out)
				}
				return out, err
			},
		}
	}
}
