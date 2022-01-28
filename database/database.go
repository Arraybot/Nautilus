package database

// Punishment is a wrapper for a punishment in the database.
type Punishment struct {
	ID     int
	Type   string
	User   int64
	Staff  int64
	Reason string
	Expiry int64
	Log    int64
}

// Disabled gets disabled commands.
func Disabled(id string) []string {
	return []string{}
}

// ReplyHidden gets the mode to respond to.
func ReplyHidden(id string) bool {
	return true
}

// MuteRole gets the mute role ID of the server.
func MuteRole(id string) string {
	return "419474184536850435"
}

// MutePermission gets the mute permission role ID of the server.
func MutePermission(id string) string {
	return "388004558909210624"
}

// Guide gets the server guide.
func Guide(id string) string {
	return ""
}

func GetPunishment(guild, id string) *Punishment {
	return &Punishment{
		ID:   1,
		Type: "KICK",
	}
}
