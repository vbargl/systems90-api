package types

type LoginRequest struct {
	UID      string `urlparams:"uid"`
	Password string `urlparams:"password"`
}

type LoginResponse struct {
	Status Status `xml:"status"`
	UID    string `xml:"uid"`
	SID    string `xml:"sid"`
}
