package address

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/craftto/go-tron/pkg/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength = 21
	// AddressLengthBase58 is the expected length of the address in base58format
	AddressLengthBase58 = 34
	// TronBytePrefix is the hex prefix to address
	TronBytePrefix = byte(0x41)
)

var (
	ZeroAddress = "410000000000000000000000000000000000000000"
)

// Address represents the 21 byte address of an Tron account.
type Address []byte

// String implements fmt.Stringer.
func (a Address) String() string {
	return a.Base58()
}

// Bytes get bytes from address
func (a Address) Bytes() []byte {
	return a[:]
}

// Hex get bytes from address in string
func (a Address) Hex() string {
	return common.Bytes2Hex(a[:])
}

// Base58 get base58 encoded string
func (a Address) Base58() string {
	if a[0] == 0 {
		return new(big.Int).SetBytes(a.Bytes()).String()
	}

	return common.EncodeBase58(a.Bytes())
}

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func Hex2Address(s string) (Address, error) {
	addr, err := common.Hex2Bytes(s)
	if err != nil {
		return nil, err
	}

	return addr, nil
}

// Base58ToAddress returns Address with byte values of s.
func Base58ToAddress(s string) (Address, error) {
	addr, err := common.DecodeBase58(s)
	if err != nil {
		return nil, err
	}

	return addr, nil
}

// PubkeyToAddress returns address from ecdsa public key
func PubkeyToAddress(p ecdsa.PublicKey) Address {
	address := crypto.PubkeyToAddress(p)

	addressTron := make([]byte, 0)
	addressTron = append(addressTron, TronBytePrefix)
	addressTron = append(addressTron, address.Bytes()...)
	return addressTron
}
