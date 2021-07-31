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

func init() {
	token = os.Getenv("AUTH_TOKEN")
	server = os.Getenv("SERVER")
}

func main() {
	log.Println("Nautilus starting...")
	log.Printf("Token %s.\n", token)
	err := loadDatabase()
	if err != nil {
		log.Fatal(err)
	}
	err = loadBot(server != "")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening...")
	err = loadHttp()
	if err != nil {
		log.Fatal(err)
	}
}

func loadDatabase() error {
	return nil
}

func loadBot(ws bool) error {
	var err error
	client, err = discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		return err
	}
	if ws {
		log.Println("Using WebSocket events.")
		client.AddHandler(slashDistributor)
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

func loadHttp() error {
	http.HandleFunc("/interact", httpInteract)
	http.HandleFunc("/register", httpRegister)
	http.HandleFunc("/unregister", httpUnregister)
	return http.ListenAndServe(":8080", nil)
}
