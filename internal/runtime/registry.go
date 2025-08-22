// Package runtime define tipos e interfaces para operadores executáveis
package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
)

// ===== Tipos de domínio (públicos do runtime) =====

// RepoRef identifica um repositório alvo.
type RepoRef struct {
	Owner string
	Name  string
	Head  string // opcional (commit SHA)
}

// ClientBundle agrupa clientes externos (GitHub/LLM/etc.).
type ClientBundle struct {
	GitHub any // use uma interface fina no operador (evita acoplamento)
	LLM    any
}

// Metric métrica objetiva e agregável.
type Metric struct {
	Name   string            `json:"name"`
	Value  float64           `json:"value"`
	Unit   string            `json:"unit"`
	Labels map[string]string `json:"labels,omitempty"`
}

// Insight achado/aviso pontual.
type Insight struct {
	Key     string         `json:"key"`
	Summary string         `json:"summary"`
	Details map[string]any `json:"details,omitempty"`
	Score   float64        `json:"score,omitempty"` // 0..1
}

// OpInput entrada comum para qualquer operador.
type OpInput struct {
	Repo    RepoRef
	Params  map[string]any
	Clients ClientBundle
	DryRun  bool
	// IdempotencyKey opcional; se vazio, o Manager gera uma chave determinística.
	IdempotencyKey string
}

// OpOutput saída padrão de operadores.
type OpOutput struct {
	Data      any
	Metrics   []Metric
	Insights  []Insight
	Artifacts map[string][]byte
}

// Operator descreve uma unidade executável plugável.
type Operator interface {
	Name() string
	Version() string
	Run(ctx context.Context, in OpInput) (OpOutput, error)
}

// TypedOperator permite ergonomia com generics sem poluir o runtime.
type TypedOperator[I any, O any] interface {
	Name() string
	Version() string
	RunTyped(ctx context.Context, in I) (O, error)
}

// Adapt converte um TypedOperator em Operator padrão.
func Adapt[I any, O any](t TypedOperator[I, O], decode func(OpInput) (I, error), encode func(O) OpOutput) Operator {
	return opFunc{
		name:    t.Name(),
		version: t.Version(),
		run: func(ctx context.Context, in OpInput) (OpOutput, error) {
			typedIn, err := decode(in)
			if err != nil {
				return OpOutput{}, err
			}
			out, err := t.RunTyped(ctx, typedIn)
			if err != nil {
				return OpOutput{}, err
			}
			return encode(out), nil
		},
	}
}

// opFunc é o invólucro base (também usado pelos middlewares).
type opFunc struct {
	name    string
	version string
	run     func(context.Context, OpInput) (OpOutput, error)
}

func (o opFunc) Name() string                                          { return o.name }
func (o opFunc) Version() string                                       { return o.version }
func (o opFunc) Run(ctx context.Context, in OpInput) (OpOutput, error) { return o.run(ctx, in) }

// ===== Registry =====

type Descriptor struct {
	Name    string
	Version string
}

type Registry interface {
	Register(op Operator)
	Get(name string) (Operator, bool)
	List() []Descriptor
}

type registry struct {
	mu  sync.RWMutex
	ops map[string]Operator
}

func NewRegistry() Registry { return &registry{ops: make(map[string]Operator)} }

func (r *registry) Register(op Operator) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ops[op.Name()] = op
}

func (r *registry) Get(name string) (Operator, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	op, ok := r.ops[name]
	return op, ok
}

func (r *registry) List() []Descriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d := make([]Descriptor, 0, len(r.ops))
	for _, op := range r.ops {
		d = append(d, Descriptor{Name: op.Name(), Version: op.Version()})
	}
	return d
}

// ===== Erros sentinela úteis =====

var (
	ErrBudgetExceeded = errors.New("budget exceeded")
	ErrRateLimited    = errors.New("rate limited")
	ErrCanceled       = context.Canceled
	ErrTimeout        = context.DeadlineExceeded
)

// MarshalParamsDeterministic serializa Params de forma estável para logs/chaves.
func MarshalParamsDeterministic(m map[string]any) []byte {
	if m == nil {
		return []byte("null")
	}
	// Encode estável (uma abordagem simples: ordenar por chave e re-montar)
	type kv struct {
		K string
		V any
	}
	pairs := make([]kv, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, kv{k, v})
	}
	// sort
	sortFn := func(i, j int) bool { return pairs[i].K < pairs[j].K }
	tmp := pairs
	for i := 0; i < len(tmp); i++ {
		for j := i + 1; j < len(tmp); j++ {
			if !sortFn(i, j) {
				tmp[i], tmp[j] = tmp[j], tmp[i]
			}
		}
	}
	ordered := make(map[string]any, len(pairs))
	for _, p := range tmp {
		ordered[p.K] = p.V
	}
	b, _ := json.Marshal(ordered)
	return b
}
