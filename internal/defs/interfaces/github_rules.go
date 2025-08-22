package interfaces

type IRules interface {
	GetRuns() IRule
	GetArtifacts() IRule
	GetReleases() IRule
	GetSecurity() IRule
	GetMonitoring() IRule

	GetRunsRule() IRunsRule
	GetArtifactsRule() IArtifactsRule
	GetReleasesRule() IReleasesRule
	GetSecurityRule() ISecurityRule
	GetMonitoringRule() IMonitoringRule

	SetRunsRule(runs IRunsRule)
	SetArtifactsRule(artifacts IArtifactsRule)
	SetReleasesRule(releases IReleasesRule)
	SetSecurityRule(security ISecurityRule)
	SetMonitoringRule(monitoring IMonitoringRule)
}
