package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net/http"

	"github.com/kunal768/trustwallet/common"
)

type ethclient struct {
	client *http.Client
}

type EthereumClient interface {
	EthGetBlockNumber(ctx context.Context) (*big.Int, string, error)
	GetBlockTransactions(blockNum string) ([]Transaction, error)
}

func ClientService(client *http.Client) EthereumClient {
	return &ethclient{
		client: client,
	}
}

func sendRPCRequest(method string, params interface{}, client http.Client) (*EthResponse, error) {
	body, err := json.Marshal(EthRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      rand.Intn(100),
	})

	if err != nil {
		return &EthResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://cloudflare-eth.com", bytes.NewReader(body))
	if err != nil {
		return &EthResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return &EthResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &EthResponse{}, err
	}

	var response EthResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return &EthResponse{}, err
	}
	if response.Error != nil {
		return &EthResponse{}, fmt.Errorf("error from Ethereum node: %v", response.Error)
	}

	return &response, nil
}

func (c ethclient) EthGetBlockNumber(ctx context.Context) (*big.Int, string, error) {
	resp, err := sendRPCRequest("eth_blockNumber", nil, *c.client)
	if err != nil {
		return &big.Int{}, "", err
	}

	hexString, ok := resp.Result.(string)
	if !ok {
		return &big.Int{}, "", err
	}
	res, err := common.HexStringToBigInt(hexString)
	if err != nil {
		return &big.Int{}, "", err
	}
	return res, hexString, nil
}

func (c ethclient) GetBlockTransactions(blockNum string) ([]Transaction, error) {

	txns := []Transaction{}

	response, err := sendRPCRequest("eth_getBlockByNumber", []interface{}{blockNum, true}, *c.client)
	if err != nil {
		return txns, err
	}

	res, ok1 := response.Result.(map[string]interface{})
	transactions, ok2 := res["transactions"].([]interface{})
	if !ok1 || !ok2 {
		return txns, err
	}

	for _, txn := range transactions {
		m, ok := txn.(map[string]interface{})
		if ok {
			newTxn := Transaction{
				Hash:  m["hash"].(string),
				From:  m["from"].(string),
				To:    m["to"].(string),
				Value: m["value"].(string),
			}
			txns = append(txns, newTxn)
		}
	}

	return txns, nil
}
