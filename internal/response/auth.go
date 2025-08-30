package response

type SignIn struct {
	Domain    string `json:"domain"`
	SessionID string `json:"sessionID"`
}

type UserProfile struct {
}
