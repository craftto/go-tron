package client

import (
	"bytes"
	"fmt"

	"github.com/craftto/go-tron/pkg/common"
	"github.com/craftto/go-tron/pkg/proto/api"
	"github.com/craftto/go-tron/pkg/proto/core"
	"google.golang.org/protobuf/proto"
)

// TotalTransaction return total transciton in network
func (g *GrpcClient) TotalTransaction() (*api.NumberMessage, error) {
	ctx, cancel := g.GetContext()
	defer cancel()

	return g.Client.TotalTransaction(ctx, new(api.EmptyMessage))
}

//GetTransactionByID returns transaction details by ID
func (g *GrpcClient) GetTransactionByID(txHash string) (*core.Transaction, error) {
	transactionID := new(api.BytesMessage)
	var err error

	transactionID.Value, err = common.Hex2Bytes(txHash)
	if err != nil {
		return nil, fmt.Errorf("get transaction by id error: %v", err)
	}

	ctx, cancel := g.GetContext()
	defer cancel()

	tx, err := g.Client.GetTransactionById(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	if size := proto.Size(tx); size == 0 {
		return nil, fmt.Errorf("transaction info not found")
	}

	return tx, nil
}

//GetTransactionInfoByID returns transaction receipt by ID
func (g *GrpcClient) GetTransactionInfoByID(txHash string) (*core.TransactionInfo, error) {
	transactionID := new(api.BytesMessage)
	var err error

	transactionID.Value, err = common.Hex2Bytes(txHash)
	if err != nil {
		return nil, fmt.Errorf("get transaction by id error: %v", err)
	}

	ctx, cancel := g.GetContext()
	defer cancel()

	tx, err := g.Client.GetTransactionInfoById(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(tx.Id, transactionID.Value) {
		return nil, fmt.Errorf("transaction info not found")

	}

	return tx, nil
}
