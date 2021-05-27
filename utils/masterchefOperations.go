package utils

import (
	"encoding/binary"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func PendingTokens(fnSignature string, poolId int64, contractAddr string, userAddr string) big.Int {
	pendingFnSignature := []byte(fnSignature)
	functionHash := crypto.Keccak256Hash(pendingFnSignature).Bytes()[:4]

	poolDataBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(poolDataBytes, uint64(poolId))
	poolDataHex := hexutil.Encode(poolDataBytes)

	dataString := hexutil.Encode(functionHash) + IntPadding + poolDataHex[2:] + TokenPadding + userAddr[2:]
	dataBytes := common.Hex2Bytes(dataString[2:])

	pendingTokenCountBytes := readOnlyTx(contractAddr, dataBytes)
	pendingTokenCount := new(big.Int).SetBytes(pendingTokenCountBytes)
	return *pendingTokenCount
}

func Harvest(poolId int64, contractAddr string) string {
	harvestFnSignature := []byte("deposit(uint256,uint256)")
	methodHash := crypto.Keccak256Hash(harvestFnSignature).Bytes()[:4]

	poolDataBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(poolDataBytes, uint64(poolId))
	poolDataHex := hexutil.Encode(poolDataBytes)

	dataString := hexutil.Encode(methodHash) + IntPadding + poolDataHex[2:] + EthMinValue
	dataBytes := common.Hex2Bytes(dataString[2:])
	return signAndSendTx(*big.NewInt(0), contractAddr, dataBytes, 200000)
}

func Approve(tokenAddr string) string {
	approveFnSignature := []byte("approve(address,uint256)")
	methodHash := crypto.Keccak256Hash(approveFnSignature).Bytes()[:4]

	dataString := hexutil.Encode(methodHash) + PancakeRouterAddrPadded + EthMaxValue
	dataBytes := common.Hex2Bytes(dataString[2:])
	return signAndSendTx(*big.NewInt(0), tokenAddr, dataBytes, 50000)
}

func IsApproved(tokenAddr string, walletAddr string) bool {
	allowanceFnSignature := []byte("allowance(address,address)")
	functionHash := crypto.Keccak256Hash(allowanceFnSignature).Bytes()[:4]

	dataString := hexutil.Encode(functionHash) + TokenPadding + walletAddr[2:] + PancakeRouterAddrPadded
	dataBytes := common.Hex2Bytes(dataString[2:])
	allowanceBytes := readOnlyTx(tokenAddr, dataBytes)
	allowance := hexutil.Encode(allowanceBytes)
	return string(allowance) != "0x"+EthMinValue
}

func BalanceOf(tokenAddr string, walletAddr string) big.Int {
	allowanceFnSignature := []byte("balanceOf(address)")
	functionHash := crypto.Keccak256Hash(allowanceFnSignature).Bytes()[:4]

	dataString := hexutil.Encode(functionHash) + TokenPadding + walletAddr[2:]
	dataBytes := common.Hex2Bytes(dataString[2:])
	balanceBytes := readOnlyTx(tokenAddr, dataBytes)

	balance := new(big.Int).SetBytes(balanceBytes)
	return *balance
}
