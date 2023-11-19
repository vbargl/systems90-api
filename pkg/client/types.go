package client

import (
	s90api "barglvojtech.net/systems90api/pkg/embi"
)

type Credentials = s90api.Credentials
type DNSRecord = s90api.DNSRecord

type DNSType = s90api.DNSType

const (
	DNSTypeA     = s90api.DNSTypeA
	DNSTypeAAAA  = s90api.DNSTypeAAAA
	DNSTypeCNAME = s90api.DNSTypeCNAME
	DNSTypeDNAME = s90api.DNSTypeDNAME
	DNSTypeLOC   = s90api.DNSTypeLOC
	DNSTypeMX    = s90api.DNSTypeMX
	DNSTypeNS    = s90api.DNSTypeNS
	DNSTypeSRV   = s90api.DNSTypeSRV
	DNSTypeSSHFP = s90api.DNSTypeSSHFP
	DNSTypeTXT   = s90api.DNSTypeTXT
	DNSTypeCAA   = s90api.DNSTypeCAA
)
