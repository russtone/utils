package dns

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

const port = 53

// Supported DNS queries types.
const (
	TypeA     = "A"
	TypeAAAA  = "AAAA"
	TypeNS    = "NS"
	TypeMX    = "MX"
	TypeTXT   = "TXT"
	TypeSRV   = "SRV"
	TypeCNAME = "CNAME"
	TypePTR   = "PTR"
)

var (
	// Types are all supported DNS query types.
	Types = []string{TypeA, TypeAAAA, TypeNS, TypeMX, TypeTXT, TypeSRV, TypeCNAME, TypePTR}

	Client = &dns.Client{
		Timeout: time.Second * 3,
	}

	qtypeMap = map[string]uint16{
		TypeA:     dns.TypeA,
		TypeAAAA:  dns.TypeAAAA,
		TypeNS:    dns.TypeNS,
		TypeMX:    dns.TypeMX,
		TypeTXT:   dns.TypeTXT,
		TypeSRV:   dns.TypeSRV,
		TypeCNAME: dns.TypeCNAME,
		TypePTR:   dns.TypePTR,
	}
)

// Query returns DNS query result for the given name using given NS.
func Query(name string, ns net.IP, qtype string) ([]string, error) {

	msg := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Id:               dns.Id(),
			RecursionDesired: true,
		},
		Question: []dns.Question{
			{
				Name:   name + ".",
				Qtype:  qtypeMap[qtype],
				Qclass: dns.ClassINET,
			},
		},
	}

	in, _, err := Client.Exchange(msg, getAddr(ns))
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)

	switch qtype {

	case TypeA:
		for _, a := range in.Answer {
			if t, ok := a.(*dns.A); ok {
				res = append(res, t.A.String())
			}
		}

	case TypeAAAA:
		for _, a := range in.Answer {
			if t, ok := a.(*dns.AAAA); ok {
				res = append(res, t.AAAA.String())
			}
		}

	case TypeNS:
		for _, a := range in.Answer {
			if t, ok := a.(*dns.NS); ok {
				res = append(res, strings.Trim(t.Ns, "."))
			}
		}

	case TypeMX:
		for _, a := range in.Answer {
			if t, ok := a.(*dns.MX); ok {
				res = append(res, strings.Trim(t.Mx, "."))
			}
		}

	case TypeTXT:
		for _, a := range in.Answer {
			if t, ok := a.(*dns.TXT); ok {
				res = append(res, t.Txt...)
			}
		}

	case TypeSRV:
		for _, a := range in.Answer {
			if t, ok := a.(*dns.SRV); ok {
				res = append(res, strings.Trim(t.Target, "."))
			}
		}

	case TypeCNAME:
		for _, a := range in.Answer {
			if t, ok := a.(*dns.CNAME); ok {
				res = append(res, strings.Trim(t.Target, "."))
			}
		}

	case TypePTR:
		for _, a := range in.Answer {
			if t, ok := a.(*dns.PTR); ok {
				res = append(res, strings.Trim(t.Ptr, "."))
			}
		}
	}

	return res, nil
}

func getAddr(ip net.IP) string {
	if ip.To4() == nil {
		// IPv6
		return fmt.Sprintf("[%s]:%d", ip, port)
	}

	return fmt.Sprintf("%s:%d", ip, port)
}
