package transaction

import (
	"crypto/sha256"

	"github.com/craftto/go-tron/pkg/proto/api"
	"google.golang.org/protobuf/proto"
)

type Transaction struct {
	TransactionHash string
	Transaction     *api.TransactionExtention
	Result          *api.Return
}

func UpdateTxHash(tx *api.TransactionExtention) error {
	rawData, err := proto.Marshal(tx.Transaction.GetRawData())
	if err != nil {
		return err
	}

	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	tx.Txid = hash

	return nil
}
