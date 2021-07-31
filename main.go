package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

var client *discordgo.Session
var token string
var server string

// Reads the internal Arraybot token as well as the server ID.
func init() {
	token = os.Getenv("AUTH_TOKEN")
	server = os.Getenv("SERVER")
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
	if dev {
		log.Println("Using WebSocket events.")
		// Use gateway based slash command handling.
		client.AddHandler(slashDistributor)
		// Register base slash commands to Discord.
		for _, command := range commands {
			log.Printf("Registering command %s.\n", command.appCommand.Name)
			_, err = client.ApplicationCommandCreate(os.Getenv("APP_ID"), server, command.appCommand)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		return client.Open()
	}
	return nil
}

// Loads the HTTP server.
func loadHttp() error {
	http.HandleFunc("/interact", httpInteract)
	http.HandleFunc("/register", httpRegister)
	http.HandleFunc("/unregister", httpUnregister)
	return http.ListenAndServe(":8080", nil)
}
