/*
Logger singleton. Initialized during startup (by main).
*/
package log

import (
	"fmt"
	"os"
	"path"
	"revolt_tui/cfgdir"

	"github.com/charmbracelet/log"
	"github.com/spf13/pflag"
)

const (
	filePermission        = 0600
	dirPermission         = 0700
	defaultLogName string = "log.txt" // in config directory
)

var Writer *log.Logger
var logPath string   // set on Initialize()
var logFile *os.File // set on Initialize()

// Initializes the writer singleton
func Initialize() error {
	loglevel, err := pflag.CommandLine.GetString("loglevel")
	if err != nil {
		panic(err) // developer error
	}
	lvl, err := log.ParseLevel(loglevel)
	if err != nil {
		return fmt.Errorf("unknown log level '%s'; see -h for help", loglevel)
	}

	logPath = path.Join(cfgdir.Get(), defaultLogName)

	// open log file
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filePermission)
	if err != nil {
		// try again in pwd
		logPath = defaultLogName
		logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filePermission)
		if err != nil {
			return fmt.Errorf("failed to open file in config directory or in pwd: %v", err)
		}
	}

	// spin up a logger
	logOpts := log.Options{Level: lvl, ReportTimestamp: true}
	if lvl == log.DebugLevel {
		logOpts.ReportCaller = true
	}

	Writer = log.NewWithOptions(logFile, logOpts)

	return nil
}

// Destroys the writer singleton.
// Should only be called on program exit.
func Destroy() {
	Writer = nil
	logFile.Close()
}
