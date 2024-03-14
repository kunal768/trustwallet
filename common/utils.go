package common

import (
	"fmt"
	"math/big"
	"strings"
)

func HexStringToBigInt(hexStr string) (*big.Int, error) {
	// Remove '0x' if present
	hexStr = strings.TrimPrefix(hexStr, "0x")

	// Convert hexadecimal string to a big.Int
	bigint := new(big.Int)
	_, success := bigint.SetString(hexStr, 16)
	if !success {
		return nil, fmt.Errorf("error converting hex string to big.Int")
	}

	return bigint, nil
}
