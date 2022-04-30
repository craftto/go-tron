package client

import (
	"fmt"

	"github.com/craftto/go-tron/pkg/common"
	"github.com/craftto/go-tron/pkg/proto/api"
	"github.com/craftto/go-tron/pkg/proto/core"
	"google.golang.org/protobuf/proto"
)

// Transfer from to base58 address
func (g *GrpcClient) Transfer(from, toAddress string, amount int64) (*api.TransactionExtention, error) {
	var err error

	contract := &core.TransferContract{}
	if contract.OwnerAddress, err = common.DecodeBase58(from); err != nil {
		return nil, err
	}
	if contract.ToAddress, err = common.DecodeBase58(toAddress); err != nil {
		return nil, err
	}
	contract.Amount = amount

	ctx, cancel := g.GetContext()
	defer cancel()

	tx, err := g.Client.CreateTransaction2(ctx, contract)
	if err != nil {
		return nil, err
	}

	if proto.Size(tx) == 0 {
		return nil, fmt.Errorf("bad transaction")
	}

	if tx.GetResult().GetCode() != 0 {
		return nil, fmt.Errorf("%s", tx.GetResult().GetMessage())
	}
	return tx, nil
}
