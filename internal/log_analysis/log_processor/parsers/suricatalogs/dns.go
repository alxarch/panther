package suricatalogs

/**
 * Panther is a Cloud-Native SIEM for the Modern Security Team.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"github.com/aws/aws-sdk-go/aws"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/numerics"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/timestamp"
)

const (
	TypeDNS = `Suricata.DNS`
	DNSDesc = `Suricata parser for the DNS event type in the EVE JSON output.
Reference: https://suricata.readthedocs.io/en/suricata-5.0.2/output/eve/eve-json-output.html#dns`
)

func init() {
	parsers.MustRegister(parsers.LogType{
		Name:        TypeDNS,
		Description: DNSDesc,
		Schema: struct {
			DNS
			parsers.PantherLog
		}{},
		NewParser: NewDNSParser,
	})
}

//nolint:lll
type DNS struct {
	CommunityID  *string                      `json:"community_id,omitempty" description:"Suricata DNS CommunityID"`
	DNS          *DNSDetails                  `json:"dns" validate:"required,dive" description:"Suricata DNS DNS"`
	DestIP       *string                      `json:"dest_ip" validate:"required" description:"Suricata DNS DestIP"`
	DestPort     *uint16                      `json:"dest_port,omitempty" description:"Suricata DNS DestPort"`
	EventType    *string                      `json:"event_type" validate:"required,eq=dns" description:"Suricata DNS EventType"`
	FlowID       *int                         `json:"flow_id,omitempty" description:"Suricata DNS FlowID"`
	PcapCnt      *int                         `json:"pcap_cnt,omitempty" description:"Suricata DNS PcapCnt"`
	PcapFilename *string                      `json:"pcap_filename,omitempty" description:"Suricata DNS PcapFilename"`
	Proto        *numerics.Integer            `json:"proto" validate:"required" description:"Suricata DNS Proto"`
	SrcIP        *string                      `json:"src_ip" validate:"required" description:"Suricata DNS SrcIP"`
	SrcPort      *uint16                      `json:"src_port,omitempty" description:"Suricata DNS SrcPort"`
	Timestamp    *timestamp.SuricataTimestamp `json:"timestamp" validate:"required" description:"Suricata DNS Timestamp"`
	Vlan         []int                        `json:"vlan,omitempty" description:"Suricata DNS Vlan"`
}

var _ parsers.PantherEventer = (*DNS)(nil)

func (event *DNS) PantherEvent() *parsers.PantherEvent {
	e := parsers.NewEvent(TypeDNS, event.Timestamp.UTC(),
		parsers.IPAddressP(event.SrcIP),
		parsers.IPAddressP(event.DestIP),
	)
	if event.DNS != nil {
		e.Extend(
			parsers.DomainNameP(event.DNS.Rrname),
			parsers.IPAddressP(event.DNS.RData),
		)
		for _, answer := range event.DNS.Answers {
			switch aws.StringValue(answer.Rrtype) {
			case "A", "AAAA":
				e.Extend(
					parsers.IPAddressP(answer.Rdata),
					parsers.DomainNameP(answer.Rrname),
				)
			case "CNAME", "MX":
				e.Extend(
					parsers.DomainNameP(answer.Rrname),
					parsers.DomainNameP(answer.Rdata),
				)
			case "PTR":
				e.Insert(parsers.DomainNameP(answer.Rdata))
			case "TXT":
				e.Insert(parsers.DomainNameP(answer.Rrname))
			}
		}
		if event.DNS.Grouped != nil {
			for _, aRecord := range event.DNS.Grouped.A {
				e.Insert(parsers.IPAddress(aRecord))
			}
			for _, aaaaRecord := range event.DNS.Grouped.Aaaa {
				e.Insert(parsers.IPAddress(aaaaRecord))
			}
			for _, cNameRecord := range event.DNS.Grouped.Cname {
				e.Insert(parsers.DomainName(cNameRecord))
			}
			for _, mxRecord := range event.DNS.Grouped.Mx {
				e.Insert(parsers.DomainName(mxRecord))
			}
		}
	}
	return e
}

//nolint:lll
type DNSDetails struct {
	Aa          *bool                   `json:"aa,omitempty" description:"Suricata DNSDetails Aa"`
	Answers     []DNSDetailsAnswers     `json:"answers,omitempty" validate:"omitempty,dive" description:"Suricata DNSDetails Answers"`
	Authorities []DNSDetailsAuthorities `json:"authorities,omitempty" validate:"omitempty,dive" description:"Suricata DNSDetails Authorities"`
	Flags       *string                 `json:"flags,omitempty" description:"Suricata DNSDetails Flags"`
	Grouped     *DNSDetailsGrouped      `json:"grouped,omitempty" validate:"omitempty,dive" description:"Suricata DNSDetails Grouped"`
	ID          *int                    `json:"id,omitempty" description:"Suricata DNSDetails ID"`
	Qr          *bool                   `json:"qr,omitempty" description:"Suricata DNSDetails Qr"`
	Ra          *bool                   `json:"ra,omitempty" description:"Suricata DNSDetails Ra"`
	Rcode       *string                 `json:"rcode,omitempty" description:"Suricata DNSDetails Rcode"`
	Rd          *bool                   `json:"rd,omitempty" description:"Suricata DNSDetails Rd"`
	Rrname      *string                 `json:"rrname,omitempty" description:"Suricata DNSDetails Rrname"`
	RData       *string                 `json:"rdata,omitempty" description:"Suricata DNSDetails RData"`
	Rrtype      *string                 `json:"rrtype,omitempty" description:"Suricata DNSDetails Rrtype"`
	TTL         *int                    `json:"ttl,omitempty" description:"Suricata DNSDetails TTL"`
	TxID        *int                    `json:"tx_id,omitempty" description:"Suricata DNSDetails TxID"`
	Type        *string                 `json:"type,omitempty" description:"Suricata DNSDetails Type"`
	Version     *int                    `json:"version,omitempty" description:"Suricata DNSDetails Version"`
}

//nolint:lll
type DNSDetailsAnswers struct {
	Rdata  *string `json:"rdata,omitempty" description:"Suricata DNSDetailsAnswers Rdata"`
	Rrname *string `json:"rrname,omitempty" description:"Suricata DNSDetailsAnswers Rrname"`
	Rrtype *string `json:"rrtype,omitempty" description:"Suricata DNSDetailsAnswers Rrtype"`
	TTL    *int    `json:"ttl,omitempty" description:"Suricata DNSDetailsAnswers TTL"`
}

//nolint:lll
type DNSDetailsGrouped struct {
	A     []string `json:"A,omitempty" description:"Suricata DNSDetailsGrouped A"`
	Aaaa  []string `json:"AAAA,omitempty" description:"Suricata DNSDetailsGrouped Aaaa"`
	Cname []string `json:"CNAME,omitempty" description:"Suricata DNSDetailsGrouped Cname"`
	Mx    []string `json:"MX,omitempty" description:"Suricata DNSDetailsGrouped Mx"`
	Ptr   []string `json:"PTR,omitempty" description:"Suricata DNSDetailsGrouped Ptr"`
	Txt   []string `json:"TXT,omitempty" description:"Suricata DNSDetailsGrouped Txt"`
}

//nolint:lll
type DNSDetailsAuthorities struct {
	Rrname *string `json:"rrname,omitempty" description:"Suricata DNSDetailsAuthorities Rrname"`
	Rrtype *string `json:"rrtype,omitempty" description:"Suricata DNSDetailsAuthorities Rrtype"`
	TTL    *int    `json:"ttl,omitempty" description:"Suricata DNSDetailsAuthorities TTL"`
}

// DNSParser parses Suricata DNS alerts in the JSON format
type DNSParser struct{}

var _ parsers.Parser = (*DNSParser)(nil)

func NewDNSParser() parsers.Parser {
	return &DNSParser{}
}

// Parse returns the parsed events or nil if parsing failed
func (p *DNSParser) Parse(log string) ([]*parsers.PantherLogJSON, error) {
	event := &DNS{}
	return parsers.QuickParseJSON(event, log)
}
