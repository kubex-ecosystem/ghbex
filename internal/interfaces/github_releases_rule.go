package interfaces

type IReleasesRule interface {
	IRule
	GetDeleteDrafts() bool
	SetDeleteDrafts(deleteDrafts bool)
}
