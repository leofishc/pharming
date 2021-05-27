package main

import (
	"./utils"
	"github.com/magiconair/properties"
	"io"
	"log"
	"math/big"
	"os"
	"time"
)

func main() {

	f, err := os.OpenFile("./pharming.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	p := properties.MustLoadFile("default.properties", properties.UTF8)
	walletAddr := p.GetString("wallet.address", utils.EmptyAddr)
	harvestToken := p.GetString("harvest.token", utils.EmptyAddr)
	masterChefContract := p.GetString("masterChef.address", utils.EmptyAddr)
	busdThreshold := p.GetInt64("busd.threshold", 20)
	poolId := p.GetInt64("pool.id", 1)
	slippage := p.GetFloat64("slippage.limit", 0.95)
	referralAddress := p.GetString("referrer.address", utils.EmptyAddr)
	isRefer := p.GetBool("refer.isEnabled", false)
	log.Println("Properties Loaded")

	isApproved := utils.IsApproved(harvestToken, walletAddr)
	if !isApproved {
		approveTx := utils.Approve(harvestToken)
		log.Println("Approved token:" + harvestToken + " for swapping with tx:" + approveTx)
	}
	log.Println(harvestToken + " is approved for swapping")

	for {
		pendingTokenWeiCount := utils.PendingTokens(p.GetString("pending.function", "0x0000"), poolId, masterChefContract, walletAddr)
		tokenUSDPrice := utils.GetTokenUSDPrice(harvestToken)
		log.Println(harvestToken + " USD price: $" + tokenUSDPrice.String())
		pendingTokenCount := new(big.Float).Quo(new(big.Float).SetInt(&pendingTokenWeiCount), big.NewFloat(utils.WEIMultiplier))
		totalHarvestValue := new(big.Float).Mul(tokenUSDPrice, pendingTokenCount)
		log.Println("Total pending harvest value: $" + totalHarvestValue.String())
		if totalHarvestValue.Cmp(new(big.Float).SetInt64(busdThreshold)) == 1 {
			log.Println("Harvesting..")
			harvestTx := utils.Harvest(poolId, masterChefContract, referralAddress, isRefer)
			log.Println("Harvested token:" + harvestToken + " with tx:" + harvestTx)
			time.Sleep(10 * time.Second)

			balance := utils.BalanceOf(harvestToken, walletAddr)
			minBalance := new(big.Float).Mul(new(big.Float).Mul(totalHarvestValue, big.NewFloat(utils.WEIMultiplier)), new(big.Float).SetFloat64(slippage))
			minBalanceInt := new(big.Int)
			minBalance.Int(minBalanceInt)

			now := time.Now()
			nanos := now.UnixNano() / 1000000
			swapTx := utils.SwapApprovedTokenToBUSD(balance, *minBalanceInt, walletAddr, *new(big.Int).SetUint64(uint64(nanos)), harvestToken)
			log.Println("Swapped with min of $" + new(big.Float).Mul(totalHarvestValue, big.NewFloat(utils.WEIMultiplier)).String() + "BUSD at tx: " + swapTx)
		}
		time.Sleep(15 * time.Second)
	}
}
