// Package analytics fornece operadores para análise de repositórios.
package analytics

import (
	"context"
	"time"

	rt "github.com/kubex-ecosystem/gemx/ghbex/internal/runtime"
)

// ===== Interface fina para GitHub (evita acoplamento ao client real) =====

type Issue struct {
	State   string // "open" | "closed"
	Updated time.Time
}

type GitHubAPI interface {
	ListIssues(ctx context.Context, owner, repo string, page, perPage int, state string, since time.Time) (items []Issue, nextPage int, err error)
}

// ===== I/O fortemente tipado para o operador de exemplo =====

type AnalyzeRepositoryInput struct {
	Repo         rt.RepoRef
	AnalysisDays int
	Clients      struct{ GitHub GitHubAPI }
}

type ScoreInsightsReport struct {
	AnalysisDays int
	IssuesTotal  int
	OpenIssues   int
	ClosedIssues int
}

// ===== Implementação TypedOperator com paginação real =====

type AnalyzeRepositoryOperator struct{}

func (o AnalyzeRepositoryOperator) Name() string    { return "analytics.analyze_repository" }
func (o AnalyzeRepositoryOperator) Version() string { return "1.0.0" }

func (o AnalyzeRepositoryOperator) RunTyped(ctx context.Context, in *AnalyzeRepositoryInput) (*ScoreInsightsReport, error) {
	if in.AnalysisDays <= 0 {
		in.AnalysisDays = 30
	}
	since := time.Now().Add(-time.Duration(in.AnalysisDays) * 24 * time.Hour)

	page := 1
	per := 100
	var total, open, closed int
	for {
		items, next, err := in.Clients.GitHub.ListIssues(ctx, in.Repo.Owner, in.Repo.Name, page, per, "all", since)
		if err != nil {
			return nil, err
		}
		for _, it := range items {
			total++
			if it.State == "open" {
				open++
			} else {
				closed++
			}
		}
		if next == 0 {
			break
		}
		page = next
	}

	rep := &ScoreInsightsReport{
		AnalysisDays: in.AnalysisDays,
		IssuesTotal:  total,
		OpenIssues:   open,
		ClosedIssues: closed,
	}
	return rep, nil
}

// RegisterExample registra o operador adaptado no Registry do runtime.
func RegisterExample(reg rt.Registry) {
	op := rt.Adapt(
		AnalyzeRepositoryOperator{},
		func(in rt.OpInput) (*AnalyzeRepositoryInput, error) {
			// Decode: extrai GitHubAPI do bundle e Params
			gh, _ := in.Clients.GitHub.(GitHubAPI)
			ai := &AnalyzeRepositoryInput{Repo: in.Repo, AnalysisDays: 0}
			if v, ok := in.Params["analysis_days"].(int); ok {
				ai.AnalysisDays = v
			}
			ai.Clients.GitHub = gh
			return ai, nil
		},
		func(out *ScoreInsightsReport) rt.OpOutput {
			return rt.OpOutput{
				Data: out,
				Metrics: []rt.Metric{
					{Name: "issues_total", Value: float64(out.IssuesTotal), Unit: "count"},
					{Name: "issues_open", Value: float64(out.OpenIssues), Unit: "count"},
					{Name: "issues_closed", Value: float64(out.ClosedIssues), Unit: "count"},
				},
			}
		},
	)
	reg.Register(op)
}
