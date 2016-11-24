package steamauth

// LinkSuccessEvent is a struct for the event called when an account is successfully linked.
type LinkSuccessEvent struct {
	Secret    string
	DiscordID string
	SteamID   string
}

// LinkFailureEvent is a struct for the event called when an account is unsuccessfully linked.
type LinkFailureEvent struct {
	Secret    string
	DiscordID string
}

// LinkAttemptEvent is a struct for the event called when an account link has been attempted.
type LinkAttemptEvent struct {
	Secret    string
	DiscordID string
}
