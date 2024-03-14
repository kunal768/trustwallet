package parser

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/kunal768/trustwallet/clients"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock(ctx context.Context) (*big.Int, error)

	// add address to observer
	Subscribe(ctx context.Context, address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(ctx context.Context, address string) (Transaction, error)

	// update subscriber txns
	UpdateSubscriberTxns(ctx context.Context)
}

type service struct {
	client      clients.EthereumClient
	subscribers map[string]Transaction
}

func NewParserService(ethClient clients.EthereumClient) Parser {
	return &service{
		client:      ethClient,
		subscribers: make(map[string]Transaction),
	}
}

func (svc service) GetCurrentBlock(ctx context.Context) (*big.Int, error) {
	block, _, err := svc.client.EthGetBlockNumber(ctx)
	if err != nil {
		return &big.Int{}, err
	}
	return block, nil
}

func (svc service) Subscribe(ctx context.Context, address string) bool {
	if _, ok := svc.subscribers[address]; ok {
		// log error
		return false
	}
	svc.subscribers[address] = Transaction{Address: address}
	return true
}

func (svc service) getAllLatestTransactions(ctx context.Context) ([]clients.Transaction, error) {
	_, blockHex, err := svc.client.EthGetBlockNumber(ctx)
	if err != nil {
		return []clients.Transaction{}, err
	}

	allTxns, err := svc.client.GetBlockTransactions(blockHex)
	if err != nil {
		return []clients.Transaction{}, err
	}

	return allTxns, nil
}

func (svc service) GetTransactions(ctx context.Context, address string) (Transaction, error) {
	allTxns, err := svc.getAllLatestTransactions(ctx)
	if err != nil {
		return Transaction{}, err
	}

	res := Transaction{Address: address}

	for _, txn := range allTxns {
		if strings.EqualFold(address, txn.From) {
			res.Outbound = append(res.Outbound, txn)
		} else if strings.EqualFold(address, txn.To) {
			res.Inbound = append(res.Inbound, txn)
		}
	}

	return res, nil
}

func (svc service) UpdateSubscriberTxns(ctx context.Context) {
	if len(svc.subscribers) == 0 {
		return
	}

	allTxns, err := svc.getAllLatestTransactions(ctx)
	if err != nil {
		return
	}

	for _, txn := range allTxns {

		if from, ok := svc.subscribers[txn.From]; ok {

			// Then we modify the copy
			from.Outbound = append(from.Outbound, txn)

			// Then we reassign map entry
			svc.subscribers[txn.From] = Transaction{Address: from.Address, Inbound: from.Inbound, Outbound: from.Outbound}
		}

		if to, ok := svc.subscribers[txn.To]; ok {
			// Then we modify the copy
			to.Inbound = append(to.Inbound, txn)

			// Then we reassign map entry
			svc.subscribers[txn.To] = Transaction{Address: to.Address, Inbound: to.Inbound, Outbound: to.Outbound}
		}

	}

	// printing in order to see on console can convert to API
	fmt.Println(svc.subscribers)

}
