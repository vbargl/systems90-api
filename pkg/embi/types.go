package embi

import "time"

type Credentials struct {
	UID      string
	Password string
}

type DomainMap map[string]string

type DNSRecord struct {
	ID       string
	Name     string
	TTL      time.Duration
	Type     DNSType
	IP       string
	Priority int
	Locked   bool
}

const (
	DNSTypeA     DNSType = "A"
	DNSTypeAAAA  DNSType = "AAAA"
	DNSTypeCNAME DNSType = "CNAME"
	DNSTypeDNAME DNSType = "DNAME"
	DNSTypeLOC   DNSType = "LOC"
	DNSTypeMX    DNSType = "MX"
	DNSTypeNS    DNSType = "NS"
	DNSTypeSRV   DNSType = "SRV"
	DNSTypeSSHFP DNSType = "SSHFP"
	DNSTypeTXT   DNSType = "TXT"
	DNSTypeCAA   DNSType = "CAA"
)

type DNSType string

func (DNSType) sealed() {}

func (t DNSType) String() string {
	return string(t)
}

func DNSTypeFromString(s string) DNSType {
	switch s {
	case "A":
		return DNSTypeA
	case "AAAA":
		return DNSTypeAAAA
	case "CNAME":
		return DNSTypeCNAME
	case "DNAME":
		return DNSTypeDNAME
	case "LOC":
		return DNSTypeLOC
	case "MX":
		return DNSTypeMX
	case "NS":
		return DNSTypeNS
	case "SRV":
		return DNSTypeSRV
	case "SSHFP":
		return DNSTypeSSHFP
	case "TXT":
		return DNSTypeTXT
	case "CAA":
		return DNSTypeCAA
	default:
		return ""
	}
}
