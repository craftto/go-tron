package keystore

import (
	"crypto/ecdsa"
	"crypto/sha256"

	"github.com/craftto/go-tron/pkg/address"
	"github.com/craftto/go-tron/pkg/proto/core"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/protobuf/proto"
)

type Keystore struct {
	Address    address.Address
	privateKey ecdsa.PrivateKey
}

func ImportFromPrivateKey(privateKey string) (*Keystore, error) {
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	return &Keystore{
		Address:    address.PubkeyToAddress(privateKeyECDSA.PublicKey),
		privateKey: *privateKeyECDSA,
	}, nil
}

func (ks *Keystore) SignTx(tx *core.Transaction) (*core.Transaction, error) {
	rawData, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		return nil, err
	}

	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	signature, err := crypto.Sign(hash, &ks.privateKey)
	if err != nil {
		return nil, err
	}

	tx.Signature = append(tx.Signature, signature)

	return tx, nil
}
