package utils

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/magiconair/properties"
	"io/ioutil"
	"log"
	"math/big"
)

func signAndSendTx(value big.Int, toAddressString string, data []byte, gasLimit int64) string {
	props := properties.MustLoadFile("default.properties", properties.UTF8)
	privateKey := props.GetString("private.key", "0000")
	privateKeyObj, err := crypto.HexToECDSA(privateKey)

	if err != nil {
		log.Fatal(err)
	}
	client, err := ethclient.Dial(BSCRPCNode)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKeyObj.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPriceGwei := props.GetInt64("gas.price", 0)
	gasPrice := big.NewInt(gasPriceGwei)

	toAddress := common.HexToAddress(toAddressString)
	tx := types.NewTransaction(nonce, toAddress, &value, uint64(gasLimit), gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		// TODO: implement fallback system
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyObj)
	if err != nil {
		log.Fatal(err)
	}

	binarySignedTx, err := signedTx.MarshalBinary()
	rawTx := hex.EncodeToString(binarySignedTx)

	var respTx *types.Transaction

	rawTxBytes, err := hex.DecodeString(rawTx)
	rlp.DecodeBytes(rawTxBytes, &respTx)

	err = client.SendTransaction(context.Background(), respTx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("tx sent: " + respTx.Hash().Hex())
	return respTx.Hash().Hex()
}

func readOnlyTx(toAddressString string, data []byte) []byte {
	client, err := ethclient.Dial(BSCRPCNode)
	if err != nil {
		log.Fatal(err)
	}
	addr := common.HexToAddress(toAddressString)

	msg := ethereum.CallMsg{To: &addr, Data: data}

	resp, err := client.CallContract(context.Background(), msg, nil)
	body, err := ioutil.ReadAll(bytes.NewReader(resp))
	return body
}
