package types

type ListDnsRequest struct {
	SID      string `urlparams:"sid"`
	DomainID string `urlparams:"domain_id"`
}

type ListDnsResponse struct {
	Status Status               `xml:"status"`
	Zone   ListDnsResponse_Zone `xml:"zone"`
}

type ListDnsResponse_Zone struct {
	Records []ListDnsResponse_Record `xml:"record"`
}

type ListDnsResponse_Record struct {
	DnsID    string `xml:"dns_id"`
	Name     string `xml:"name"`
	TTL      string `xml:"ttl"`
	Type     string `xml:"type"`
	IP       string `xml:"ip"`
	Priority string `xml:"priority"`
	Locked   bool   `xml:"locked"`
}
