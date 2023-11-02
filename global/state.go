package global

import "sync"

var (
	HostTargetMap = &sync.Map{}
	CacheKeyPath  = &sync.Map{}
)
