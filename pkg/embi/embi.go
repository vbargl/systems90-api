package embi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"barglvojtech.net/systems90api/internal/request"
	"barglvojtech.net/systems90api/internal/types"
)

var (
	ErrInvalidSession = errors.New("to use this endpoint, you need to be logged in")
)

// Systems90Api is a basic interface for communication with the Systems90 API.
type Systems90Api struct {
	client         *http.Client
	requestBuilder func(*request.Builder)
}

// NewSystems90Api creates a new Systems90Api.
func NewSystems90Api() *Systems90Api {
	return &Systems90Api{
		client: &http.Client{},
		requestBuilder: func(b *request.Builder) {
			b.Url("https://admin.systems90.cz/api/")
			b.Header(http.Header{
				"Content-Type": []string{"application/x-www-form-urlencoded"},
			})
		},
	}
}

// Login logs in to the API and returns a session ID.
func (api *Systems90Api) Login(cred Credentials) (SID string, err error) {
	resp, err := request.Fetch[types.LoginResponse](api.client, func(b *request.Builder) {
		api.requestBuilder(b)
		b.Method(http.MethodPost)
		b.Path("login")
		b.Payload(types.LoginRequest{
			UID:      cred.UID,
			Password: cred.Password,
		})
	})

	switch {
	case err != nil && resp != nil:
		return "", fmt.Errorf("systems90: %w: %s", err, resp.Status.Text)
	case err != nil:
		return "", err
	}

	return resp.SID, nil
}

// Logout invalidates the session ID.
func (api *Systems90Api) Logout(SID string) error {
	if isInvalidSession(SID) {
		return ErrInvalidSession
	}

	resp, err := request.Fetch[types.LogoutResponse](api.client, func(b *request.Builder) {
		api.requestBuilder(b)
		b.Path("logout")
		b.UrlParams(types.LogoutRequest{
			SID: SID,
		})
	})

	switch {
	case err != nil && resp != nil:
		return fmt.Errorf("systems90: %w: %s", err, resp.Status.Text)
	case err != nil:
		return err
	}

	return nil
}

// ListDomains lists all domains registred for logged in UID.
func (api *Systems90Api) ListDomains(SID string) ([]Domain, error) {
	if isInvalidSession(SID) {
		return nil, ErrInvalidSession
	}

	resp, err := request.Fetch[types.ListDomainsResponse](api.client, func(b *request.Builder) {
		api.requestBuilder(b)
		b.Path("domain_list")
		b.Method(http.MethodGet)
		b.UrlParams(types.ListDomainsRequest{
			SID: SID,
		})
	})

	switch {
	case err != nil && resp != nil:
		return nil, fmt.Errorf("systems90: %w: %s", err, resp.Status.Text)
	case err != nil:
		return nil, err
	}

	domains := make([]Domain, len(resp.Domains.Domains))
	for i, domain := range resp.Domains.Domains {
		domains[i] = Domain{
			DomainID: domain.DomainID,
			Zone:     domain.Name,
		}
	}
	return domains, nil
}

// ListDNS lists all DNS records for the given domain.
func (api *Systems90Api) ListDNS(sd SessionDomain) ([]DNSRecord, error) {
	if isInvalidSession(sd) {
		return nil, ErrInvalidSession
	}

	resp, err := request.Fetch[types.ListDnsResponse](api.client, func(b *request.Builder) {
		api.requestBuilder(b)
		b.Path("domain_list_dns")
		b.Method(http.MethodGet)
		b.UrlParams(types.ListDnsRequest{
			SID:      sd.SID,
			DomainID: sd.DomainID,
		})
	})

	switch {
	case err != nil && resp != nil:
		return nil, fmt.Errorf("systems90: %w: %s", err, resp.Status.Text)
	case err != nil:
		return nil, err
	}

	records := make([]DNSRecord, len(resp.Zone.Records))
	for i, rec := range resp.Zone.Records {
		ttl, err := strconv.ParseUint(rec.TTL, 10, 64)
		if err != nil {
			ttl = 0
		}

		priority, err := strconv.ParseUint(rec.Priority, 10, 64)
		if err != nil {
			priority = 0
		}

		records[i] = DNSRecord{
			ID:       rec.DnsID,
			Name:     rec.Name,
			TTL:      time.Duration(ttl) * time.Second,
			Type:     DNSTypeFromString(rec.Type),
			IP:       rec.IP,
			Priority: int(priority),
			Locked:   rec.Locked,
		}
	}
	return records, nil
}

// AddDNS adds a new DNS record to the given domain.
func (api *Systems90Api) AddDNS(sd SessionDomain, dnsRecord *DNSRecord) (string, error) {
	if isInvalidSession(sd) {
		return "", ErrInvalidSession
	}

	resp, err := request.Fetch[types.AddDnsResponse](api.client, func(b *request.Builder) {
		api.requestBuilder(b)
		b.Path("domain_add_dns")
		b.Method(http.MethodPost)
		b.UrlParams(types.AddDnsRequest_UrlParams{
			SID:      sd.SID,
			DomainID: sd.DomainID,
		})
		b.Payload(types.AddDnsRequest_Payload{
			Name:     dnsRecord.Name,
			TTL:      strconv.FormatUint(uint64(dnsRecord.TTL.Seconds()), 10),
			Type:     dnsRecord.Type.String(),
			IP:       dnsRecord.IP,
			Priority: strconv.FormatInt(int64(dnsRecord.Priority), 10),
		})
	})

	switch {
	case err != nil && resp != nil:
		return "", fmt.Errorf("systems90: %w: %s", err, resp.Status.Text)
	case err != nil:
		return "", err
	}

	return resp.DNSID, nil
}

// DeleteDNS deletes a DNS record from the given domain.
func (api *Systems90Api) DeleteDNS(sid string, dnsRecordId string) error {
	if isInvalidSession(sid) {
		return ErrInvalidSession
	}

	resp, err := request.Fetch[types.DeleteDnsResponse](api.client, func(b *request.Builder) {
		api.requestBuilder(b)
		b.Path("domain_delete_dns")
		b.Method(http.MethodGet)
		b.UrlParams(types.DeleteDnsRequest{
			SID:   sid,
			DNSID: dnsRecordId,
		})
	})

	switch {
	case err != nil && resp != nil:
		return fmt.Errorf("systems90: %w: %s", err, resp.Status.Text)
	case err != nil:
		return err
	}

	return nil
}
