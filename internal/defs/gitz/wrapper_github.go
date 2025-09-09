// Package gitz provides a set of types and functions for working with GitHub.
package gitz

import "github.com/kubex-ecosystem/ghbex/internal/defs/interfaces"

type GitHub struct {
	*GitHubAuth `yaml:"auth" json:"auth"`
	Repos       []*RepoCfg `yaml:"repos" json:"repos"`
}

func NewGitHubType(auth *GitHubAuth, repos []interfaces.IRepoCfg) *GitHub {
	rps := make([]*RepoCfg, len(repos))
	for i, repo := range repos {
		rps[i] = repo.(*RepoCfg)
	}

	return &GitHub{
		GitHubAuth: auth,
		Repos:      rps,
	}
}

func NewGitHub(auth interfaces.IGitHubAuth, repos []interfaces.IRepoCfg) interfaces.IGitHub {
	if auth == nil {
		return NewGitHubType(nil, repos)
	}
	return NewGitHubType(auth.(*GitHubAuth), repos)
}

func (g *GitHub) GetAuth() interfaces.IGitHubAuth {
	return g.GitHubAuth
}

func (g *GitHub) SetAuth(auth interfaces.IGitHubAuth) {
	if auth == nil {
		g.GitHubAuth = nil
	} else {
		g.GitHubAuth = auth.(*GitHubAuth)
	}
}

func (g *GitHub) GetRepos() []interfaces.IRepoCfg {
	if g.Repos == nil {
		return nil
	}
	var result []interfaces.IRepoCfg
	for _, repo := range g.Repos {
		result = append(result, repo)
	}
	return result
}

func (g *GitHub) SetRepos(repos []interfaces.IRepoCfg) {
	rps := make([]*RepoCfg, len(repos))
	for i, repo := range repos {
		rps[i] = repo.(*RepoCfg)
	}
	g.Repos = rps
}
