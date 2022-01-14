package database

// Disabled gets disabled commands.
func Disabled(id string) []string {
	return []string{}
}

// ReplyHidden gets the mode to respond to.
func ReplyHidden(id string) bool {
	return true
}

// MutePermission gets the mute permission role ID of the server.
func MutePermission(id string) string {
	return ""
}

// Guide gets the server guide.
func Guide(id string) string {
	return ""
}
