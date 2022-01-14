package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func handlerClear(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasPermission(i, discordgo.PermissionManageMessages) {
		s.InteractionRespond(i.Interaction, respondText(permissionDenyPermission, i))
		return
	}
	// Load all the command options.
	count := option(i.ApplicationCommandData().Options, "amount")
	c := count.IntValue()
	user := option(i.ApplicationCommandData().Options, "user")
	var u *discordgo.User = nil
	if user != nil {
		u = user.UserValue(nil)
	}
	// Validate that the user range is valid.
	if c < 2 || c > 100 {
		s.InteractionRespond(i.Interaction, respondText("You can only clear at least 2 and at most 100 messages.", i))
		return
	}
	// Try to get the IDs of messages to delete.
	var ids []string
	// Try to load all messages.
	messages, err := s.ChannelMessages(i.ChannelID, int(c), "", "", "")
	if err != nil {
		log.Println(err)
		s.InteractionRespond(i.Interaction, respondText("Error getting messages.", i))
		return
	}
	// Filter to only eligible messages.
	nt := time.Now()
	for _, message := range messages {
		ct := creationTime(message.ID)
		diff := nt.Sub(ct)
		hrs := diff.Hours()
		wks := hrs / (24 * 7)
		// Only include messages less than 2 weeks old.
		if wks <= 2.0 {
			// If user filtering is enabled, do that.
			if u != nil && message.Author.ID != u.ID {
				continue
			}
			ids = append(ids, message.ID)
		}
	}
	// Delete all the messages.
	err = s.ChannelMessagesBulkDelete(i.ChannelID, ids)
	if err != nil {
		log.Println(err)
		s.InteractionRespond(i.Interaction, respondText("An error ocurred deleting messages", i))
	} else {
		s.InteractionRespond(i.Interaction, respondText(fmt.Sprintf("Deleted %d messages.", len(ids)), i))
	}
}

func handlerMute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasMutePermission(i) {
		s.InteractionRespond(i.Interaction, respondText(permissionDenyMuteRole, i))
		return
	}
	s.InteractionRespond(i.Interaction, respondText("Mute", i))
}

func handlerUnMute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasMutePermission(i) {
		s.InteractionRespond(i.Interaction, respondText(permissionDenyMuteRole, i))
		return
	}
	s.InteractionRespond(i.Interaction, respondText("Unmute", i))
}

func handlerExpire(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: Look up punishment type and determine permission dynamically.
	s.InteractionRespond(i.Interaction, respondText("Expire", i))
}

func handlerRevoke(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: Look up punishment type and determine permission dynamically.
	s.InteractionRespond(i.Interaction, respondText("Revoke", i))
}

func handlerLookup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasModerator(i) {
		s.InteractionRespond(i.Interaction, respondText(permissionDenyModerator, i))
		return
	}
	s.InteractionRespond(i.Interaction, respondText("Lookup", i))
}

func handlerHistory(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasModerator(i) {
		s.InteractionRespond(i.Interaction, respondText(permissionDenyModerator, i))
		return
	}
	s.InteractionRespond(i.Interaction, respondText("History", i))
}
