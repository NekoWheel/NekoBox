package response

type SignIn struct {
	Profile   *SignInUserProfile `json:"profile"`
	SessionID string             `json:"sessionID"`
}

type SignInUserProfile struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}
