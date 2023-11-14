package request

import (
	"errors"
	"net/http"
	"net/url"
	"path"
	"strings"

	"barglvojtech.net/webglobeapi/internal/urlparams"
)

type Builder struct {
	errs []error

	url    *url.URL
	header http.Header
	method string
	path   string

	urlParam string
	payload  string
}

func (b *Builder) Url(rawUrl string) {
	var err error
	b.url, err = url.ParseRequestURI(rawUrl)
	b.errs = append(b.errs, err)
}

func (b *Builder) Header(header http.Header) {
	b.header = header
}

func (b *Builder) Method(method string) {
	b.method = method
}

func (b *Builder) Path(path string) {
	b.path = path
}

func (b *Builder) UrlParams(params any) {
	switch params.(type) {
	case string:
		b.urlParam = params.(string)
	default:
		var err error
		b.urlParam, err = urlparams.Marshal(params)
		if err != nil {
			b.errs = append(b.errs, err)
		}
	}
}

func (b *Builder) Payload(payload any) {
	switch payload.(type) {
	case string:
		b.payload = payload.(string)
	default:
		var err error
		b.payload, err = urlparams.Marshal(payload)
		if err != nil {
			b.errs = append(b.errs, err)
		}
	}
}

func (b *Builder) build() (*http.Request, error) {
	if err := errors.Join(b.errs...); err != nil {
		return nil, err
	}

	url := *b.url
	url.Path = path.Join(url.Path, b.path)
	url.RawQuery = b.urlParam

	req, err := http.NewRequest(b.method, url.String(), strings.NewReader(b.payload))
	if err != nil {
		return nil, err
	}

	req.Header = b.header

	return req, nil
}
