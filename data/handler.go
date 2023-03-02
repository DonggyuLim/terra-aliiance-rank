package data

import (
	"fmt"
	"strconv"

	"github.com/DonggyuLim/Alliance-Rank/request"
	"github.com/DonggyuLim/Alliance-Rank/utils"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/imroc/req/v3"
	module "github.com/terra-money/alliance/x/alliance/types"
)

func GetEndopoint(a int) string {
	switch a {
	case 0:
		return "http://localhost:1317"
	case 1:
		return "http://localhost:2317"
	case 2:
		return "http://localhost:3317"
	case 3:
		return "http://localhost:4317"
	}
	return ""
}

func GetAddress(chainCode int, address string) string {
	switch chainCode {
	case 0:
		return utils.MakeAddressPrefix(address, "atreides")
	case 1:
		return utils.MakeAddressPrefix(address, "harkonnen")
	case 2:
		return utils.MakeAddressPrefix(address, "corrino")
	case 3:
		return utils.MakeAddressPrefix(address, "ordos")
	}
	return ""
}

func GetDelegation(height, chainCode int) request.DelegationRequest {

	value := fmt.Sprintf("%v", height)
	// fmt.Println(value)
	client := req.R().
		SetHeader("x-cosmos-block-height", value).SetHeader("Content-Type", "application/json")
	endpoint := fmt.Sprintf("%s/terra/alliances/delegations?pagination.limit=10000000",
		GetEndopoint(chainCode),
		// GetAddress(chainCode, address),
	)

	var req request.DelegationRequest
	_, err := client.SetSuccessResult(&req).Get(endpoint)
	utils.PanicError(err)
	return req
}

func GetLastBlock(chainCode int) int {
	client := req.R()

	endpoint := fmt.Sprintf("%s/cosmos/base/tendermint/v1beta1/blocks/latest",
		GetEndopoint(chainCode),
	)
	var lastBlock request.LastBlock
	_, err := client.SetSuccessResult(&lastBlock).Get(endpoint)
	utils.PanicError(err)
	latestHeight, err := strconv.Atoi(lastBlock.Block.Header.Height)
	utils.PanicError(err)
	return latestHeight

}

func GetClaim(c module.QueryClient, address, validator string, height int) []types.Coin {

	var coinSlice []types.Coin
	lastHeight := height
	claimCount := 0

	for {

		res := GetDelegationsByValidatorHeight(c, address, validator, lastHeight-1)
		// utils.PrettyJson(res)
		if len(res) == 0 {
			break
		}
		re := res[0]

		coins, err := GetRewardHeight(c, address, validator, re.Balance.Denom, lastHeight-1)
		claimCount += 1
		lastHeight = int(re.Delegation.LastRewardClaimHeight)
		// fmt.Println(lastHeight, coins, address, validator)
		if err != nil || len(coins.Rewards) == 0 {
			break
		}

		coinSlice = append(coinSlice, coins.Rewards...)

	}
	// if claimCount < 10 {
	// 	fmt.Printf("Delegator :%s ,Validator : %s ,Height: %v, ClaimCount : %v\n", address, validator, height, claimCount)
	// 	utils.PrettyJson(coinSlice)
	// }

	return coinSlice
}

func GetHeight(chainCode int) int {
	switch chainCode {
	case 0:
		return 230240
	case 1:
		return 218630
	case 2:
		return 265950
	case 3:
		return 100515
	}
	return 0
}
