package parser

import "github.com/kunal768/trustwallet/clients"

type Transaction struct {
	Address  string                `json:"address"`
	Inbound  []clients.Transaction `json:"inbound"`
	Outbound []clients.Transaction `json:"outbound"`
}
