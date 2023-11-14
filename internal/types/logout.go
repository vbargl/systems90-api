package types

type LogoutRequest struct {
	SID string `urlparams:"sid"`
}

type LogoutResponse struct {
	Status Status `xml:"status"`
}
