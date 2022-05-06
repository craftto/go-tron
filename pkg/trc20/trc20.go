package trc20

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/craftto/go-tron/pkg/abi"
	"github.com/craftto/go-tron/pkg/address"
	"github.com/craftto/go-tron/pkg/client"
	"github.com/craftto/go-tron/pkg/common"
	"github.com/craftto/go-tron/pkg/contract"
	"github.com/craftto/go-tron/pkg/keystore"
	"github.com/craftto/go-tron/pkg/proto/api"
	"github.com/craftto/go-tron/pkg/proto/core"
	"github.com/craftto/go-tron/pkg/transaction"
)

const (
	methodName         = "0x06fdde03"
	methodSymbol       = "0x95d89b41"
	methodDecimals     = "0x313ce567"
	methodTotalSupply  = "0x18160ddd"
	methodBalanceOf    = "0x70a08231"
	methodAllowance    = "0xdd62ed3e"
	methodApprove      = "0x095ea7b3"
	methodTransfer     = "0xa9059cbb"
	methodTransferFrom = "0x23b872dd"
)

var (
	feeLimit int64 = 30_000_000
)

type TRC20 struct {
	ContractAddress address.Address
	*client.GrpcClient
}

func NewTrc20(g *client.GrpcClient, contractAddr string) (*TRC20, error) {
	addr, err := address.Base58ToAddress(contractAddr)
	if err != nil {
		return nil, err
	}

	return &TRC20{
		ContractAddress: addr,
		GrpcClient:      g,
	}, nil
}

func (t *TRC20) GetName() (string, error) {
	data, err := common.Hex2Bytes(methodName)
	if err != nil {
		return "", err
	}

	result, err := t.callConstant(data)
	if err != nil {
		return "", err
	}

	return contract.ParseString(result.GetConstantResult()[0])
}

func (t *TRC20) GetSymbol() (string, error) {
	data, err := common.Hex2Bytes(methodSymbol)
	if err != nil {
		return "", err
	}

	result, err := t.callConstant(data)
	if err != nil {
		return "", err
	}

	return contract.ParseString(result.GetConstantResult()[0])
}

func (t *TRC20) GetDecimals() (*big.Int, error) {
	data, err := common.Hex2Bytes(methodDecimals)
	if err != nil {
		return nil, err
	}

	result, err := t.callConstant(data)
	if err != nil {
		return nil, err
	}

	return contract.ParseInt(result.GetConstantResult()[0]), nil
}

func (t *TRC20) GetTotalSupply() (*big.Int, error) {
	data, err := common.Hex2Bytes(methodTotalSupply)
	if err != nil {
		return nil, err
	}

	result, err := t.callConstant(data)
	if err != nil {
		return nil, err
	}

	return contract.ParseInt(result.GetConstantResult()[0]), nil
}

func (t *TRC20) GetBalance(addr string) (*big.Int, error) {
	addrA, err := address.Base58ToAddress(addr)
	if err != nil {
		return nil, err
	}

	data := methodBalanceOf + abi.PaddedParam(addrA.Hex()[2:])
	dataBytes, err := common.Hex2Bytes(data)
	if err != nil {
		return nil, err
	}

	result, err := t.callConstant(dataBytes)
	if err != nil {
		return nil, err
	}

	return contract.ParseInt(result.GetConstantResult()[0]), nil
}

func (t *TRC20) GetAllowance(owner, spender string) (*big.Int, error) {
	ownerAddr, err := address.Base58ToAddress(owner)
	if err != nil {
		return nil, err
	}

	spwnderAddr, err := address.Base58ToAddress(spender)
	if err != nil {
		return nil, err
	}

	data := methodAllowance + abi.PaddedParam(ownerAddr.Hex()[2:]) + abi.PaddedParam(spwnderAddr.Hex()[2:])
	dataBytes, err := common.Hex2Bytes(data)
	if err != nil {
		return nil, err
	}

	result, err := t.callConstant(dataBytes)
	if err != nil {
		return nil, err
	}

	return contract.ParseInt(result.GetConstantResult()[0]), nil
}

func (t *TRC20) Approve(ks *keystore.Keystore, spender string, amount *big.Int) (*transaction.Transaction, error) {
	param, err := abi.GetParams([]abi.Param{
		{"address": spender},
		{"uint256": amount},
	})
	if err != nil {
		return nil, err
	}

	data, err := common.Hex2Bytes(methodApprove)
	if err != nil {
		return nil, err
	}

	data = append(data, param...)

	ct := &core.TriggerSmartContract{
		OwnerAddress:    ks.Address.Bytes(),
		ContractAddress: t.ContractAddress.Bytes(),
		Data:            data,
	}

	return t.call(ks, ct)
}

func (t *TRC20) Transfer(ks *keystore.Keystore, to string, amount *big.Int) (*transaction.Transaction, error) {
	param, err := abi.GetParams([]abi.Param{
		{"address": to},
		{"uint256": amount},
	})
	if err != nil {
		return nil, err
	}

	data, err := common.Hex2Bytes(methodTransfer)
	if err != nil {
		return nil, err
	}

	data = append(data, param...)

	ct := &core.TriggerSmartContract{
		OwnerAddress:    ks.Address.Bytes(),
		ContractAddress: t.ContractAddress.Bytes(),
		Data:            data,
	}

	return t.call(ks, ct)
}

func (t *TRC20) TransferFrom(ks *keystore.Keystore, from, to string, amount *big.Int) (*transaction.Transaction, error) {
	param, err := abi.GetParams([]abi.Param{
		{"address": from},
		{"address": to},
		{"uint256": amount},
	})
	if err != nil {
		return nil, err
	}

	data, err := common.Hex2Bytes(methodTransferFrom)
	if err != nil {
		return nil, err
	}

	data = append(data, param...)

	ct := &core.TriggerSmartContract{
		OwnerAddress:    ks.Address.Bytes(),
		ContractAddress: t.ContractAddress.Bytes(),
		Data:            data,
	}

	return t.call(ks, ct)
}

func (t *TRC20) Call(ks *keystore.Keystore, method string, params []byte) (*transaction.Transaction, error) {
	signature := abi.MethodSignature(method)
	data := append(signature, params...)

	ct := &core.TriggerSmartContract{
		OwnerAddress:    ks.Address.Bytes(),
		ContractAddress: t.ContractAddress.Bytes(),
		Data:            data,
	}

	return t.call(ks, ct)
}

func (t *TRC20) CallConstant(method string, params []byte) (string, error) {
	signature := abi.MethodSignature(method)
	data := append(signature, params...)

	result, err := t.callConstant(data)
	if err != nil {
		return "", err
	}

	return contract.ParseString(result.GetConstantResult()[0])
}

func (t *TRC20) callConstant(data []byte) (*api.TransactionExtention, error) {
	addr, _ := address.Hex2Address(address.ZeroAddress)
	ct := &core.TriggerSmartContract{
		OwnerAddress:    addr.Bytes(),
		ContractAddress: t.ContractAddress.Bytes(),
		Data:            data,
	}

	ctx, cancel := t.GrpcClient.GetContext()
	defer cancel()

	result, err := t.Client.TriggerConstantContract(ctx, ct)
	if err != nil {
		return nil, err
	}

	if result.Result.Code > 0 {
		return nil, errors.New(string(result.Result.Message))
	}

	return result, nil
}

func (t *TRC20) call(ks *keystore.Keystore, ct *core.TriggerSmartContract) (*transaction.Transaction, error) {
	ctx, cancel := t.GetContext()
	defer cancel()

	tx, err := t.Client.TriggerContract(ctx, ct)
	if err != nil {
		return nil, err
	}

	if tx.Result.Code > 0 {
		return nil, fmt.Errorf("%s", string(tx.Result.Message))
	}

	tx.Transaction.RawData.FeeLimit = feeLimit
	if err := transaction.UpdateTxHash(tx); err != nil {
		return nil, err
	}

	signedTx, err := ks.SignTx(tx.Transaction)
	if err != nil {
		return nil, err
	}

	result, err := t.Broadcast(signedTx)
	if err != nil {
		return nil, err
	}

	return &transaction.Transaction{
		TransactionHash: common.Bytes2Hex(tx.GetTxid()),
		Result:          result,
	}, nil
}
