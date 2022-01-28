package main

import (
	"net/http"

	"github.com/arraybot/nautilus/commands"
	"github.com/arraybot/nautilus/config"
)

// The endpoint that Discord calls when there is an interaction.
func epInteract(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	slashHandler(w, req)
}

// The endpoint that registers a new custom command for a guild.
func epRegister(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if c := handleAuthorization(w, req); !c {
		return
	}
	// TODO: Register custom command.
}

// The endpoint that removes a custom command from a guild.
func epUnregister(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if c := handleAuthorization(w, req); !c {
		return
	}
	// TODO: Unregister custom command.
}

// The endpoint that invalidates the cache for a guild.
func epInvalidate(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if c := handleAuthorization(w, req); !c {
		return
	}
	params := req.URL.Query()
	guild := params.Get("guild")
	if guild == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	commands.CacheInvisibility.Invalidate(guild)
}

// Checks if the request actually came from the web panel, using authorization tokens.
func handleAuthorization(w http.ResponseWriter, req *http.Request) bool {
	auth := req.Header.Get("Authorization")
	if auth == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	if auth != config.CommandsSecret {
		w.WriteHeader(http.StatusForbidden)
		return false
	}
	return true
}
