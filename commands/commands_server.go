package commands

import (
	"github.com/arraybot/nautilus/database"
	"github.com/bwmarrin/discordgo"
)

// The guide command.
// This shows the guide set by the server administrators.
func handlerGuide(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: Get the guide.
	guide := database.Guide(i.GuildID)
	if guide == "" {
		guide = "The server's administrator(s) has/have not set a guide yet. Perhaps ask them to?"
	}
	s.InteractionRespond(i.Interaction, respondText(guide, i))
}
