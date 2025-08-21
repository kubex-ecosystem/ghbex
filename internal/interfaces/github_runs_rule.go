package interfaces

type IRunsRule interface {
	IRule
	GetMaxAgeDays() int
	SetMaxAgeDays(days int)
	GetKeepSuccessLast() int
	SetKeepSuccessLast(days int)
	GetOnlyWorkflows() []string
	SetOnlyWorkflows(workflows []string)
}
