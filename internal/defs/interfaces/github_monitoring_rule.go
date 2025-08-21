package interfaces

type IMonitoringRule interface {
	IRule
	GetCheckInactivity() bool
	SetCheckInactivity(check bool)
	GetInactiveDaysThreshold() int
	SetInactiveDaysThreshold(days int)
	GetMonitorPRs() bool
	SetMonitorPRs(monitor bool)
}
