package client

import (
	"fmt"

	"github.com/craftto/go-tron/pkg/common"
	"github.com/craftto/go-tron/pkg/proto/api"
	"github.com/craftto/go-tron/pkg/proto/core"
	"google.golang.org/protobuf/proto"
)

// FreezeBalance from base58 address
func (g *GrpcClient) FreezeBalance(from, delegateTo string, resource core.ResourceCode, frozenBalance int64) (*api.TransactionExtention, error) {
	var err error

	contract := &core.FreezeBalanceContract{}
	if contract.OwnerAddress, err = common.DecodeBase58(from); err != nil {
		return nil, err
	}

	contract.FrozenBalance = frozenBalance
	contract.FrozenDuration = 3 // Tron Only allows 3 days freeze

	if len(delegateTo) > 0 {
		if contract.ReceiverAddress, err = common.DecodeBase58(delegateTo); err != nil {
			return nil, err
		}

	}
	contract.Resource = resource

	ctx, cancel := g.GetContext()
	defer cancel()

	tx, err := g.Client.FreezeBalance2(ctx, contract)
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

// UnfreezeBalance from base58 address
func (g *GrpcClient) UnfreezeBalance(from, delegateTo string, resource core.ResourceCode) (*api.TransactionExtention, error) {
	var err error
	contract := &core.UnfreezeBalanceContract{}

	if contract.OwnerAddress, err = common.DecodeBase58(from); err != nil {
		return nil, err
	}

	if len(delegateTo) > 0 {
		if contract.ReceiverAddress, err = common.DecodeBase58(delegateTo); err != nil {
			return nil, err
		}

	}
	contract.Resource = resource

	ctx, cancel := g.GetContext()
	defer cancel()

	tx, err := g.Client.UnfreezeBalance2(ctx, contract)
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
