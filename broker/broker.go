// Singleton instance that manages the current state of the TUI, from the known dimensions of the
// tty to its current pwd.
// Primarily a mechanism for passing parameters and ensuring data consistency between modes
// (operating as the master copy).
// While controller is the primary operator, this is the primary data broker.
package broker

import (
	"sync"

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

var server *revoltgo.Server
var serverLock sync.Mutex

// Updates the current server the user is interacting with.
func SetServer(svr *revoltgo.Server) {
	serverLock.Lock()
	server = svr
	serverLock.Unlock()
}

// Returns the current server the user is interacting with.
// If the user is not currently interacting with a server, this may be nil.
func GetServer() *revoltgo.Server {
	serverLock.Lock()
	defer serverLock.Unlock()
	return server
}

//#endregion current server
