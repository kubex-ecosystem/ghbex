package interfaces

type IArtifactsRule interface {
	IRule
	GetMaxAgeDays() int
	SetMaxAgeDays(days int)
}
