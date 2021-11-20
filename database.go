package main

// Gets the mode to respond to.
func databaseReplyHidden(id string) bool {
	return true
}

// Gets the server guide.
func databaseGuide(id string) string {
	return ""
}

// Gets disabled commands.
func databaseDisabled(id string) []string {
	return []string{}
}
