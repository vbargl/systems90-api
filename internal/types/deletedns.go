package types

type DeleteDnsRequest struct {
	SID   string `urlparams:"sid"`
	DNSID string `urlparams:"dns_id"`
}

type DeleteDnsResponse struct {
	Status Status `xml:"status"`
}
