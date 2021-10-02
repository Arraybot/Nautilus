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
func respondEmbed(e *EmbedBuilder, i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	var flag uint64 = 0
	if cacheInvisibility.get(i.GuildID) {
		flag = 64
	}
	return respondEmbedRaw(e, flag)
}

// Helper function that responds to an interaction using an embed and a manual flag.
func respondEmbedRaw(e *EmbedBuilder, flag uint64) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{e.Embed},
			Flags:  flag,
		},
	}
}

// Helper method that responds to an interaction using just text.
func respondText(s string, i *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	var flag uint64 = 0
	if cacheInvisibility.get(i.GuildID) {
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
