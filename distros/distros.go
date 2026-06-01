package distros

import "time"

type Distro interface {
	GetSecurityUpdates() int
	GetTotalUpdates() int
	GetRebootRequired() bool
	GetLatestCacheChange() time.Time
}
