package interfaces

type IRuntime interface {
	GetDebug() bool
	GetDryRun() bool
	GetReportDir() string
	SetDebug(debug bool)
	SetDryRun(dryRun bool)
	SetReportDir(reportDir string)
}
