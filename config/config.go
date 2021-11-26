package config

import (
	"os"
	"strings"
	"time"
)

// Defined in a separate file to prevent cycles.
var Admins []string
var AppID string
var AppPubKey string
var BotToken string
var CommandsSecret string
var DevServer string
var HostMonitor string
var PortCommands string
var PortMonitor string
var SchemeMonitor string
var StartTime time.Time

// Load initializes the environment variables.
func Load() {
	Admins = strings.Split(os.Getenv("ADMINS"), ";")
	AppID = os.Getenv("APP_ID")
	AppPubKey = os.Getenv("APP_PUBKEY")
	BotToken = os.Getenv("BOT_TOKEN")
	CommandsSecret = os.Getenv("COMMANDS_SECRET")
	DevServer = os.Getenv("SERVER")
	HostMonitor = os.Getenv("HOST_MONITOR")
	PortCommands = os.Getenv("PORT_COMMANDS")
	PortMonitor = os.Getenv("PORT_MONITOR")
	SchemeMonitor = os.Getenv("SCHEME_MONITOR")
	StartTime = time.Now()
}
