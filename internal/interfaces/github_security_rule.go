package interfaces

type ISecurityRule interface {
	IRule
	GetRotateSSHKeys() bool
	SetRotateSSHKeys(rotate bool)
	GetRemoveOldKeys() bool
	SetRemoveOldKeys(remove bool)
	GetKeyPattern() string
	SetKeyPattern(pattern string)
}
