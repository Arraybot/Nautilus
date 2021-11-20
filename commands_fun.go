package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var eightballResponses = []string{
	"It is certain.", "It is decidedly so.", "Without a doubt.", "Yes definitely.", "You may rely on it.",
	"As I see it, yes.", "Most likely.", "Outlook good.", "Yes.", "Signs point to yes.",
	"Reply hazy, try again.", "Ask again later.", "Better not tell you now.", "Cannot predict now.", "Concentrate and ask again.",
	"Don't count on it.", "My reply is no.", "My sources say no.", "Outlook not so good.", "Very doubtful.",
}

// The cat command.
// Sends a link to a cute cat.
func handleCat(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cat := cat{}
	link, err := requestPet(&cat)
	var toSend string
	if err == nil {
		toSend = link
	} else {
		toSend = "Error getting cat."
		log.Println(err)
	}
	s.InteractionRespond(i.Interaction, respondText(toSend, i))
}

// The dog command.
// Sends a link to a cute dog.
func handleDog(s *discordgo.Session, i *discordgo.InteractionCreate) {
	dog := dog{}
	link, err := requestPet(&dog)
	var toSend string
	if err == nil {
		toSend = link
	} else {
		toSend = "Error getting dog."
		log.Println(err)
	}
	s.InteractionRespond(i.Interaction, respondText(toSend, i))
}

// The eightball command.
// Repeats the user input and gives its opinion on it.
func handleEightBall(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get the question.
	q := commandGet2(i.ApplicationCommandData().Options, "question")
	// Truncate if it is really long.
	str := q.StringValue()
	if len(str) > 1024 {
		str = str[0:1024]
	}
	respose := fmt.Sprintf("> %s\n\n**%s**", str, eightballResponses[rand.Intn(len(eightballResponses))])
	s.InteractionRespond(i.Interaction, respondText(respose, i))
}

// The urban command.
// Looks up the given phrase in the Urban Dictionary and returns the definition.
func handleUrban(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p := commandGet2(i.ApplicationCommandData().Options, "phrase")
	str := url.QueryEscape(p.StringValue())
	var toSend string
	definition, err := requestUrban(str)
	if err == nil {
		if definition != nil {
			toSend = fmt.Sprintf("> %s\n\n**👍 %d | %d 👎**", definition.Description, definition.Upvotes, definition.Downvotes)
			toSend = strings.Replace(toSend, "[", "", -1)
			toSend = strings.Replace(toSend, "]", "", -1)
			if len(toSend) > 1024 {
				toSend = toSend[0:1024]
			}
		} else {
			toSend = "No definition found. Maybe you can add it!"
		}

	} else {
		toSend = "Error getting definition."
		log.Println(err)
	}
	s.InteractionRespond(i.Interaction, respondText(toSend, i))
}
