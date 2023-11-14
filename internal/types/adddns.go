package types

type AddDnsRequest_UrlParams struct {
	SID      string `urlparams:"sid"`
	DomainID string `urlparams:"domain_id"`
}

type AddDnsRequest_Payload struct {
	Name     string `urlparams:"name"`
	TTL      string `urlparams:"ttl"`
	Type     string `urlparams:"type"`
	IP       string `urlparams:"ip"`
	Priority string `urlparams:"priority,omitempty"`
}

type AddDnsResponse struct {
	Status Status `xml:"status"`
	DNSID  string `xml:"dns_id"`
}
