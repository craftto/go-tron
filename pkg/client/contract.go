package client

import (
	"fmt"
	"strconv"

	"github.com/craftto/go-tron/pkg/abi"
	"github.com/craftto/go-tron/pkg/address"
	"github.com/craftto/go-tron/pkg/common"
	"github.com/craftto/go-tron/pkg/keystore"
	"github.com/craftto/go-tron/pkg/proto/core"
	"github.com/craftto/go-tron/pkg/transaction"
)

type TokenAmount struct {
	TokenId string
	Amount  int64
}

func (g *GrpcClient) TriggerConstantContract(contractAddress, from, method string, param []byte) (*transaction.Transaction, error) {
	var err error
	fromDesc, _ := address.Hex2Address(address.ZeroAddress)

	if len(from) > 0 {
		fromDesc, err = address.Base58ToAddress(from)
		if err != nil {
			return nil, err
		}
	}
	contractDesc, err := address.Base58ToAddress(contractAddress)
	if err != nil {
		return nil, err
	}

	signature := abi.MethodSignature(method)
	data := append(signature, param...)

	ct := &core.TriggerSmartContract{
		OwnerAddress:    fromDesc.Bytes(),
		ContractAddress: contractDesc.Bytes(),
		Data:            data,
	}

	ctx, cancel := g.GetContext()
	defer cancel()

	result, err := g.Client.TriggerConstantContract(ctx, ct)
	if err != nil {
		return nil, err
	}

	return &transaction.Transaction{
		TransactionHash: common.Bytes2Hex(result.Txid),
		Transaction:     result,
		Result:          result.Result,
	}, nil
}

func (g *GrpcClient) TriggerContract(ks *keystore.Keystore, contractAddress, method string, paramData []byte, feeLimit, amount int64, tokenAmount *TokenAmount) (*transaction.Transaction, error) {
	contractDesc, err := address.Base58ToAddress(contractAddress)
	if err != nil {
		return nil, err
	}

	signature := abi.MethodSignature(method)
	data := append(signature, paramData...)

	ct := &core.TriggerSmartContract{
		OwnerAddress:    ks.Address.Bytes(),
		ContractAddress: contractDesc.Bytes(),
		Data:            data,
		CallValue:       amount,
	}

	if tokenAmount != nil {
		ct.CallTokenValue = tokenAmount.Amount
		ct.TokenId, err = strconv.ParseInt(tokenAmount.TokenId, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return g.triggerContract(ks, ct, feeLimit)
}

func (g *GrpcClient) triggerContract(ks *keystore.Keystore, ct *core.TriggerSmartContract, feeLimit int64) (*transaction.Transaction, error) {
	ctx, cancel := g.GetContext()
	defer cancel()

	tx, err := g.Client.TriggerContract(ctx, ct)
	if err != nil {
		return nil, err
	}

	if tx.Result.Code > 0 {
		return nil, fmt.Errorf("%s", string(tx.Result.Message))
	}

	if feeLimit > 0 {
		tx.Transaction.RawData.FeeLimit = feeLimit
		if err := transaction.UpdateTxHash(tx); err != nil {
			return nil, err
		}
	}

	signedTx, err := ks.SignTx(tx.Transaction)
	if err != nil {
		return nil, err
	}

	result, err := g.Broadcast(signedTx)
	if err != nil {
		return nil, err
	}

	return &transaction.Transaction{
		TransactionHash: common.Bytes2Hex(tx.GetTxid()),
		Transaction:     tx,
		Result:          result,
	}, nil
}

// GetContractABI return smartContract
func (g *GrpcClient) GetContractABI(contractAddress string) (*core.SmartContract_ABI, error) {
	var err error

	contractDesc, err := address.Base58ToAddress(contractAddress)
	if err != nil {
		return nil, err
	}

	ctx, cancel := g.GetContext()
	defer cancel()

	sm, err := g.Client.GetContract(ctx, GetMessageBytes(contractDesc))
	if err != nil {
		return nil, err
	}
	if sm == nil {
		return nil, fmt.Errorf("invalid contract abi")
	}

	return sm.Abi, nil
}
