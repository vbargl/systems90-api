package types

// Status is a type that represents the status of a response.
type Status struct {
	Code string `xml:"status"`
	Text string `xml:"text"`
}
