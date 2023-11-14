package request

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

func Fetch[T any](client *http.Client, fn func(*Builder)) (*T, error) {
	b := &Builder{}
	fn(b)

	req, err := b.build()
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var val T
	defer resp.Body.Close()
	if err := xml.NewDecoder(resp.Body).Decode(&val); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return &val, fmt.Errorf("request: failed %s: %s", req.URL, resp.Status)
	}
	return &val, nil
}
