package abi

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/craftto/go-tron/pkg/address"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

type Param map[string]interface{}

// MethodSignature get the signature of a method
func MethodSignature(method string) []byte {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(method))
	b := hasher.Sum(nil)

	return b[:4]
}

func GetParams(param []Param) ([]byte, error) {
	values := make([]interface{}, 0)
	arguments := abi.Arguments{}

	for _, p := range param {
		if len(p) != 1 {
			return nil, fmt.Errorf("invalid param %+v", p)
		}
		for k, v := range p {
			ty, err := abi.NewType(k, "", nil)
			if err != nil {
				return nil, fmt.Errorf("invalid param %+v: %+v", p, err)
			}
			arguments = append(arguments,
				abi.Argument{
					Name:    "",
					Type:    ty,
					Indexed: false,
				},
			)

			if ty.T == abi.SliceTy || ty.T == abi.ArrayTy {
				if ty.Elem.T == abi.AddressTy {
					tmp := v.([]string)
					v = make([]common.Address, 0)
					for i := range tmp {
						addr, err := toAddress(tmp[i])
						if err != nil {
							return nil, err
						}
						v = append(v.([]common.Address), addr)
					}
				}

				if (ty.Elem.T == abi.IntTy || ty.Elem.T == abi.UintTy) &&
					ty.Elem.Size > 64 &&
					reflect.TypeOf(v).Elem().Kind() == reflect.String {
					tmp := make([]*big.Int, 0)
					for _, s := range v.([]string) {
						var value *big.Int
						// check for hex char
						if strings.HasPrefix(s, "0x") {
							value, _ = new(big.Int).SetString(s[2:], 16)
						} else {
							value, _ = new(big.Int).SetString(s, 10)
						}
						tmp = append(tmp, value)
					}
					v = tmp
				}
			}
			if ty.T == abi.AddressTy {
				if v, err = toAddress(v); err != nil {
					return nil, err
				}
			}
			if (ty.T == abi.IntTy || ty.T == abi.UintTy) && reflect.TypeOf(v).Kind() == reflect.String {
				v = toInt(ty, v)
			}

			if ty.T == abi.BytesTy || ty.T == abi.FixedBytesTy {
				var err error
				if v, err = toBytes(ty, v); err != nil {
					return nil, err
				}
			}

			values = append(values, v)
		}
	}
	// convert params to bytes
	return arguments.PackValues(values)
}

func PaddedParam(param string) string {
	return "0000000000000000000000000000000000000000000000000000000000000000"[len(param):] + param
}

func toAddress(v interface{}) (common.Address, error) {
	switch v.(type) {
	case string:
		addr, err := address.Base58ToAddress(v.(string))
		if err != nil {
			return common.Address{}, fmt.Errorf("invalid address %s: %+v", v.(string), err)
		}
		return common.BytesToAddress(addr.Bytes()[len(addr.Bytes())-20:]), nil
	}
	return common.Address{}, fmt.Errorf("invalid address %v", v)
}

func toInt(ty abi.Type, v interface{}) interface{} {
	if ty.T == abi.IntTy && ty.Size <= 64 {
		tmp, _ := strconv.ParseInt(v.(string), 10, ty.Size)
		switch ty.Size {
		case 8:
			v = int8(tmp)
		case 16:
			v = int16(tmp)
		case 32:
			v = int32(tmp)
		case 64:
			v = int64(tmp)
		}
	} else if ty.T == abi.UintTy && ty.Size <= 64 {
		tmp, _ := strconv.ParseUint(v.(string), 10, ty.Size)
		switch ty.Size {
		case 8:
			v = uint8(tmp)
		case 16:
			v = uint16(tmp)
		case 32:
			v = uint32(tmp)
		case 64:
			v = uint64(tmp)
		}
	} else {
		s := v.(string)
		// check for hex char
		if strings.HasPrefix(s, "0x") {
			v, _ = new(big.Int).SetString(s[2:], 16)
		} else {
			v, _ = new(big.Int).SetString(s, 10)
		}
	}
	return v
}

func toBytes(ty abi.Type, v interface{}) (interface{}, error) {
	// if string
	if data, ok := v.(string); ok {
		// convert from hex string
		dataBytes, err := hex.DecodeString(data)
		if err != nil {
			// try with base64
			dataBytes, err = base64.StdEncoding.DecodeString(data)
			if err != nil {
				return nil, err
			}
		}
		// if array and size == 0
		if ty.T == abi.BytesTy || ty.Size == 0 {
			return dataBytes, nil
		}
		if len(dataBytes) != ty.Size {
			return nil, fmt.Errorf("invalid size: %d/%d", ty.Size, len(dataBytes))
		}
		switch ty.Size {
		case 1:
			value := [1]byte{}
			copy(value[:], dataBytes[:1])
			return value, nil
		case 2:
			value := [2]byte{}
			copy(value[:], dataBytes[:2])
			return value, nil
		case 8:
			value := [8]byte{}
			copy(value[:], dataBytes[:8])
			return value, nil
		case 16:
			value := [16]byte{}
			copy(value[:], dataBytes[:16])
			return value, nil
		case 32:
			value := [32]byte{}
			copy(value[:], dataBytes[:32])
			return value, nil
		}
	}
	return v, nil
}
