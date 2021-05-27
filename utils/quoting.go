package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

func GetTokenUSDPrice(tokenAddr string) *big.Float {
	dataString := GetPairFnHash + BUSDTokenPadded + TokenPadding + tokenAddr[2:]
	data := common.Hex2Bytes(dataString)

	pairContractAddrBytes := readOnlyTx(PancakeFactoryAddr, data)
	pairContractAddr := hexutil.Encode(pairContractAddrBytes[12:])

	dataString = GetToken0FnHash
	data = common.Hex2Bytes(dataString)
	token0AddrBytes := readOnlyTx(pairContractAddr, data)
	token0Addr := hexutil.Encode(token0AddrBytes)

	dataString = GetReservesFnHash
	data = common.Hex2Bytes(dataString)

	pairReserves := readOnlyTx(pairContractAddr, data)
	BUSDReserveBytes := new([]byte)
	tokenReserveBytes := new([]byte)

	if token0Addr == "0x"+BUSDTokenPadded {
		*BUSDReserveBytes = pairReserves[:32]
		*tokenReserveBytes = pairReserves[32:64]
	} else {
		*tokenReserveBytes = pairReserves[:32]
		*BUSDReserveBytes = pairReserves[32:64]
	}

	BUSDReserve := new(big.Int)
	tokenReserve := new(big.Int)
	BUSDReserve.SetBytes(*BUSDReserveBytes)
	tokenReserve.SetBytes(*tokenReserveBytes)

	BUSDReserveFloat := new(big.Float).SetInt(BUSDReserve)
	tokenReserveFloat := new(big.Float).SetInt(tokenReserve)

	USDPrice := new(big.Float).Quo(BUSDReserveFloat, tokenReserveFloat)
	return USDPrice
}
