package cache

import "sync"

var (
	TypeAndNumberToLineNameCache     = make(map[string]string)
	TypeAndNumberToLineNameCacheLock = sync.Mutex{}
	StopIdForDirectionCache          = make(map[string]string)
	StopIdForDirectionCacheLock      = sync.Mutex{}
)
