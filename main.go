package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"revolt_tui/broker"
	"revolt_tui/cfgdir"
	"revolt_tui/controller"
	"revolt_tui/credentials"
	"revolt_tui/log"
	"revolt_tui/modes"
	"revolt_tui/modes/server"
	serverselection "revolt_tui/modes/serverSelection"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sentinelb51/revoltgo"
	"github.com/spf13/pflag"
)

const (
	tokenFileName string = "token"
)

// main's init just defines flags
func init() {
	pflag.String("loglevel", "DEBUG",
		"set the log level.\n"+
			"Viable options (from most verbose to least) are: debug, info, warn, error, fatal")
}

func main() {
	// consume flags
	pflag.Parse()

	// set up the logger singleton
	if err := log.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	// attempt to login via token, fallback to credentials on failure
	var session *revoltgo.Session = loginViaToken()
	if session == nil {
		var killed bool
		session, killed = loginViaCredentials()
		if killed {
			fmt.Fprintln(os.Stdout, "You must authenticate to use RevoltTUI")
			log.Destroy()
			return
		}
		if session == nil { // die on failure
			fmt.Fprintln(os.Stderr, "An error has occurred. Sorry, friend.")
			log.Destroy()
			return
		}
		// write the token from the session so we do not need to prompt next time
		tknpth := path.Join(cfgdir.Get(), tokenFileName)
		log.Writer.Debugf("creating token at path '%v'", tknpth)
		tokenFile, err := os.OpenFile(tknpth, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			log.Writer.Warn("failed to create token file", "error", err)
		} else {
			tokenFile.WriteString(session.Token)
		}
		if tokenFile != nil {
			tokenFile.Close()
		}
	}

	// register modes
	modes.Add(modes.ServerSelection, &serverselection.Action{})
	modes.Add(modes.Server, server.New())

	/*func(session *revoltgo.Session, r *revoltgo.EventReady) {
		log.Writer.Info("Ready to handle commands from %v user(s) across %d servers from %d channels",
			len(r.Users), len(r.Servers), len(r.Channels))
		cache.OnEventReadyFunc(session, r)
		// update all dependencies of model.cache

		//#region server selection
		// cast to items
		var items []list.Item = make([]list.Item, len(model.cache.Servers))
		for i, server := range model.cache.Servers {
			items[i] = serverItem{
				title:       server.Name,
				id:          server.ID,
				description: server.Description,
			}
		}

		if !serverSelectionImpl.initialized { // first update
			if term.width == 0 || term.height == 0 { // cannot initialize until the first WindowSizeMsg
				serverSelectionImpl.list = list.New(items, list.NewDefaultDelegate(), term.width, term.height)
			} else {
				model.log.Debug("server list ready, but terminal dimensions have not been recieved")
			}
		} else { // later updates
			// TODO this returns a tea.Cmd for filtering
			serverSelectionImpl.list.SetItems(items)
		}
		serverSelectionImpl.initialized = true
		//#endregion
	}*/

	// spin up program
	p := tea.NewProgram(controller.Initial(session))

	// attach ready handler to our revoltgo session so we can inject messages into bubble tea
	session.AddHandler(func(session *revoltgo.Session, r *revoltgo.EventReady) {
		// update the cache
		broker.OnEventReadyFunc(session, r)
		log.Writer.Info("cache updated")
		p.Send(broker.CacheUpdatedMsg{})
	})

	_, err := p.Run()
	if err != nil {
		log.Writer.Error("error running the main model", "error", err)
	}

	// on completion, clean up resources
	session.Close()
	log.Destroy()
}

// Attempts to authenticate via the existing token, found in the config directory.
// Automatically logs to the given logger.
// Returns an authenticated session or nil.
func loginViaToken() (session *revoltgo.Session) {
	tknPth := path.Join(cfgdir.Get(), tokenFileName)
	f, err := os.Open(tknPth)
	if err != nil {
		log.Writer.Warn("failed to open token file. Skipping token login.",
			"error", err, "file path", tknPth)
		return nil
	}

	token, err := io.ReadAll(f)
	if err != nil {
		log.Writer.Warn("failed to read token. Skipping token login.",
			"error", err, "file path", tknPth)
		return nil
	}
	return revoltgo.New(string(token))
}

// Attempts to authenticate via email and password.
// Automatically logs to the given logger.
// Returns an authenticated session.
// Main is expected to exit if nil is returned; the error will already have been logged.
func loginViaCredentials() (session *revoltgo.Session, killed bool) {
	// spawn the login dialog
	credProg := tea.NewProgram(credentials.InitialModel(), tea.WithAltScreen())
	finalCredModelRaw, err := credProg.Run()
	if err != nil {
		log.Writer.Error(err)
		return nil, false
	}

	// cast the raw model to its actual model
	credModel, ok := finalCredModelRaw.(credentials.Model)
	if !ok {
		log.Writer.Error("failed to cast final credential model")
		return nil, false
	}

	// if the program was killed, do not try to authenticate
	if credModel.Killed {
		return nil, true
	}

	return credModel.Session, false
}
