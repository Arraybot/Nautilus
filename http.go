package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// A wrapper for the panel health.
type panelHealth struct {
	Connections int `json:"connected"`
}

// An interface for pet links.
type pet interface {
	source() string
	link() string
}

// A cat.
type cat struct {
	File string `json:"file"`
}

// A dog.
type dog struct {
	Url string `json:"url"`
}

// A custom HTTP client with a short timeout to avoid unecessary hangs.
var httpClient = http.Client{
	Timeout: 1 * time.Second,
}

// Implement pet for cat.
func (p *cat) link() string {
	return p.File
}

func (p *cat) source() string {
	return "https://aws.random.cat/meow"
}

// Implement pet for dog.
func (p *dog) link() string {
	return p.Url
}

func (p *dog) source() string {
	return "https://random.dog/woof.json"
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

// Requests a pet and loads it in.
func requestPet(target pet) (string, error) {
	resp, err := httpClient.Get(target.source())
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, target)
	if err != nil {
		return "", err
	}
	return target.link(), nil
}
