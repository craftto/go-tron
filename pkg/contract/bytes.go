package contract

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"unicode/utf8"

	"github.com/craftto/go-tron/pkg/address"
	"github.com/craftto/go-tron/pkg/common"
)

func ParseString(data []byte) (string, error) {
	dataStr := common.Bytes2Hex(data)

	if common.Has0xPrefix(dataStr) {
		dataStr = dataStr[2:]
	}

	if len(dataStr) > 128 {
		dataN := dataStr[64:128]
		if len(dataN) != 64 {
			return "", errors.New("Cannot parse data")
		}
		var n big.Int
		if _, ok := n.SetString(dataN, 16); !ok {
			return "", errors.New("Cannot parse data")
		}

		l := n.Uint64()
		if 2*int(l) <= len(dataStr)-128 {
			b, err := hex.DecodeString(dataStr[128 : 128+2*l])
			if err == nil {
				return string(b), nil
			}
		}
	} else if len(dataStr) == 64 {
		// allow string properties as 32 bytes of UTF-8 data
		b, err := hex.DecodeString(dataStr)
		if err == nil {
			i := bytes.Index(b, []byte{0})
			if i > 0 {
				b = b[:i]
			}
			if utf8.Valid(b) {
				return string(b), nil
			}
		}
	}
	return "", errors.New("Cannot parse data")
}

func ParseInt(data []byte) *big.Int {
	return new(big.Int).SetBytes(data)
}

func ParseAddress(data []byte) (address.Address, error) {
	if len(data) < 20 {
		return *new(address.Address), errors.New("Cannot parse data")
	}

	data = data[len(data)-20:]

	addr := make([]byte, 0)
	addr = append(addr, byte(0x41))
	addr = append(addr, data...)

	return address.Address(addr), nil
}

func ParseStringArray(data []byte) ([]string, error) {
	if len(data) < 64 {
		return nil, errors.New("Cannot parse data")
	}

	data = data[32:]
	size := ParseInt(data[:32])
	data = data[32:]

	arr := make([]string, 0, size.Int64())
	for i := 0; i < int(size.Int64()); i++ {
		index := 32 * i
		value, err := ParseString(data[index : index+32])
		if err != nil {
			return nil, errors.New("Cannot parse data")
		}

		arr = append(arr, value)
	}

	return arr, nil
}

func ParseIntArray(data []byte) ([]*big.Int, error) {
	if len(data) < 64 {
		return nil, errors.New("Cannot parse data")
	}

	data = data[32:]
	size := ParseInt(data[:32])
	data = data[32:]

	arr := make([]*big.Int, 0, size.Int64())
	for i := 0; i < int(size.Int64()); i++ {
		index := 32 * i
		arr = append(arr, ParseInt(data[index:index+32]))
	}

	return arr, nil
}

func ParseAddressArray(data []byte) ([]address.Address, error) {
	if len(data) < 64 {
		return nil, errors.New("Cannot parse data")
	}

	data = data[32:]
	size := ParseInt(data[:32])
	data = data[32:]

	arr := make([]address.Address, 0, size.Int64())
	for i := 0; i < int(size.Int64()); i++ {
		index := 32 * i
		value, err := ParseAddress(data[index : index+32])
		if err != nil {
			return nil, errors.New("Cannot parse data")
		}

		arr = append(arr, value)
	}

	return arr, nil
}
