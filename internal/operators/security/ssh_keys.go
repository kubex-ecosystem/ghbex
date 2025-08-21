// Package security provides functions to manage GitHub repository security.
package security

import (
	"context"
	"fmt"

	"github.com/google/go-github/v61/github"
)

// RotateSSHKeys rotates SSH deploy keys for a repository
func RotateSSHKeys(ctx context.Context, cli *github.Client, owner, repo string, dry bool) (*SSHKeyPair, error) {
	if dry {
		// In dry-run mode, just generate a key pair without applying changes
		return generateSSHKeyPair()
	}

	// Generate new SSH key pair
	newKeyPair, err := generateSSHKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate SSH key pair: %w", err)
	}

	// Add new deploy key to repository
	keyTitle := fmt.Sprintf("ghbex-auto-generated-%d", newKeyPair.KeyID)
	deployKey := &github.Key{
		Title:    &keyTitle,
		Key:      &newKeyPair.PublicKey,
		ReadOnly: github.Bool(false), // Allow write access
	}

	createdKey, _, err := cli.Repositories.CreateKey(ctx, owner, repo, deployKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create deploy key: %w", err)
	}

	newKeyPair.KeyID = createdKey.GetID()

	// TODO: Remove old auto-generated keys (implement in next iteration)
	// This would require tracking which keys were auto-generated

	return newKeyPair, nil
}

// ListDeployKeys returns all deploy keys for a repository
func ListDeployKeys(ctx context.Context, cli *github.Client, owner, repo string) ([]*github.Key, error) {
	opt := &github.ListOptions{PerPage: 100}
	var allKeys []*github.Key

	for {
		keys, resp, err := cli.Repositories.ListKeys(ctx, owner, repo, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list deploy keys: %w", err)
		}

		allKeys = append(allKeys, keys...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allKeys, nil
}

// RemoveOldDeployKeys removes deploy keys that match a specific pattern (auto-generated ones)
func RemoveOldDeployKeys(ctx context.Context, cli *github.Client, owner, repo string, pattern string, dry bool) (removed int, err error) {
	keys, err := ListDeployKeys(ctx, cli, owner, repo)
	if err != nil {
		return 0, err
	}

	for _, key := range keys {
		if key.GetTitle() != "" && contains(key.GetTitle(), pattern) {
			if dry {
				removed++
				continue
			}

			_, err := cli.Repositories.DeleteKey(ctx, owner, repo, key.GetID())
			if err != nil {
				return removed, fmt.Errorf("failed to delete key %d: %w", key.GetID(), err)
			}
			removed++
		}
	}

	return removed, nil
}
