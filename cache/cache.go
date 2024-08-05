/*
The Cache package provides a singleton containing the most recently EventReady data and manages
locking and updating it.
*/
package cache

import (
	"revolt_tui/log"
	"sync"

	"github.com/sentinelb51/revoltgo"
)

var cache *revoltgo.EventReady
var ready bool // cache has been populated at least once
var cacheMTX *sync.RWMutex

// updates the cache on Ready event; registered during the AddHandler
func OnEventReadyfunc(_ *revoltgo.Session, r *revoltgo.EventReady) {
	cacheMTX.Lock()

	isCacheNil := cache == nil

	go log.Writer.Debug("Session is ready", "nil cache?", isCacheNil, "EventReady", r)
	cache = r
	ready = true
	cacheMTX.Unlock()
}

//#region public getters

func Ready() bool {
	return ready
}

func Servers() []*revoltgo.Server {
	if cache == nil {
		return nil
	}
	cacheMTX.RLock()
	defer cacheMTX.RUnlock()
	// TODO may need to duplicate cache.Servers to prevent data destruction by callers
	return cache.Servers
}

//#endregion public getters
