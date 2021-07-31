package main

import "github.com/bwmarrin/discordgo"

// Helper method that responds to an interaction using just text.
// Optionally, it can be specified if it should be executed quietly or not.
func respondText(s string, q bool) *discordgo.InteractionResponse {
	var flag uint64 = 0
	if q {
		flag = 64
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: s,
			Flags:   flag,
		},
	}
}
