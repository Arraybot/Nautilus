package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
)

// A custom HTTP client with a short timeout to avoid unecessary hangs.
var httpClientShort = http.Client{
	Timeout: 500 * time.Millisecond,
}
var httpClientLong = http.Client{
	Timeout: 2 * time.Second,
}

// The URLs.
var panelUrl = fmt.Sprintf("%s://%s:%s", os.Getenv("SCHEME_MONITOR"), os.Getenv("HOST_MONITOR"), os.Getenv("PORT_MONITOR"))
var listenerUrl = fmt.Sprintf("%s://%s:%s", os.Getenv("SCHEME_LISTENER"), os.Getenv("HOST_LISTENER"), os.Getenv("PORT_LISTENER"))

// Gets the panel's health.
func PanelHealthcheck() (*PanelHealth, error) {
	resp, err := httpClientShort.Get(panelUrl)
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
	if _, err := httpClientShort.Do(req); err != nil {
		return err
	}
	return nil
}

// Gets the listener's health.
func ListenerHealthcheck() (*ListenerHealth, error) {
	resp, err := httpClientShort.Get(listenerUrl)
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
	if _, err := httpClientShort.Do(req); err != nil {
		return err
	}
	return nil
}

func ListenerCommand(i *discordgo.InteractionCreate) error {
	data, err := json.Marshal(i)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/interaction", listenerUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	res, err := httpClientShort.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("command interaction failed")
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
	res, err := httpClientShort.Do(req)
	if err != nil {
		return false, err
	}
	return res.StatusCode != 400, nil
}

// Schedules the revocation of a punishment
func ListenerExpire(guild string, c int64, t int64) error {
	url := fmt.Sprintf("%s/revoke?guild=%s&case=%d&t=%d", listenerUrl, guild, c, t)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	_, err = httpClientShort.Do(req)
	return err
}

// Requests a pet and loads it in.
func PetOnDemand(target Pet) (string, error) {
	resp, err := httpClientLong.Get(target.source())
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
	resp, err := httpClientLong.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 431 {
		return nil, nil
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
