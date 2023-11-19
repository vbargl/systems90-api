package client

import (
	"errors"
	"fmt"

	s90api "barglvojtech.net/systems90api/pkg/embi"
)

var (
	// ErrDomainNotManaged is returned when the domain is not managed by UID.
	ErrDomainNotManaged = errors.New("domain not managed")
)

type Client struct {
	api *s90api.Systems90Api
	sid string
}

// NewClient creates a new client for the Systems90 API.
func NewClient(cred s90api.Credentials) (*Client, error) {
	api := s90api.NewSystems90Api()
	sid, err := api.Login(cred)

	if err != nil {
		return nil, err
	}

	return &Client{
		api: api,
		sid: sid,
	}, nil
}

// Close closes the client and logs out from the API.
func (c *Client) Close() error {
	return c.api.Logout(c.sid)
}

// Domain returns a DomainClient for the given domain.
func (c *Client) Domain(zone string) (*DomainClient, error) {
	domains, err := c.api.ListDomains(c.sid)

	if err != nil {
		return nil, err
	}

	var domainID string
	for _, d := range domains {
		if d.Zone == zone {
			domainID = d.DomainID
			break
		}
	}

	if domainID == "" {
		return nil, fmt.Errorf("systems90: %w (%s)", ErrDomainNotManaged, zone)
	}

	return &DomainClient{
		api:      c.api,
		sid:      c.sid,
		domainID: domainID,
	}, nil
}
