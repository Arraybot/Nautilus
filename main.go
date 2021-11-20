package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var flagRegister bool
var flagProduction bool
var client *discordgo.Session
var appId string
var token string
var server string
var admins []string
var port string
var startTime time.Time

// Reads the internal Arraybot token as well as the server ID.
func init() {
	// Parse the runtime flags.
	flag.BoolVar(&flagRegister, "register", false, "Whether to re-register all commands")
	flag.BoolVar(&flagProduction, "prod", false, "Whether to run in production mode")
	flag.Parse()
	// Parse the environment variables.
	admins = strings.Split(os.Getenv("ADMINS"), ";")
	appId = os.Getenv("APP_ID")
	server = os.Getenv("DEV_SERVER")
	port = os.Getenv("PORT_COMMANDS")
	startTime = time.Now()
}

// Starts the application.
func main() {
	log.Println("Nautilus starting...")
	log.Printf("Token %s.\n", token)
	// First load the database.
	err := loadDatabase()
	if err != nil {
		log.Fatal(err)
	}
	// Then load the bot.
	err = loadBot(server != "")
	if err != nil {
		log.Fatal(err)
	}
	// Lastly, boot up the HTTP server.
	log.Println("Listening...")
	err = loadHttp()
	if err != nil {
		log.Fatal(err)
	}
}

// Loads the database.
func loadDatabase() error {
	return nil
}

// Loads the bot.
// This also connects to the gateway if development mode is enabled.
// This is because the endpoint can't be served in development mode.
func loadBot(dev bool) error {
	var err error
	// Create the client.
	client, err = discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		return err
	}
	// Check which mode we are running in.
	if flagProduction {
		log.Println("Using production mode; using REST interactions")
		// Check if we need to re-register commands.
		if flagRegister {
			log.Panicln("Force re-register global commands")
			commandsRegister(func(ac *discordgo.ApplicationCommand) error {
				_, err2 := client.ApplicationCommandCreate(appId, "", ac)
				return err2
			})
		}
		return nil
	} else {
		log.Println("Using non-production mode; falling back to WebSocket events")
		// We are in development mode, so the gateway based slash command handler should be used.
		client.AddHandler(slashDistributor)
		// Check if we need to re-register commands.
		if flagRegister {
			log.Printf("Force re-register guild (%s) commands\n", appId)
			commandsRegister(func(ac *discordgo.ApplicationCommand) error {
				_, err2 := client.ApplicationCommandCreate(appId, server, ac)
				return err2
			})
		}
		// Establish WS connection.
		return client.Open()
	}
}

// Loads the HTTP server.
func loadHttp() error {
	http.HandleFunc("/interact", httpInteract)
	http.HandleFunc("/register", httpRegister)
	http.HandleFunc("/unregister", httpUnregister)
	http.HandleFunc("/invalidate", httpInvalidate)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
