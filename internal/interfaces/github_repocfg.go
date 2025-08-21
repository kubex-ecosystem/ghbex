package interfaces

type IRepoCfg interface {
	GetOwner() string
	GetName() string
	GetRules() IRules
	GetMonitoring() IMonitoringRule
	SetOwner(owner string)
	SetName(name string)
	SetRules(rules IRules)
	SetMonitoring(monitoring IMonitoringRule)
}
