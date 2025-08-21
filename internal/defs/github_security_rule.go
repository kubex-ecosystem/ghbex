package defs

import "github.com/rafa-mori/ghbex/internal/interfaces"

type SecurityRule struct {
	RotateSSHKeys bool   `yaml:"rotate_ssh_keys" json:"rotate_ssh_keys"`
	RemoveOldKeys bool   `yaml:"remove_old_keys" json:"remove_old_keys"`
	KeyPattern    string `yaml:"key_pattern" json:"key_pattern"`
}

func NewSecurityRuleType(rotateSSHKeys, removeOldKeys bool, keyPattern string) *SecurityRule {
	return &SecurityRule{
		RotateSSHKeys: rotateSSHKeys,
		RemoveOldKeys: removeOldKeys,
		KeyPattern:    keyPattern,
	}
}

func NewSecurityRule(rotateSSHKeys, removeOldKeys bool, keyPattern string) interfaces.ISecurityRule {
	return NewSecurityRuleType(rotateSSHKeys, removeOldKeys, keyPattern)
}

func (r *SecurityRule) GetRotateSSHKeys() bool {
	return r.RotateSSHKeys
}

func (r *SecurityRule) SetRotateSSHKeys(rotate bool) {
	r.RotateSSHKeys = rotate
}

func (r *SecurityRule) GetRemoveOldKeys() bool {
	return r.RemoveOldKeys
}

func (r *SecurityRule) SetRemoveOldKeys(remove bool) {
	r.RemoveOldKeys = remove
}

func (r *SecurityRule) GetKeyPattern() string {
	return r.KeyPattern
}

func (r *SecurityRule) SetKeyPattern(pattern string) {
	r.KeyPattern = pattern
}

func (r *SecurityRule) GetRuleName() string {
	return "security"
}

func (r *SecurityRule) SetRuleName(name string) {
	// No-op for security rule
}
