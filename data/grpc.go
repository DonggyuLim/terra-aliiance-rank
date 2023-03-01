package data

import (
	"context"
	"fmt"
	"time"

	"github.com/DonggyuLim/Alliance-Rank/request"
	"github.com/DonggyuLim/Alliance-Rank/utils"
	"github.com/imroc/req/v3"
	module "github.com/terra-money/alliance/x/alliance/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GetDelegations(c module.QueryClient) []module.DelegationResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	req := &module.QueryAllAlliancesDelegationsRequest{}
	// md := metadata.New(map[string]string{"x-cosmos-block-height": fmt.Sprintf("%v", height)})
	// ctx = metadata.NewOutgoingContext(ctx, md)
	// var header metadata.MD
	res, err := c.AllAlliancesDelegations(ctx, req)
	utils.PanicError(err)
	return res.Delegations
}

func GetDelegationsHeight(c module.QueryClient, height int) []module.DelegationResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	req := &module.QueryAllAlliancesDelegationsRequest{}
	md := metadata.New(map[string]string{"x-cosmos-block-height": fmt.Sprintf("%v", height)})
	ctx = metadata.NewOutgoingContext(ctx, md)
	var header metadata.MD
	res, err := c.AllAlliancesDelegations(ctx, req, grpc.Header(&header))
	utils.PanicError(err)
	return res.Delegations
}

func GetDelegationsByDelegatorHeight(c module.QueryClient, delegator string, height int) []module.DelegationResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	req := &module.QueryAlliancesDelegationsRequest{
		DelegatorAddr: delegator,
	}
	md := metadata.New(map[string]string{"x-cosmos-block-height": fmt.Sprintf("%v", height)})
	ctx = metadata.NewOutgoingContext(ctx, md)
	var header metadata.MD
	res, err := c.AlliancesDelegation(ctx, req, grpc.Header(&header))
	utils.PanicError(err)
	return res.Delegations
}

func GetDelegationsByValidatorHeight(c module.QueryClient, delegator, validator string, height int) []module.DelegationResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	req := &module.QueryAlliancesDelegationByValidatorRequest{
		DelegatorAddr: delegator,
		ValidatorAddr: validator,
	}
	md := metadata.New(map[string]string{"x-cosmos-block-height": fmt.Sprintf("%v", height)})
	ctx = metadata.NewOutgoingContext(ctx, md)
	var header metadata.MD
	res, err := c.AlliancesDelegationByValidator(ctx, req, grpc.Header(&header))
	utils.PanicError(err)
	return res.Delegations
}

func GetDelegationsByDelegator(c module.QueryClient, delegator string) []module.DelegationResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	req := &module.QueryAlliancesDelegationsRequest{
		DelegatorAddr: delegator,
	}
	// md := metadata.New(map[string]string{"x-cosmos-block-height": fmt.Sprintf("%v", height)})
	// ctx = metadata.NewOutgoingContext(ctx, md)
	// var header metadata.MD
	res, err := c.AlliancesDelegation(ctx, req)
	utils.PanicError(err)
	return res.Delegations
}

func GetRewards(c module.QueryClient, delegator, validator, denom string) (*module.QueryAllianceDelegationRewardsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	req := &module.QueryAllianceDelegationRewardsRequest{
		DelegatorAddr: delegator,
		ValidatorAddr: validator,
		Denom:         denom,
	}
	// md := metadata.New(map[string]string{"x-cosmos-block-height": fmt.Sprintf("%v", height)})
	// ctx = metadata.NewOutgoingContext(ctx, md)

	res, err := c.AllianceDelegationRewards(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, err
}

func GetRewardHeight(c module.QueryClient, delegator, validator, denom string, height int) (*module.QueryAllianceDelegationRewardsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	req := &module.QueryAllianceDelegationRewardsRequest{
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

func GetValidators(c module.QueryClient) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	req := &module.QueryAllAllianceValidatorsRequest{}
	res, err := c.AllAllianceValidators(ctx, req)
	utils.PanicError(err)
	a := []string{}
	for _, el := range res.Validators {
		a = append(a, el.ValidatorAddr)
	}
	return a
}

func GetAccounts(chainCode int) []string {

	client := req.R()

	var req request.AccountsReq
	endpoint := fmt.Sprintf("%s/cosmos/auth/v1beta1/accounts",
		GetEndopoint(chainCode),
	)

	_, err := client.SetSuccessResult(&req).Get(endpoint)
	utils.PanicError(err)
	a := []string{}
	for _, el := range req.Account {
		a = append(a, el.Address)
	}
	return a
}
