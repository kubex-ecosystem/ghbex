package runtime

import (
	"context"
	"fmt"
	"sync"
)

// Manager despacha operadores com middlewares, idempotência e cancelamento por job.
type Manager struct {
	reg Registry
	mws []Middleware

	muJobs sync.Mutex
	jobs   map[string]context.CancelFunc // key = IdempotencyKey
}

func NewManager(reg Registry, mws ...Middleware) *Manager {
	return &Manager{reg: reg, mws: mws, jobs: make(map[string]context.CancelFunc)}
}

// Dispatch executa um operador pelo nome, aplicando middlewares e idempotência.
func (m *Manager) Dispatch(ctx context.Context, name string, in OpInput) (OpOutput, error) {
	op, ok := m.reg.Get(name)
	if !ok {
		return OpOutput{}, fmt.Errorf("unknown operator: %s", name)
	}
	op = Chain(op, m.mws...)
	if in.IdempotencyKey == "" {
		in.IdempotencyKey = MakeIDKey(op, in)
	}

	ctx, cancel := context.WithCancel(ctx)
	m.track(in.IdempotencyKey, cancel)
	defer m.untrack(in.IdempotencyKey)

	return op.Run(ctx, in)
}

// OperatorStatus representa o status de um operador em execução.
type OperatorStatus struct {
	Status   string // pending|running|done|failed
	Error    error
	Metadata map[string]any
}

// Monitor executa e emite status em canal (forma simples).
func (m *Manager) Monitor(ctx context.Context, name string, in OpInput) (<-chan *OperatorStatus, error) {
	ch := make(chan *OperatorStatus, 8)
	go func() {
		defer close(ch)
		ch <- &OperatorStatus{Status: "running"}
		out, err := m.Dispatch(ctx, name, in)
		if err != nil {
			ch <- &OperatorStatus{Status: "failed", Error: err}
			return
		}
		ch <- &OperatorStatus{Status: "done", Metadata: map[string]any{"metrics": out.Metrics}}
	}()
	return ch, nil
}

func (m *Manager) Cancel(idem string) {
	m.muJobs.Lock()
	defer m.muJobs.Unlock()
	if c, ok := m.jobs[idem]; ok {
		c()
		delete(m.jobs, idem)
	}
}

func (m *Manager) track(idem string, c context.CancelFunc) {
	m.muJobs.Lock()
	m.jobs[idem] = c
	m.muJobs.Unlock()
}
func (m *Manager) untrack(idem string) { m.muJobs.Lock(); delete(m.jobs, idem); m.muJobs.Unlock() }

// Defaults para uso imediato.
var (
	DefaultRegistry = NewRegistry()
	DefaultManager  = NewManager(DefaultRegistry)
)
