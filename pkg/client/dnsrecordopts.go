package client

import (
	"time"

	s90api "barglvojtech.net/systems90api/pkg/embi"
)

func applyDNSRecordOptions(rec *s90api.DNSRecord, options []dnsRecordOption) {
	for _, opt := range dnsRecordDefaults {
		opt(rec)
	}
	for _, opt := range options {
		opt(rec)
	}
}

type dnsRecordOption func(*s90api.DNSRecord)

var dnsRecordDefaults = []dnsRecordOption{
	DNSRecordTTL(30 * time.Second),
}

// DNSRecordTTL sets the TTL of the DNS record.
func DNSRecordTTL(duration time.Duration) dnsRecordOption {
	return func(rec *s90api.DNSRecord) {
		rec.TTL = duration
	}
}

// DNSRecordPriority sets the priority of the DNS record.
func DNSRecordPriority(priority int) dnsRecordOption {
	return func(rec *s90api.DNSRecord) {
		rec.Priority = priority
	}
}
