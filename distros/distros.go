package distros

type Distro interface {
    GetSecurityUpdates() int
    GetRebootRequired() bool
}
