package distros

type Distro interface {
	GetSecurityUpdates() int
	GetTotalUpdates() int
	GetRebootRequired() bool
}
