// Singleton instance that manages the current state of the TUI, from the known dimensions of the
// tty to its current pwd.
// Primarily a mechanism for passing parameters and ensuring data consistency between modes
// (operating as the master copy).
// While controller is the primary operator, this is the primary data broker.
package broker

import (
	"revolt_tui/log"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sentinelb51/revoltgo"
)

//#region tty dimensions

var width, height int
var dimensionLock sync.Mutex

// Set the dimensions of the tty.
func SetDimensions(w, h int) {
	dimensionLock.Lock()
	width = w
	height = h
	dimensionLock.Unlock()
}

// Returns the current terminal width
func Width() int {
	dimensionLock.Lock()
	defer dimensionLock.Unlock()
	return width
}

// Returns the current terminal height
func Height() int {
	dimensionLock.Lock()
	defer dimensionLock.Unlock()
	return height
}

//#endregion tty dimensions

//#region current server

var curServer *revoltgo.Server
var serverLock sync.Mutex

// Updates the current server the user is interacting with.
func SetCurrentServer(svr *revoltgo.Server) {
	serverLock.Lock()
	curServer = svr
	serverLock.Unlock()
}

// Returns the current server the user is interacting with.
// If the user is not currently interacting with a server, this may be nil.
func GetCurrentServer() *revoltgo.Server {
	serverLock.Lock()
	defer serverLock.Unlock()
	return curServer
}

//#endregion current server

//#region cache

var cache *revoltgo.EventReady
var ready bool // cache has been populated at least once
var cacheMTX sync.RWMutex

// updates the cache on Ready event; registered during the AddHandler
func OnEventReadyFunc(_ *revoltgo.Session, r *revoltgo.EventReady) {
	cacheMTX.Lock()

	isCacheNil := cache == nil

	go log.Writer.Debug("Session is ready", "nil cache?", isCacheNil, "EventReady", r)
	cache = r
	ready = true
	cacheMTX.Unlock()
}

func CacheReady() bool {
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

type CacheUpdatedMsg struct {
	tea.Msg
}

//#endregion cache
