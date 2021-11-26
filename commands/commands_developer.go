package commands

import (
	"log"
	"os"
	"strings"

	"github.com/arraybot/nautilus/requests"
	"github.com/bwmarrin/discordgo"
)

// The kill command.
// Will kill one of the provided services.
func handleKill(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasDeveloper(i) {
		s.InteractionRespond(i.Interaction, respondText(permissionDenyDeveloper, i))
		return
	}
	o := i.ApplicationCommandData().Options
	sub(o, "services", func(o []*commandOption) {
		n := option(o, "name").StringValue()
		switch strings.ToLower(n) {
		case "carbon":
			// Kill the web panel.
			var response string
			err := requests.PanelKill()
			if err == nil {
				response = "Carbon terminating..."
			} else {
				log.Println(err)
				response = "Could not send termination request. Please attempt to do so manually."
			}
			s.InteractionRespond(i.Interaction, respondText(response, i))
		case "nautilus":
			// Kill the command handler.
			s.InteractionRespond(i.Interaction, respondText("Nautilus terminating...", i))
			os.Exit(0)
		case "mantis":
			// Kill the gateway listener.
			s.InteractionRespond(i.Interaction, respondText("Mantis terminating...", i))
		default:
			s.InteractionRespond(i.Interaction, respondText("Unknown service.", i))
		}
	})
	sub(o, "shard", func(o []*commandOption) {
		// TODO: Tell Mantis to restart a shard.
		// v := commandGet2(o, "id").IntValue()
		s.InteractionRespond(i.Interaction, respondText("Attempting to re-connect shard, if it exists.", i))
	})
}
