package embi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"barglvojtech.net/webglobeapi/internal/request"
	"barglvojtech.net/webglobeapi/internal/types"
)

var (
	ErrNotLoggedIn = errors.New("to use this endpoint, you need to be logged in")
)

type WebGlobeApi struct {
	sess   *session
	client *http.Client

	requestBuilder func(*request.Builder)
}

type session struct {
	uid string // user id
	sid string // session id
}

func NewWebGlobeApi() *WebGlobeApi {
	return &WebGlobeApi{
		client: &http.Client{},
		requestBuilder: func(b *request.Builder) {
			b.Url("https://admin.systems90.cz/api/")
			b.Header(http.Header{
				"Content-Type": []string{"application/x-www-form-urlencoded"},
			})
		},
	}
}

func (api *WebGlobeApi) hasSession() bool {
	if api.sess == nil {
		return false
	}

	return true
}

func (api *WebGlobeApi) Login(cred Credentials) error {
	if api.hasSession() {
		return errors.New("already logged in")
	}

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
		return fmt.Errorf("%w: %s", err, resp.Status.Text)
	case err != nil:
		return err
	}

	api.sess = &session{
		uid: resp.UID,
		sid: resp.SID,
	}
	return nil
}

func (api *WebGlobeApi) Logout() error {
	if !api.hasSession() {
		return ErrNotLoggedIn
	}

	resp, err := request.Fetch[types.LogoutResponse](api.client, func(b *request.Builder) {
		api.requestBuilder(b)
		b.Path("logout")
		b.UrlParams(types.LogoutRequest{
			SID: api.sess.sid,
		})
	})

	switch {
	case err != nil && resp != nil:
		return fmt.Errorf("%w: %s", err, resp.Status.Text)
	case err != nil:
		return err
	}

	api.sess = nil
	return nil
}

func (api *WebGlobeApi) ListDomains() (DomainMap, error) {
	if !api.hasSession() {
		return nil, ErrNotLoggedIn
	}

	resp, err := request.Fetch[types.ListDomainsResponse](api.client, func(b *request.Builder) {
		api.requestBuilder(b)
		b.Path("domain_list")
		b.Method(http.MethodGet)
		b.UrlParams(types.ListDomainsRequest{
			SID: api.sess.sid,
		})
	})

	switch {
	case err != nil && resp != nil:
		return nil, fmt.Errorf("%w: %s", err, resp.Status.Text)
	case err != nil:
		return nil, err
	}

	domains := make(map[string]string)
	for _, domain := range resp.Domains.Domains {
		domains[domain.Name] = domain.DomainID
	}
	return domains, nil
}

func (api *WebGlobeApi) ListDNS(domainID string) ([]DNSRecord, error) {
	if !api.hasSession() {
		return nil, ErrNotLoggedIn
	}

	resp, err := request.Fetch[types.ListDnsResponse](api.client, func(b *request.Builder) {
		api.requestBuilder(b)
		b.Path("domain_list_dns")
		b.Method(http.MethodGet)
		b.UrlParams(types.ListDnsRequest{
			SID:      api.sess.sid,
			DomainID: domainID,
		})
	})

	switch {
	case err != nil && resp != nil:
		return nil, fmt.Errorf("%w: %s", err, resp.Status.Text)
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

func (api *WebGlobeApi) AddDNS(domainID string, dnsRecord *DNSRecord) (string, error) {
	if !api.hasSession() {
		return "", ErrNotLoggedIn
	}

	resp, err := request.Fetch[types.AddDnsResponse](api.client, func(b *request.Builder) {
		api.requestBuilder(b)
		b.Path("domain_add_dns")
		b.Method(http.MethodPost)
		b.UrlParams(types.AddDnsRequest_UrlParams{
			SID:      api.sess.sid,
			DomainID: domainID,
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
		return "", fmt.Errorf("%w: %s", err, resp.Status.Text)
	case err != nil:
		return "", err
	}

	return resp.DNSID, nil
}

func (api *WebGlobeApi) DeleteDNS(domainID string, dnsRecordId string) error {
	if !api.hasSession() {
		return ErrNotLoggedIn
	}

	resp, err := request.Fetch[types.DeleteDnsResponse](api.client, func(b *request.Builder) {
		api.requestBuilder(b)
		b.Path("domain_delete_dns")
		b.Method(http.MethodGet)
		b.UrlParams(types.DeleteDnsRequest{
			SID:   api.sess.sid,
			DNSID: dnsRecordId,
		})
	})

	switch {
	case err != nil && resp != nil:
		return fmt.Errorf("%w: %s", err, resp.Status.Text)
	case err != nil:
		return err
	}

	return nil
}
