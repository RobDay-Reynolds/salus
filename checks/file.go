package checks

type FileCheck struct {
	Name           string
	Path           string
	IfChanged      string
	FailedSocket   FailedSocket
	FailedHost     FailedHost
	TotalMemChecks []MemUsage
	Group          string
	DependsOn      string
}
