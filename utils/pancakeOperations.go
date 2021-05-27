package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func SwapApprovedTokenToBUSD(amountIn big.Int, minAmountOut big.Int, userAddr string, epochTime big.Int, token string) string {
	timeout := hexutil.EncodeBig(new(big.Int).Add(&epochTime, new(big.Int).SetInt64(20000)))

	swapFnSignature := []byte("swapExactTokensForTokens(uint256,uint256,address[],address,uint256)")
	fnHash := crypto.Keccak256Hash(swapFnSignature).Bytes()[:4]

	amountInDataBytes := make([]byte, 32)
	amountInDataBytes = amountIn.Bytes()
	amountInParam := hexutil.Encode(amountInDataBytes)

	minAmountOutBytes := make([]byte, 32)
	minAmountOutBytes = minAmountOut.Bytes()
	minAmountOutParam := hexutil.Encode(minAmountOutBytes)

	amountInParam = amountInParam[2:]
	minAmountOutParam = minAmountOutParam[2:]

	for len(amountInParam) < 64 {
		amountInParam = "0" + amountInParam
	}
	for len(minAmountOutParam) < 64 {
		minAmountOutParam = "0" + minAmountOutParam
	}

	timeoutBlockParam := timeout[2:]
	for len(timeoutBlockParam) < 64 {
		timeoutBlockParam = "0" + timeoutBlockParam
	}

	dataString := hexutil.Encode(fnHash) + amountInParam + minAmountOutParam +
		ParamWrap2 + TokenPadding + userAddr[2:] + timeoutBlockParam + ParamWrapArr + TokenPadding + token + WBNBTokenPadded + BUSDTokenPadded

	data := common.Hex2Bytes(dataString[2:])

	tx := signAndSendTx(*big.NewInt(0), PancakeRouterAddr, data, 350000)
	return tx
}
