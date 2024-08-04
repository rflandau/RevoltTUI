package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"revolt_tui/src/credentials"
	"revolt_tui/src/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/sentinelb51/revoltgo"
	"github.com/spf13/pflag"
)

const (
	defaultConfigSubDirName string = "revolttui"
	defaultLogName          string = "log.txt" // in config directory
	defaultTokenName        string = "token"
	filePermission                 = 0644
)

func init() {
	pflag.String("log", "", "set the log path. Defaults to '"+defaultLogName+"' in your OS's default config directory.")
	pflag.String("loglevel", "DEBUG", "set the log level.")
}

func main() {
	// consume flags
	pflag.Parse()

	logpath, err := pflag.CommandLine.GetString("log")
	if err != nil {
		panic(err)
	}
	if logpath == "" {
		cfgDir, err := os.UserConfigDir()
		if err != nil {
			panic(err)
		}
		logpath = path.Join(cfgDir, defaultConfigSubDirName, defaultLogName)
	}

	loglevel, err := pflag.CommandLine.GetString("loglevel")
	if err != nil {
		panic(err)
	}
	lvl, err := log.ParseLevel(loglevel)
	if err != nil {
		log.Warn(err)
		lvl = log.InfoLevel
	}

	// open log file
	logFile, err := os.OpenFile(logpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, filePermission)
	if err != nil {
		fmt.Printf("Failed to open file to log to @ %v: %v", logpath, err)
		return
	}

	// spin up a logger
	logOpts := log.Options{Level: lvl, ReportTimestamp: true}
	if lvl == log.DebugLevel {
		logOpts.ReportCaller = true
	}

	l := log.NewWithOptions(logFile, logOpts)
	l.SetLevel(lvl)

	// attempt to login via token, then credentials.
	var session *revoltgo.Session = loginViaToken(l)
	if session == nil {
		session = loginViaCredentials(l)
	}
	// die on failure
	if session == nil {
		fmt.Fprintln(os.Stdout, "An error has occurred. Sorry, friend.")
		return
	}
	// spin up program
	tea.NewProgram(model.Initial(session, l))
}

// Attempts to authenticate via the existing token, found in the config directory.
// Automatically logs to the given logger.
// Returns an authenticated session or nil.
func loginViaToken(l *log.Logger) (session *revoltgo.Session) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		l.Warn("failed to open user config directory. Skipping token login.",
			"error", err, "config path", cfgDir)
		return nil
	}
	tknPth := path.Join(cfgDir, defaultConfigSubDirName, defaultTokenName)
	f, err := os.Open(tknPth)
	if err != nil {
		l.Warn("failed to open token file. Skipping token login.",
			"error", err, "file path", tknPth)
		return nil
	}

	token, err := io.ReadAll(f)
	if err != nil {
		l.Warn("failed to read token. Skipping token login.",
			"error", err, "file path", tknPth)
		return nil
	}
	return revoltgo.New(string(token))
}

// Attempts to authenticate via email and password.
// Automatically logs to the given logger.
// Returns an authenticated session. Main is expected to exit if nil is returned; the error will
// already have been logged.
func loginViaCredentials(l *log.Logger) (session *revoltgo.Session) {
	credProg := tea.NewProgram(credentials.InitialModel(), tea.WithAltScreen())
	finalCredModelRaw, err := credProg.Run()
	if err != nil {
		l.Error(err)
		return nil
	}

	// cast the raw model to its actual model
	credModel, ok := finalCredModelRaw.(credentials.Model)
	if !ok {
		l.Error("failed to cast final credential model")
		return nil
	}
	email := credModel.EmailTI.Value()
	pass := credModel.PassTI.Value()

	sess, lr, err := revoltgo.NewWithLogin(revoltgo.LoginData{
		Email:        email,
		Password:     pass,
		FriendlyName: "TUIFriendly",
	})
	if err != nil {
		l.Error(err)
		return nil
	}
	l.Debug("completed login attempt", "loginResponse", lr)

	return sess
}
