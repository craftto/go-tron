package client

import (
	"fmt"

	"github.com/craftto/go-tron/pkg/common"
	"github.com/craftto/go-tron/pkg/proto/api"
	"github.com/craftto/go-tron/pkg/proto/core"
	"google.golang.org/protobuf/proto"
)

// ProposalsList return all network proposals
func (g *GrpcClient) ProposalsList() (*api.ProposalList, error) {
	ctx, cancel := g.GetContext()
	defer cancel()

	return g.Client.ListProposals(ctx, new(api.EmptyMessage))
}

// ProposalCreate create proposal based on parameter list
func (g *GrpcClient) ProposalCreate(from string, parameters map[int64]int64) (*api.TransactionExtention, error) {
	var err error

	contract := &core.ProposalCreateContract{
		Parameters: parameters,
	}

	if contract.OwnerAddress, err = common.DecodeBase58(from); err != nil {
		return nil, err
	}

	ctx, cancel := g.GetContext()
	defer cancel()

	tx, err := g.Client.ProposalCreate(ctx, contract)
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

// ProposalApprove change URL info
func (g *GrpcClient) ProposalApprove(from string, id int64, confirm bool) (*api.TransactionExtention, error) {
	var err error

	contract := &core.ProposalApproveContract{
		ProposalId:    id,
		IsAddApproval: confirm,
	}

	if contract.OwnerAddress, err = common.DecodeBase58(from); err != nil {
		return nil, err
	}

	ctx, cancel := g.GetContext()
	defer cancel()

	tx, err := g.Client.ProposalApprove(ctx, contract)
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

func (g *GrpcClient) ProposalWithdraw(from string, id int64) (*api.TransactionExtention, error) {
	var err error

	contract := &core.ProposalDeleteContract{
		ProposalId: id,
	}

	if contract.OwnerAddress, err = common.DecodeBase58(from); err != nil {
		return nil, err
	}

	ctx, cancel := g.GetContext()
	defer cancel()

	tx, err := g.Client.ProposalDelete(ctx, contract)
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
