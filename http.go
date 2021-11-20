package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type panelHealth struct {
	Connections int `json:"connected"`
}

var httpClient = http.Client{
	Timeout: 1 * time.Second,
}

// The endpoint that Discord calls when there is an interaction.
func httpInteract(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	slashHandler(w, req)
}

// The endpoint that registers a new custom command for a guild.
func httpRegister(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if c := handleAuthorization(w, req); !c {
		return
	}
}

// The endpoint that removes a custom command from a guild.
func httpUnregister(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if c := handleAuthorization(w, req); !c {
		return
	}
}

// The endpoint that invalidates the cache for a guild.
func httpInvalidate(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if c := handleAuthorization(w, req); !c {
		return
	}
	// TODO: Invalidate cache.
}

// Checks if the request actually came from the web panel, using authorization tokens.
func handleAuthorization(w http.ResponseWriter, req *http.Request) bool {
	auth := req.Header.Get("Authorization")
	if auth == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	if auth != token {
		w.WriteHeader(http.StatusForbidden)
		return false
	}
	return true
}

// Gets the panel's health.
func requestPanelHealth() (*panelHealth, error) {
	url := fmt.Sprintf("%s://%s:%s", os.Getenv("SCHEME_MONITOR"), os.Getenv("HOST_MONITOR"), os.Getenv("PORT_MONITOR"))
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result panelHealth
	err = json.Unmarshal(body, &result)
	return &result, err
}

// Requests to kill the panel.
func requestPanelKill() error {
	url := fmt.Sprintf("%s://%s:%s", os.Getenv("SCHEME_MONITOR"), os.Getenv("HOST_MONITOR"), os.Getenv("PORT_MONITOR"))
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	if _, err := httpClient.Do(req); err != nil {
		return err
	}
	return nil
}
