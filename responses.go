package main

import "github.com/bwmarrin/discordgo"

type EmbedBuilder struct {
	Embed *discordgo.MessageEmbed
}

func embed() *EmbedBuilder {
	author := discordgo.MessageEmbedAuthor{
		URL:     "https://arraybot.xyz",
		Name:    "Arraybot",
		IconURL: "https://i.imgur.com/1JAkQbj.png",
	}
	embed := discordgo.MessageEmbed{
		URL:    "https://arraybot.xyz",
		Color:  0xFFDD57,
		Author: &author,
	}
	return &EmbedBuilder{
		Embed: &embed,
	}
}

// Appends to the description.
func (e *EmbedBuilder) description(s string) {
	e.Embed.Description = e.Embed.Description + s
}

// Appends a field.
func (e *EmbedBuilder) field(t, v string, b bool) {
	field := discordgo.MessageEmbedField{
		Name:   t,
		Value:  v,
		Inline: b,
	}
	e.Embed.Fields = append(e.Embed.Fields, &field)
}

// Helper function that responds to an interaction using an embed.
// Optionally, it can be specified if it should be executed quietly or not.
func respondEmbed(e *EmbedBuilder, q bool) *discordgo.InteractionResponse {
	var flag uint64 = 0
	if q {
		flag = 64
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{e.Embed},
			Flags:  flag,
		},
	}
}

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
