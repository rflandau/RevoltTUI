// Singleton tracking global information about the current tty, particularly dimensions.
package terminal

import "sync"

var width, height int
var dimensionLock sync.Mutex

func SetDimensions(w, h int) {
	dimensionLock.Lock()
	width = w
	height = h
	dimensionLock.Unlock()
}

func Width() int {
	dimensionLock.Lock()
	defer dimensionLock.Unlock()
	return width
}

func Height() int {
	dimensionLock.Lock()
	defer dimensionLock.Unlock()
	return height
}
