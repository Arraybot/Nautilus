package database

// Disabled gets disabled commands.
func Disabled(id string) []string {
	return []string{}
}

// ReplyHidden gets the mode to respond to.
func ReplyHidden(id string) bool {
	return true
}

// Moderator gets the moderator role ID of the server.
func Moderator(id string) string {
	return ""
}

// Guide gets the server guide.
func Guide(id string) string {
	return ""
}
