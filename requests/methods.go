package requests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"
)

// A custom HTTP client with a short timeout to avoid unecessary hangs.
var httpClient = http.Client{
	Timeout: 500 * time.Millisecond,
}

// The URLs.
var panelUrl = fmt.Sprintf("%s://%s:%s", os.Getenv("SCHEME_MONITOR"), os.Getenv("HOST_MONITOR"), os.Getenv("PORT_MONITOR"))
var listenerUrl = fmt.Sprintf("%s://%s:%s", os.Getenv("SCHEME_LISTENER"), os.Getenv("HOST_LISTENER"), os.Getenv("PORT_LISTENER"))

// Gets the panel's health.
func PanelHealthcheck() (*PanelHealth, error) {
	resp, err := httpClient.Get(panelUrl)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result PanelHealth
	err = json.Unmarshal(body, &result)
	return &result, err
}

// Requests to kill the panel.
func PanelKill() error {
	req, err := http.NewRequest("DELETE", panelUrl, nil)
	if err != nil {
		return err
	}
	if _, err := httpClient.Do(req); err != nil {
		return err
	}
	return nil
}

// Gets the listener's health.
func ListenerHealthcheck() (*ListenerHealth, error) {
	resp, err := httpClient.Get(listenerUrl)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result ListenerHealth
	err = json.Unmarshal(body, &result)
	return &result, err
}

// Requests to kill the listener.
func ListenerKill() error {
	req, err := http.NewRequest("DELETE", listenerUrl, nil)
	if err != nil {
		return err
	}
	if _, err := httpClient.Do(req); err != nil {
		return err
	}
	return nil
}

// Requests to kill the shard.
func ListenerShard(id int64) (bool, error) {
	url := fmt.Sprintf("%s/shard?id=%d", listenerUrl, id)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return false, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	return res.StatusCode != 400, nil
}

// Requests a pet and loads it in.
func PetOnDemand(target Pet) (string, error) {
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

// Gets the top urban dictionary definition.
func Urban(p string) (*UrbanDefinition, error) {
	url := fmt.Sprintf("http://api.urbandictionary.com/v0/define?term=%s", p)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result UrbanResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	definitions := result.List[:]
	if len(definitions) == 0 {
		return nil, nil
	}
	sort.Slice(definitions, func(a, b int) bool {
		return definitions[a].Upvotes > definitions[b].Upvotes
	})
	return definitions[0], nil
}
