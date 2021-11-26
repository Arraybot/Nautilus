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
	Timeout: 1 * time.Second,
}

// Gets the panel's health.
func PanelHealthcheck() (*PanelHealth, error) {
	url := fmt.Sprintf("%s://%s:%s", os.Getenv("SCHEME_MONITOR"), os.Getenv("HOST_MONITOR"), os.Getenv("PORT_MONITOR"))
	resp, err := httpClient.Get(url)
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
