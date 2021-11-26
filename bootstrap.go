package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arraybot/nautilus/commands"
	"github.com/arraybot/nautilus/config"
	"github.com/bwmarrin/discordgo"
)

var flagRegister bool
var flagProduction bool
var flagNoServe bool
var client *discordgo.Session

// Reads the config.
func init() {
	// Parse the runtime flags.
	flag.BoolVar(&flagRegister, "register", false, "Whether to re-register all commands")
	flag.BoolVar(&flagProduction, "prod", false, "Whether to run in production mode")
	flag.BoolVar(&flagNoServe, "noserve", false, "Whether to only register commands (if enabled) and nothing else")
	flag.Parse()
	// Parse the environment variables.
	config.Load()
}

// Starts the application.
func main() {
	log.Println("Nautilus starting...")
	if !flagNoServe {
		// First load the database.
		log.Println("Establishing database connection")
		err := loadDatabase()
		if err != nil {
			log.Fatal(err)
		}
	}
	// Then load the bot.
	err := loadBot()
	if err != nil {
		log.Fatal(err)
	}
	if !flagNoServe {
		// Lastly, boot up the HTTP server.
		log.Println("Listening...")
		err = loadHttp()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Loads the database.
func loadDatabase() error {
	return nil
}

// Loads the bot.
// This also connects to the gateway if development mode is enabled.
// This is because the endpoint can't be served in development mode.
func loadBot() error {
	var err error
	// Create the client.
	client, err = discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		return err
	}
	appId := config.AppID
	server := config.DevServer
	// Check which mode we are running in.
	if flagProduction {
		log.Println("Using production mode; using REST interactions")
		// Check if we need to re-register commands.
		if flagRegister {
			log.Println("Force re-register global commands")
			commands.Register(func(ac *discordgo.ApplicationCommand) error {
				_, err2 := client.ApplicationCommandCreate(appId, "", ac)
				return err2
			})
		}
		return nil
	} else {
		log.Println("Using non-production mode; falling back to WebSocket events")
		// We are in development mode, so the gateway based slash command handler should be used.
		client.AddHandler(commands.Distributor)
		// Check if we need to re-register commands.
		if flagRegister {
			log.Printf("Force re-register guild (%s) commands\n", server)
			commands.Register(func(ac *discordgo.ApplicationCommand) error {
				_, err2 := client.ApplicationCommandCreate(appId, server, ac)
				return err2
			})
		}
		// If we don't want to serve, skip opening connection.
		if flagNoServe {
			log.Println("Not connecting to WebSocket as configured")
			return nil
		}
		// Establish WS connection.
		return client.Open()
	}
}

// Loads the HTTP server.
func loadHttp() error {
	http.HandleFunc("/interact", epInteract)
	http.HandleFunc("/register", epRegister)
	http.HandleFunc("/unregister", epUnregister)
	http.HandleFunc("/invalidate", epInvalidate)
	return http.ListenAndServe(fmt.Sprintf(":%s", config.PortCommands), nil)
}
