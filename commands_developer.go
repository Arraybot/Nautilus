package main

import (
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// The kill command.
// Will kill one of the provided services.
func handleKill(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !commandPermissionDeveloper(i) {
		s.InteractionRespond(i.Interaction, respondText(permissionDenyDeveloper, i))
		return
	}
	o := i.ApplicationCommandData().Options
	commandWhen(o, "services", func(o []*commandOption) {
		n := commandGet2(o, "name").StringValue()
		switch strings.ToLower(n) {
		case "carbon":
			// Restart the panel.
			var response string
			err := requestPanelKill()
			if err == nil {
				response = "Carbon terminating..."
			} else {
				log.Println(err)
				response = "Could not send termination request. Please attempt to do so manually."
			}
			s.InteractionRespond(i.Interaction, respondText(response, i))
		case "nautilus":
			// Restart the command handler.
			s.InteractionRespond(i.Interaction, respondText("Restarting Nautilus...", i))
			os.Exit(0)
		case "mantis":
			// TODO: Restart the listener.
			s.InteractionRespond(i.Interaction, respondText("Restarting Mantis...", i))
		default:
			s.InteractionRespond(i.Interaction, respondText("Unknown service.", i))
		}
	})
	commandWhen(o, "shard", func(o []*commandOption) {
		// TODO: Tell Mantis to restart a shard.
		// v := commandGet2(o, "id").IntValue()
		s.InteractionRespond(i.Interaction, respondText("Attempting to re-connect shard, if it exists.", i))
	})
}
