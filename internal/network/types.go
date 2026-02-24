package network

import (
	"net/netip"
)

type NetInterface struct {
	Name string       `json:"name"`
	IPs  []netip.Addr `json:"ips"`
	Mac  string       `json:"mac,omitempty"`
}
