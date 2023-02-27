package data

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/DonggyuLim/Alliance-Rank/account"
	"github.com/DonggyuLim/Alliance-Rank/db"
	"github.com/DonggyuLim/Alliance-Rank/utils"
	"github.com/dariubs/percent"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
)

const (
	sCOR = "ibc/D7AA592A1C1C00FE7C9E15F4BB7ADB4B779627DD3FBB3C877CD4DB27F56E35B4"
	sORD = "ibc/3FA98D26F2D6CCB58D8E4D1B332C6EB8EE4AC7E3F0AD5B5B05201155CEB1AD1D"
	uatr = "uatr"
	uhar = "uhar"
	ucor = "ucor"
	uord = "uord"
)
const (
	ATREIDES = iota
	Harkonnen
	CORRINO
	ORDOS
)

func Main(wg *sync.WaitGroup) {
	defer wg.Done()
	w := &sync.WaitGroup{}
	w.Add(5)
	go MakeReward(w, ATREIDES)
	go MakeReward(w, Harkonnen)
	go MakeReward(w, CORRINO)
	go MakeReward(w, ORDOS)
	go MakeTotal(w)
	wg.Wait()
}

func MakeTotal(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		time.Sleep(time.Second * 60)
		fmt.Println("===========Total!================")
		accountList, err := db.FindAll()
		if err != nil || len(accountList) == 0 {
			fmt.Println("Make Total None")
			continue
		}
		for _, account := range accountList {
			account.CalculateTotal()
			filter := bson.D{{Key: "address", Value: account.Address}}
			update := bson.M{
				"$set": bson.M{
					"atreides.total":  account.Atreides.Total,
					"harkonnen.total": account.Harkonnen.Total,
					"corrino.total":   account.Corrino.Total,
					"ordos.total":     account.Ordos.Total,
					"total":           account.Total,
				},
			}
			db.UpdateOneMap(filter, update)
		}

	}

}

func MakeReward(wg *sync.WaitGroup, chainCode int) {
	defer wg.Done()
	height := ReturnHeight(chainCode)
	// height := 10000
	for {

		lastBlock := GetLastBlock(chainCode)
		if height > lastBlock {
			fmt.Printf("height : %v lastblock:%v time lock\n", height, lastBlock)
			time.Sleep(time.Minute * 5)
			continue
		}

		delegationsData, err := GetDelegations(height, chainCode)
		delegations := delegationsData.Deligations
		if len(delegations) == 0 || err != nil {
			// fmt.Printf("chain : %v height: %v lastBlock: %v Not Delegate \n", chainCode, height, lastBlock)
			height += 1
			WriteHeight(chainCode, height)
			continue
		} else {
			fmt.Printf("chain: %v  height: %v delecount: %v  lastblock:%v \n", chainCode, height, len(delegations), lastBlock)
		}

		g, _ := errgroup.WithContext(context.Background())
		for i := 0; i <= len(delegations)-1; i++ {
			delegation := delegations[i].Delegation

			g.Go(func() error {
				resReward, err := GetRewards(
					chainCode,
					height,
					delegation.DelegatorAddress,
					delegation.ValidatorAddress,
					delegation.Denom,
				)
				if err != nil || len(resReward) == 0 {
					fmt.Printf("chain: %v height:%v Not Reward!\n", chainCode, height)
					return err
				}

				reward := account.Reward{
					LastHeight: height,
					UAtr:       0,
					UHar:       0,
					UCor:       0,
					UOrd:       0,
					SCOR:       0,
					SORD:       0,
				}
				// fmt.Println("reward loop start!")
				for _, re := range resReward {
					switch re.Denom {
					case sCOR:
						amount, err := strconv.Atoi(re.Amount)
						utils.PanicError(err)
						reward.SCOR = amount
					case sORD:
						amount, err := strconv.Atoi(re.Amount)
						utils.PanicError(err)
						reward.SORD = amount
					case uatr:
						amount, err := strconv.Atoi(re.Amount)
						utils.PanicError(err)
						reward.UAtr = amount
					case uhar:
						amount, err := strconv.Atoi(re.Amount)
						utils.PanicError(err)
						reward.UHar = amount
					case ucor:
						amount, err := strconv.Atoi(re.Amount)
						utils.PanicError(err)
						reward.UCor = amount
					case uord:
						amount, err := strconv.Atoi(re.Amount)
						utils.PanicError(err)
						reward.UOrd = amount
					}
				}
				filter := bson.D{
					{Key: "address", Value: utils.MakeKey(delegation.DelegatorAddress)},
				}
				a, ok := db.FindOne(filter)

				switch ok {
				case nil:
					switch chainCode {
					case 0:
						o := a.Atreides.Rewards[delegation.ValidatorAddress]
						claimAtr := (o.UAtr - reward.UAtr)

						if o.UAtr > reward.UAtr && (percent.PercentOf(claimAtr, o.UAtr) >= 10) {
							//Tax 제외
							fmt.Printf("Claim! chain : %v height :%v account :%v \n", chainCode, height, delegation.DelegatorAddress)
							// utils.PrettyJson(o)
							// utils.PrettyJson(resReward)
							claimSCOR := (o.SCOR - reward.SCOR) + a.Atreides.Claim.SCOR
							claimSORD := (o.SORD - reward.SORD) + a.Atreides.Claim.SORD
							claimUpdate := bson.D{
								{
									Key: "$set", Value: bson.D{
										{Key: "atreides.claim", Value: bson.D{
											{Key: "uatr", Value: claimAtr + a.Atreides.Claim.UAtr},
											{Key: "scor", Value: claimSCOR},
											{Key: "sord", Value: claimSORD},
										},
										},
									},
								},
								{
									Key: "$set", Value: bson.D{
										{Key: fmt.Sprintf("atreides.rewards.%s.uatr", delegation.ValidatorAddress), Value: reward.UAtr},
										{Key: fmt.Sprintf("atreides.rewards.%s.scor", delegation.ValidatorAddress), Value: reward.SCOR},
										{Key: fmt.Sprintf("atreides.rewards.%s.sord", delegation.ValidatorAddress), Value: reward.SORD},
									},
								},
							}
							db.UpdateOne(filter, claimUpdate)

						} else {
							// fmt.Printf("Reward Update!! chain : %v height :%v account :%v\n ", chainCode, height, delegation.DelegatorAddress)
							update := bson.D{
								{
									Key: "$set", Value: bson.D{
										{Key: fmt.Sprintf("atreides.rewards.%s.uatr", delegation.ValidatorAddress), Value: reward.UAtr},
										{Key: fmt.Sprintf("atreides.rewards.%s.scor", delegation.ValidatorAddress), Value: reward.SCOR},
										{Key: fmt.Sprintf("atreides.rewards.%s.sord", delegation.ValidatorAddress), Value: reward.SORD},
									},
								},
							}

							db.UpdateOne(filter, update)
						}

					case 1:
						o := a.Harkonnen.Rewards[delegation.ValidatorAddress]
						claimhar := (o.UHar - reward.UHar)
						if o.UHar > reward.UHar && (percent.PercentOf(claimhar, o.UHar) >= 10) {

							fmt.Printf("Claim! chain: %v height: %v account: %v \n", chainCode, height, delegation.DelegatorAddress)

							claimSCOR := (o.SCOR - reward.SCOR) + a.Harkonnen.Claim.SCOR
							claimSORD := (o.SORD - reward.SORD) + a.Harkonnen.Claim.SORD
							claimUpdate := bson.D{
								{
									Key: "$set", Value: bson.D{
										{Key: "harkonnen.claim", Value: bson.D{
											{Key: fmt.Sprintf("harkonnen.claim.%s.uhar", delegation.ValidatorAddress), Value: claimhar + a.Harkonnen.Claim.UHar},
											{Key: fmt.Sprintf("harkonnen.claim.%s.scor", delegation.ValidatorAddress), Value: claimSCOR},
											{Key: fmt.Sprintf("harkonnen.claim.%s.sord", delegation.ValidatorAddress), Value: claimSORD},
										},
										},
									},
								},
								{
									Key: "$set", Value: bson.D{
										{Key: fmt.Sprintf("harkonnen.rewards.%s.uhar", delegation.ValidatorAddress), Value: reward.UHar},
										{Key: fmt.Sprintf("harkonnen.rewards.%s.scor", delegation.ValidatorAddress), Value: reward.SCOR},
										{Key: fmt.Sprintf("harkonnen.rewards.%s.sord", delegation.ValidatorAddress), Value: reward.SORD},
									},
								},
							}
							db.UpdateOne(filter, claimUpdate)

						} else {
							// fmt.Printf("Reward Update!! chain : %v height :%v account :%v\n ", chainCode, height, delegation.DelegatorAddress)
							update := bson.D{
								{
									Key: "$set", Value: bson.D{
										{Key: fmt.Sprintf("harkonnen.rewards.%s.uhar", delegation.ValidatorAddress), Value: reward.UHar},
										{Key: fmt.Sprintf("harkonnen.rewards.%s.scor", delegation.ValidatorAddress), Value: reward.SCOR},
										{Key: fmt.Sprintf("harkonnen.rewards.%s.sord", delegation.ValidatorAddress), Value: reward.SORD},
									},
								},
							}
							db.UpdateOne(filter, update)
						}

					case 2:

						o := a.Corrino.Rewards[delegation.ValidatorAddress]
						claimCor := (o.UCor - reward.UCor)
						if o.UCor > reward.UCor && (percent.PercentOf(claimCor, o.UCor) >= 10) {

							fmt.Printf("Claim! chain : %v height :%v account :%v \n", chainCode, height, delegation.DelegatorAddress)
							// utils.PrettyJson(o)
							// utils.PrettyJson(resReward)
							claimSCOR := (o.SCOR - reward.SCOR) + a.Corrino.Claim.SCOR
							claimSORD := (o.SORD - reward.SORD) + a.Corrino.Claim.SORD
							claimUpdate := bson.D{
								{
									Key: "$set", Value: bson.D{
										{Key: "corrino.claim", Value: bson.D{
											{Key: "ucor", Value: claimCor + a.Corrino.Claim.UCor},
											{Key: "scor", Value: claimSCOR},
											{Key: "sord", Value: claimSORD},
										},
										},
									},
								},
								{Key: "$set", Value: bson.D{
									{Key: fmt.Sprintf("corrino.rewards.%s.ucor", delegation.ValidatorAddress), Value: reward.UCor},
									{Key: fmt.Sprintf("corrino.rewards.%s.scor", delegation.ValidatorAddress), Value: reward.SCOR},
									{Key: fmt.Sprintf("corrino.rewards.%s.sord", delegation.ValidatorAddress), Value: reward.SORD},
								},
								},
							}
							db.UpdateOne(filter, claimUpdate)

						} else {
							// fmt.Printf("Reward Update!! chain : %v height :%v account :%v\n ", chainCode, height, delegation.DelegatorAddress)
							update := bson.D{
								{
									Key: "$set", Value: bson.D{
										{Key: fmt.Sprintf("corrino.rewards.%s.ucor", delegation.ValidatorAddress), Value: reward.UCor},
										{Key: fmt.Sprintf("corrino.rewards.%s.scor", delegation.ValidatorAddress), Value: reward.SCOR},
										{Key: fmt.Sprintf("corrino.rewards.%s.sord", delegation.ValidatorAddress), Value: reward.SORD},
									},
								},
							}

							db.UpdateOne(filter, update)
						}
					case 3:
						o := a.Ordos.Rewards[delegation.ValidatorAddress]
						claimOrd := (o.UOrd - reward.UOrd)
						if o.UOrd > reward.UOrd && (percent.PercentOf(claimOrd, o.UOrd) >= 10) {

							fmt.Printf("Claim! chain : %v height :%v account :%v \n", chainCode, height, delegation.DelegatorAddress)
							// utils.PrettyJson(o)
							// utils.PrettyJson(resReward)

							claimSCOR := (o.SCOR - reward.SCOR) + a.Ordos.Claim.SCOR
							claimSORD := (o.SORD - reward.SORD) + a.Ordos.Claim.SORD
							claimUpdate := bson.D{
								{
									Key: "$set", Value: bson.D{
										{Key: "ordos.claim", Value: bson.D{
											{Key: "ucor", Value: claimOrd + a.Corrino.Claim.UOrd},
											{Key: "scor", Value: claimSCOR},
											{Key: "sord", Value: claimSORD},
										},
										},
									},
								},
								{
									Key: "$set", Value: bson.D{
										{Key: fmt.Sprintf("ordos.rewards.%s.uord", delegation.ValidatorAddress), Value: reward.UOrd},
										{Key: fmt.Sprintf("ordos.rewards.%s.scor", delegation.ValidatorAddress), Value: reward.SCOR},
										{Key: fmt.Sprintf("ordos.rewards.%s.sord", delegation.ValidatorAddress), Value: reward.SORD},
									},
								},
							}

							db.UpdateOne(filter, claimUpdate)

						} else {
							// fmt.Printf("Reward Update!! chain : %v height :%v account :%v\n ", chainCode, height, delegation.DelegatorAddress)
							update := bson.D{
								{
									Key: "$set", Value: bson.D{
										{Key: fmt.Sprintf("ordos.rewards.%s.uord", delegation.ValidatorAddress), Value: reward.UOrd},
										{Key: fmt.Sprintf("ordos.rewards.%s.scor", delegation.ValidatorAddress), Value: reward.SCOR},
										{Key: fmt.Sprintf("ordos.rewards.%s.sord", delegation.ValidatorAddress), Value: reward.SORD},
									},
								},
							}

							db.UpdateOne(filter, update)
						}
					}
				case mongo.ErrNoDocuments:
					fmt.Println("New Account!")
					key := utils.MakeKey(delegation.DelegatorAddress)
					a.SetAccount(key, delegation.ValidatorAddress, reward, chainCode)
					db.Insert(a)
				}
				return nil
			})

		}
		height += 1
		WriteHeight(chainCode, height+1)
		if err := g.Wait(); err != nil {
			log.Panicln(err.Error())
		}

	}

}
