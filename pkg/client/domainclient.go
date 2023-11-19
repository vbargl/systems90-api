package client

import (
	"errors"
	"fmt"

	s90api "barglvojtech.net/systems90api/pkg/embi"
)

var (
	// ErrDNSRecordNotFound is returned when the DNS record is not found.
	ErrDNSRecordNotFound = errors.New("dns record not found")
)

// DomainClient is a client for a specific domain.
type DomainClient struct {
	api      *s90api.Systems90Api
	sid      string
	domainID string
}

func (dc *DomainClient) sessionDomain() s90api.SessionDomain {
	return s90api.SessionDomain{
		SID:      dc.sid,
		DomainID: dc.domainID,
	}
}

// AddDNSRecord adds a DNS record.
func (dc *DomainClient) AddDNSRecord(name, value string, typ DNSType, options ...dnsRecordOption) (dnsID string, err error) {
	rec := &s90api.DNSRecord{
		Name: name,
		Type: typ,
		IP:   value,
	}

	applyDNSRecordOptions(rec, options)
	return dc.api.AddDNS(dc.sessionDomain(), rec)
}

// RemoveDNSRecord removes a DNS record.
func (dc *DomainClient) RemoveDNSRecordByID(id string) error {
	return dc.api.DeleteDNS(dc.sid, id)
}

// RemoveDNSRecordByName removes a DNS record.
func (dc *DomainClient) RemoveDNSRecordByName(name string) error {
	dnsRecords, err := dc.api.ListDNS(dc.sessionDomain())
	if err != nil {
		return err
	}

	var rec *s90api.DNSRecord
	for i, r := range dnsRecords {
		if r.Name == name {
			rec = &dnsRecords[i]
		}
	}

	if rec == nil {
		return fmt.Errorf("systems90: %w (%s)", ErrDNSRecordNotFound, name)
	}

	return dc.api.DeleteDNS(dc.sid, rec.ID)
}
