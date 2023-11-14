package urlparams

import (
	"errors"
	"strings"
)

type Flags int

const (
	OmitEmptyFlag Flags = 1 << iota
)

func parseTag(data string) (name string, flags Flags, err error) {
	parts := strings.Split(data, ",")

Parsing:
	switch len(parts) {
	case 0:
		err = errors.New("urlparams: empty tag")

	case 1:
		name = strings.TrimSpace(parts[0])

	default:
		name = strings.TrimSpace(parts[0])
		for _, flag := range parts[1:] {
			switch strings.ToLower(strings.TrimSpace(flag)) {
			case "omitempty":
				flags |= OmitEmptyFlag
			default:
				err = errors.New("urlparams: unknown flag")
				break Parsing
			}
		}
	}

	if err != nil {
		return "", 0, err
	}

	return name, flags, nil
}
