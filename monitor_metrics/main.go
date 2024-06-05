package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error"`
	ID      int             `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// URL for blockchain network
const (
	// url = "http://127.0.0.1:8545/" // Use it when running main file without docker
	url = "http://anvil-node:8545" // Use it when using docker
)

func sendRPCRequest(method string, params interface{}) (*JSONRPCResponse, error) {
	reqBody := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result JSONRPCResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if result.Error != nil {
		return nil, fmt.Errorf("RPC Error: %d %s", result.Error.Code, result.Error.Message)
	}

	return &result, nil
}

func main() {
	for {
		fmt.Println("\n*** Monitoring the block time ***")
		block_time, err := monitorBlockTime()
		if err != nil {
			log.Fatalf("Failed to retrieve block time: %v", err)
		}
		fmt.Println(block_time)

		fmt.Println("\n*** Monitoring the transaction throughput ***")
		tx_count, err := monitorTransactionThroughput()
		if err != nil {
			log.Fatalf("Failed to retrieve block time: %v", err)
		}
		fmt.Println(tx_count)

		fmt.Println("\n*** Monitoring the pending transaction ***")
		pending, queued, err := monitorTransactionPoolStatus()
		if err != nil {
			log.Fatalf("Failed to retrieve block time: %v", err)
		}
		fmt.Println("Pending transactions: ", pending)
		fmt.Println("Queued transactions: ", queued)

		// Method not found at local
		// fmt.Println("\n*** Monitoring the pending transaction ***")
		// getPeerCount()

		fmt.Println("\n*** Monitoring the transaction reciept ***")
		tx_details, err := getTransactionReciept("0x4cd2c562d4fbf475a549daad43b64135d5b4ec62c0eef57305ff809a3fd790c7")
		if err != nil {
			log.Fatalf("Failed to retrieve transaction details: %v", err)
		}
		fmt.Println("Transaction hash: ", tx_details[0])
		fmt.Println("Sender: ", tx_details[1])
		fmt.Println("Receiver: ", tx_details[2])
		fmt.Println("Status: ", tx_details[3])

		fmt.Println("\n*** Monitoring the transaction count in a block of given address ***")
		tx_count_of_addr, err := getTransactionCountOfContractInBlock("0xef11D1c2aA48826D4c41e54ab82D1Ff5Ad8A64Ca", "latest")
		if err != nil {
			log.Fatalf("Failed to retrieve transaction count od given address: %v", err)
		}
		fmt.Println("Transaction counts: ", tx_count_of_addr)

		time.Sleep(10 * time.Second)
	}
}

func monitorBlockTime() (string, error) {
	fmt.Println("calling rpc")
	result, err := sendRPCRequest("eth_blockNumber", []interface{}{})
	if err != nil {
		log.Fatalf("Failed to retrieve block number: %v", err)
	}

	var blockNumber string
	err = json.Unmarshal(result.Result, &blockNumber)
	if err != nil {
		log.Fatalf("Failed to parse block number: %v", err)
	}

	result, err = sendRPCRequest("eth_getBlockByNumber", []interface{}{blockNumber, false})
	if err != nil {
		log.Fatalf("Failed to retrieve block: %v", err)
	}

	var block struct {
		Timestamp string `json:"timestamp"`
	}
	err = json.Unmarshal(result.Result, &block)
	if err != nil {
		log.Fatalf("Failed to parse block: %v", err)
	}

	return block.Timestamp, nil
}

func monitorTransactionThroughput() (int, error) {
	result, err := sendRPCRequest("eth_blockNumber", []interface{}{})
	if err != nil {
		log.Fatalf("Failed to retrieve block number: %v", err)
	}

	var blockNumber string
	err = json.Unmarshal(result.Result, &blockNumber)
	if err != nil {
		log.Fatalf("Failed to parse block number: %v", err)
	}

	result, err = sendRPCRequest("eth_getBlockByNumber", []interface{}{blockNumber, true})
	if err != nil {
		log.Fatalf("Failed to retrieve block: %v", err)
	}

	var block struct {
		Transactions []interface{} `json:"transactions"`
	}
	err = json.Unmarshal(result.Result, &block)
	if err != nil {
		log.Fatalf("Failed to parse block: %v", err)
	}

	txCount := len(block.Transactions)

	return txCount, nil
}

func monitorTransactionPoolStatus() (string, string, error) {
	result, err := sendRPCRequest("txpool_status", []interface{}{})
	if err != nil {
		log.Fatalf("Failed to get monitor transaction pool status: %v", err)
	}

	var pool struct {
		PendingTx string `json:"pending"`
		QueuedTx  string `json:"queued"`
	}
	err = json.Unmarshal(result.Result, &pool)
	if err != nil {
		log.Fatalf("Failed to parse pool: %v", err)
	}
	return pool.PendingTx, pool.QueuedTx, nil
}

/*  ----- DAPP Metrics ------ */
func getTransactionReciept(t_id string) ([]string, error) {
	tx_res, err := sendRPCRequest("eth_getTransactionReceipt", []interface{}{t_id})
	if err != nil {
		log.Fatalf("Failed to retrieve transaction reciept: %v", err)
	}

	var tx_details struct {
		Tx_id  string `json:"transactionHash"`
		Sender string `json:"from"`
		To     string `json:"to"`
		Status string `json:"status"`
	}
	err = json.Unmarshal(tx_res.Result, &tx_details)
	if err != nil {
		log.Fatalf("Failed to retrieve transaction details: %v", err)
	}

	res := []string{tx_details.Tx_id, tx_details.Sender, tx_details.To, tx_details.Status}
	return res, nil
}

func getTransactionCountOfContractInBlock(address string, blockNo string) (string, error) {
	tx_count, err := sendRPCRequest("eth_getTransactionCount", []interface{}{address, blockNo})
	if err != nil {
		log.Fatalf("Failed to retrieve transaction count of given address: %v", err)
	}
	return string(tx_count.Result), nil
}
