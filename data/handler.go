package data

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/DonggyuLim/Alliance-Rank/request"
	"github.com/DonggyuLim/Alliance-Rank/utils"
	"github.com/imroc/req/v3"
	alliancemoduletypes "github.com/terra-money/alliance/x/alliance/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func GetDelegations2(height, chainCode int) (request.DelegationRequest, error) {

	value := fmt.Sprintf("%v", height)

	client := req.R().
		SetHeader("x-cosmos-block-height", value).SetHeader("Content-Type", "application/json")
	endpoint := fmt.Sprintf("%s/terra/alliances/delegations?pagination.limit=1000000",
		GetEndopoint(chainCode),
		// GetAddress(chainCode, address),
	)

	var req request.DelegationRequest
	_, err := client.SetSuccessResult(&req).Get(endpoint)

	return req, err
}
func GetRewards2(chainCode, height int, delegator, validator, denom string) ([]request.Reward, error) {

	client := req.R().
		SetHeader("x-cosmos-block-height", fmt.Sprintf("%v", height)).SetHeader("Content-Type", "application/json")
	var req request.RewardRequest
	endpoint := fmt.Sprintf("%s/terra/alliances/rewards/%s/%s/%s",
		GetEndopoint(chainCode),
		delegator,
		validator,
		denom,
	)
	//{delegator_addr}/{validator_addr}/{denom}
	// endpoint := fmt.Sprintf("%s/terra/alliances/rewards/%s/{validator_addr}/{denom}", chain, el.deligator, validator, denom)
	_, err := client.SetSuccessResult(&req).Get(endpoint)

	return req.Rewards, err
}
func GetDelegations(c alliancemoduletypes.QueryClient, height int) (*alliancemoduletypes.QueryAlliancesDelegationsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := &alliancemoduletypes.QueryAllAlliancesDelegationsRequest{}
	md := metadata.New(map[string]string{"x-cosmos-block-height": fmt.Sprintf("%v", height)})
	ctx = metadata.NewOutgoingContext(ctx, md)
	var header metadata.MD
	res, err := c.AllAlliancesDelegations(ctx, req, grpc.Header(&header))
	if err != nil {
		return nil, err
	}
	return res, err

}

func GetRewards(c alliancemoduletypes.QueryClient, height int, delegator, validator, denom string) (*alliancemoduletypes.QueryAllianceDelegationRewardsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := &alliancemoduletypes.QueryAllianceDelegationRewardsRequest{
		DelegatorAddr: delegator,
		ValidatorAddr: validator,
		Denom:         denom,
	}
	md := metadata.New(map[string]string{"x-cosmos-block-height": fmt.Sprintf("%v", height)})
	ctx = metadata.NewOutgoingContext(ctx, md)
	var header metadata.MD
	res, err := c.AllianceDelegationRewards(ctx, req, grpc.Header(&header))
	if err != nil {
		return nil, err
	}
	return res, err

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

func ReturnHeight(chainCode int) int {
	var height int
	switch chainCode {
	case 0:
		i, err := strconv.Atoi(utils.LoadENV("HEIGHT", "atr.env"))
		utils.PanicError(err)
		height = i
	case 1:
		i, err := strconv.Atoi(utils.LoadENV("HEIGHT", "har.env"))
		utils.PanicError(err)
		height = i
	case 2:
		i, err := strconv.Atoi(utils.LoadENV("HEIGHT", "cor.env"))
		utils.PanicError(err)
		height = i
	case 3:
		i, err := strconv.Atoi(utils.LoadENV("HEIGHT", "ord.env"))
		utils.PanicError(err)
		height = i
	}
	return height
}

func WriteHeight(chainCode, height int) {
	switch chainCode {
	case 0:
		utils.WriteENV("HEIGHT", strconv.Itoa(height), "atr.env")
	case 1:
		utils.WriteENV("HEIGHT", strconv.Itoa(height), "har.env")
	case 2:
		utils.WriteENV("HEIGHT", strconv.Itoa(height), "cor.env")
	case 3:
		utils.WriteENV("HEIGHT", strconv.Itoa(height), "ord.env")
	}
}
