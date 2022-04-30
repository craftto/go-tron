package transaction

import (
	"github.com/craftto/go-tron/pkg/address"
	"github.com/craftto/go-tron/pkg/common"
	"github.com/craftto/go-tron/pkg/proto/core"
)

type TransactionReceipt struct {
	TxId                          string
	Result                        core.TransactionInfoCode
	Fee                           int64
	BlockNumber                   int64
	BlockTimeStamp                int64
	ComtractAddress               string
	Receipt                       *core.ResourceReceipt
	Logs                          []Log
	ResMessage                    []byte
	AssetIssueID                  string
	WithdrawAmount                int64
	UnfreezeAmount                int64
	InternalTransactions          []*core.InternalTransaction
	ExchangeReceivedAmount        int64
	ExchangeInjectAnotherAmount   int64
	ExchangeWithdrawAnotherAmount int64
	ExchangeId                    int64
	ShieldedTransactionFee        int64
}

type Log struct {
	Address address.Address
	Topics  []string
	Data    string
}

func GetTransactionReceipt(tx *core.TransactionInfo) (*TransactionReceipt, error) {
	receipt := &TransactionReceipt{
		TxId:                          common.Bytes2Hex(tx.GetId()),
		Fee:                           tx.Fee,
		BlockNumber:                   tx.BlockNumber,
		BlockTimeStamp:                tx.BlockTimeStamp,
		Result:                        tx.Result,
		Receipt:                       tx.Receipt,
		ResMessage:                    tx.ResMessage,
		AssetIssueID:                  tx.AssetIssueID,
		WithdrawAmount:                tx.WithdrawAmount,
		UnfreezeAmount:                tx.UnfreezeAmount,
		InternalTransactions:          tx.InternalTransactions,
		ExchangeReceivedAmount:        tx.ExchangeReceivedAmount,
		ExchangeInjectAnotherAmount:   tx.ExchangeInjectAnotherAmount,
		ExchangeWithdrawAnotherAmount: tx.ExchangeWithdrawAnotherAmount,
		ExchangeId:                    tx.ExchangeId,
		ShieldedTransactionFee:        tx.ShieldedTransactionFee,
	}

	if tx.GetContractAddress() != nil {
		receipt.ComtractAddress = address.Address(tx.GetContractAddress()).Base58()
	}

	if tx.GetLog() != nil {
		logs := make([]Log, 0, len(tx.Log))
		for _, log := range tx.Log {
			l := Log{
				Address: address.Address(log.Address),
				Topics:  make([]string, 0, len(log.Topics)),
				Data:    common.Bytes2Hex(log.Data),
			}

			for _, topic := range log.Topics {
				l.Topics = append(l.Topics, common.Bytes2Hex(topic))
			}

			logs = append(logs, l)
		}

		receipt.Logs = logs
	}

	return receipt, nil
}
