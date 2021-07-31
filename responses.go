package main

import "github.com/bwmarrin/discordgo"

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
