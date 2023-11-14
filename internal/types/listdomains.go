package types

type ListDomainsRequest struct {
	SID string `urlparams:"sid"`
}

type ListDomainsResponse struct {
	Status  Status                      `xml:"status"`
	Domains ListDomainsResponse_Domains `xml:"domains"`
}

type ListDomainsResponse_Domains struct {
	Domains []ListDomainsResponse_Domain `xml:"domain"`
}

type ListDomainsResponse_Domain struct {
	DomainID string `xml:"domain_id"`
	Name     string `xml:"name"`
}
