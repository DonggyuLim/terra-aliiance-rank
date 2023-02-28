package data

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/DonggyuLim/Alliance-Rank/account"
	"github.com/DonggyuLim/Alliance-Rank/client"
	"github.com/DonggyuLim/Alliance-Rank/db"
	"github.com/DonggyuLim/Alliance-Rank/utils"
	"github.com/dariubs/percent"
	"golang.org/x/sync/errgroup"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	sCOR = "ibc/D7AA592A1C1C00FE7C9E15F4BB7ADB4B779627DD3FBB3C877CD4DB27F56E35B4"
	sORD = "ibc/3FA98D26F2D6CCB58D8E4D1B332C6EB8EE4AC7E3F0AD5B5B05201155CEB1AD1D"
	sATR = "ibc/95287CFB16A09D3FE1D0B1E34B6725A380DD2A40AEF4F496B3DAF6F0D901695B"
	sHAR = "ibc/51B1594844CCB9438C4EF3720B7ADD4398AC5D52E073CA7E592E675C6E4163EF"
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
	w.Add(1)
	// go MakeReward(w, ATREIDES)
	go MakeReward(w, Harkonnen)
	// go MakeReward(w, CORRINO)
	// go MakeReward(w, ORDOS)
	// go MakeTotal(w)
	wg.Wait()
}

func MakeTotal(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		time.Sleep(time.Second * 600)
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
	var height int

	switch chainCode {
	case 0:
		height = 54247
	case 1:
		height = 67227
	case 2:
		height = 253258
	case 3:
		height = 108588
	}

	c := client.QueryClient(chainCode)

	// lastBlock := GetLastBlock(chainCode)
	// if height > lastBlock {
	// 	fmt.Printf("height : %v lastblock:%v time lock\n", height, lastBlock)
	// 	time.Sleep(time.Minute * 5)
	// 	continue
	// }
	for height < 300000 {
		// delegationsData, err := GetDelegations(height, chainCode)
		delegationData, err := GetDelegations(c, height)

		delegations := delegationData.Delegations
		// fmt.Println(delegations)
		if len(delegations) == 0 || err != nil {
			// fmt.Printf("chain : %v height: %v lastBlock: %v Not Delegate \n", chainCode, height, lastBlock)
			height += 1
			continue
		}
		fmt.Printf("chain: %v  height: %v delecount: %v  Start!\n", chainCode, height, len(delegations))
		g, _ := errgroup.WithContext(context.Background())
		for i := 0; i <= len(delegations)-1; i++ {
			delegation := delegations[i].Delegation

			go func() error {
				// resReward, err := GetRewards(
				// 	chainCode,
				// 	height,
				// 	delegation.DelegatorAddress,
				// 	delegation.ValidatorAddress,
				// 	delegation.Denom,
				// )
				resReward, err := GetRewards(
					c,
					height,
					delegation.DelegatorAddress,
					delegation.ValidatorAddress,
					delegation.Denom,
				)
				if err != nil || len(resReward.Rewards) == 0 {
					// fmt.Printf("chain: %v height:%v Not Reward!\n", chainCode, height)
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
					SHAR:       0,
					SATR:       0,
				}
				// fmt.Println("reward loop start!")
				for _, re := range resReward.Rewards {
					switch re.Denom {
					case sATR:
						reward.SATR = int(re.Amount.Int64())
					case sHAR:

						reward.SHAR = int(re.Amount.Int64())
					case sCOR:

						reward.SCOR = int(re.Amount.Int64())
					case sORD:

						reward.SORD = int(re.Amount.Int64())
					case uatr:

						reward.UAtr = int(re.Amount.Int64())
					case uhar:

						reward.UHar = int(re.Amount.Int64())
					case ucor:

						reward.UCor = int(re.Amount.Int64())
					case uord:

						reward.UOrd = int(re.Amount.Int64())

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

						if o.UAtr > reward.UAtr && (percent.PercentOf(claimAtr, o.UAtr) >= 90) {
							//Tax 제외
							fmt.Printf("Claim! chain : %v height :%v account :%v \n", chainCode, height, delegation.DelegatorAddress)
							// utils.PrettyJson(o)
							// utils.PrettyJson(resReward)
							claimSCOR := (o.SCOR - reward.SCOR) + a.Atreides.Claim.SCOR
							claimSORD := (o.SORD - reward.SORD) + a.Atreides.Claim.SORD
							claimSHAR := (o.SHAR - reward.SHAR) + a.Atreides.Claim.SHAR
							claimSATR := (o.SATR - reward.SATR) + a.Atreides.Claim.SATR
							claimUpdate := bson.D{
								{
									Key: "$set", Value: bson.D{
										{Key: "atreides.claim", Value: bson.D{
											{Key: "uatr", Value: claimAtr + a.Atreides.Claim.UAtr},
											{Key: "scor", Value: claimSCOR},
											{Key: "sord", Value: claimSORD},
											{Key: "shar", Value: claimSHAR},
											{Key: "satr", Value: claimSATR},
										},
										},
									},
								},
								{
									Key: "$set", Value: bson.D{
										{Key: fmt.Sprintf("atreides.rewards.%s.uatr", delegation.ValidatorAddress), Value: reward.UAtr},
										{Key: fmt.Sprintf("atreides.rewards.%s.scor", delegation.ValidatorAddress), Value: reward.SCOR},
										{Key: fmt.Sprintf("atreides.rewards.%s.sord", delegation.ValidatorAddress), Value: reward.SORD},
										{Key: fmt.Sprintf("atreides.rewards.%s.shar", delegation.ValidatorAddress), Value: reward.SHAR},
										{Key: fmt.Sprintf("atreides.rewards.%s.satr", delegation.ValidatorAddress), Value: reward.SATR},
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
										{Key: fmt.Sprintf("atreides.rewards.%s.shar", delegation.ValidatorAddress), Value: reward.SHAR},
										{Key: fmt.Sprintf("atreides.rewards.%s.satr", delegation.ValidatorAddress), Value: reward.SATR},
									},
								},
							}

							db.UpdateOne(filter, update)
						}

					case 1:
						o := a.Harkonnen.Rewards[delegation.ValidatorAddress]
						claimhar := (o.UHar - reward.UHar)
						if o.UHar > reward.UHar && (percent.PercentOf(claimhar, o.UHar) >= 90) {

							fmt.Printf("Claim! chain: %v height: %v account: %v \n", chainCode, height, delegation.DelegatorAddress)

							claimSCOR := (o.SCOR - reward.SCOR) + a.Harkonnen.Claim.SCOR
							claimSORD := (o.SORD - reward.SORD) + a.Harkonnen.Claim.SORD
							claimSHAR := (o.SHAR - reward.SHAR) + a.Harkonnen.Claim.SHAR
							claimSATR := (o.SATR - reward.SATR) + a.Harkonnen.Claim.SATR
							claimUpdate := bson.D{
								{
									Key: "$set", Value: bson.D{
										{Key: "harkonnen.claim", Value: bson.D{
											{Key: fmt.Sprintf("harkonnen.claim.%s.uhar", delegation.ValidatorAddress), Value: claimhar + a.Harkonnen.Claim.UHar},
											{Key: fmt.Sprintf("harkonnen.claim.%s.scor", delegation.ValidatorAddress), Value: claimSCOR},
											{Key: fmt.Sprintf("harkonnen.claim.%s.sord", delegation.ValidatorAddress), Value: claimSORD},
											{Key: fmt.Sprintf("harkonnen.claim.%s.shar", delegation.ValidatorAddress), Value: claimSHAR},
											{Key: fmt.Sprintf("harkonnen.claim.%s.satr", delegation.ValidatorAddress), Value: claimSATR},
										},
										},
									},
								},
								{
									Key: "$set", Value: bson.D{
										{Key: fmt.Sprintf("harkonnen.rewards.%s.uhar", delegation.ValidatorAddress), Value: reward.UHar},
										{Key: fmt.Sprintf("harkonnen.rewards.%s.scor", delegation.ValidatorAddress), Value: reward.SCOR},
										{Key: fmt.Sprintf("harkonnen.rewards.%s.sord", delegation.ValidatorAddress), Value: reward.SORD},
										{Key: fmt.Sprintf("harkonnen.rewards.%s.shar", delegation.ValidatorAddress), Value: reward.SHAR},
										{Key: fmt.Sprintf("harkonnen.rewards.%s.satr", delegation.ValidatorAddress), Value: reward.SATR},
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
										{Key: fmt.Sprintf("harkonnen.rewards.%s.shar", delegation.ValidatorAddress), Value: reward.SHAR},
										{Key: fmt.Sprintf("harkonnen.rewards.%s.satr", delegation.ValidatorAddress), Value: reward.SATR},
									},
								},
							}
							db.UpdateOne(filter, update)
						}

					case 2:

						o := a.Corrino.Rewards[delegation.ValidatorAddress]
						claimCor := (o.UCor - reward.UCor)
						if o.UCor > reward.UCor && (percent.PercentOf(claimCor, o.UCor) >= 90) {

							fmt.Printf("Claim! chain : %v height :%v account :%v \n", chainCode, height, delegation.DelegatorAddress)
							// utils.PrettyJson(o)
							// utils.PrettyJson(resReward)
							claimSCOR := (o.SCOR - reward.SCOR) + a.Corrino.Claim.SCOR
							claimSORD := (o.SORD - reward.SORD) + a.Corrino.Claim.SORD
							claimSHAR := (o.SHAR - reward.SHAR) + a.Corrino.Claim.SHAR
							claimSATR := (o.SATR - reward.SATR) + a.Corrino.Claim.SATR
							claimUpdate := bson.D{
								{
									Key: "$set", Value: bson.D{
										{Key: "corrino.claim", Value: bson.D{
											{Key: "ucor", Value: claimCor + a.Corrino.Claim.UCor},
											{Key: "scor", Value: claimSCOR},
											{Key: "sord", Value: claimSORD},
											{Key: "shar", Value: claimSHAR},
											{Key: "satr", Value: claimSATR},
										},
										},
									},
								},
								{Key: "$set", Value: bson.D{
									{Key: fmt.Sprintf("corrino.rewards.%s.ucor", delegation.ValidatorAddress), Value: reward.UCor},
									{Key: fmt.Sprintf("corrino.rewards.%s.scor", delegation.ValidatorAddress), Value: reward.SCOR},
									{Key: fmt.Sprintf("corrino.rewards.%s.sord", delegation.ValidatorAddress), Value: reward.SORD},
									{Key: fmt.Sprintf("corrino.rewards.%s.shar", delegation.ValidatorAddress), Value: reward.SHAR},
									{Key: fmt.Sprintf("corrino.rewards.%s.satr", delegation.ValidatorAddress), Value: reward.SATR},
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
										{Key: fmt.Sprintf("corrino.rewards.%s.shar", delegation.ValidatorAddress), Value: reward.SHAR},
										{Key: fmt.Sprintf("corrino.rewards.%s.satr", delegation.ValidatorAddress), Value: reward.SATR},
									},
								},
							}

							db.UpdateOne(filter, update)
						}
					case 3:
						o := a.Ordos.Rewards[delegation.ValidatorAddress]
						claimOrd := (o.UOrd - reward.UOrd)
						if o.UOrd > reward.UOrd && (percent.PercentOf(claimOrd, o.UOrd) >= 90) {

							fmt.Printf("Claim! chain : %v height :%v account :%v \n", chainCode, height, delegation.DelegatorAddress)
							// utils.PrettyJson(o)
							// utils.PrettyJson(resReward)

							claimSCOR := (o.SCOR - reward.SCOR) + a.Ordos.Claim.SCOR
							claimSORD := (o.SORD - reward.SORD) + a.Ordos.Claim.SORD
							claimSHAR := (o.SHAR - reward.SHAR) + a.Ordos.Claim.SHAR
							claimSATR := (o.SATR - reward.SATR) + a.Ordos.Claim.SATR
							claimUpdate := bson.D{
								{
									Key: "$set", Value: bson.D{
										{Key: "ordos.claim", Value: bson.D{
											{Key: "ucor", Value: claimOrd + a.Corrino.Claim.UOrd},
											{Key: "scor", Value: claimSCOR},
											{Key: "sord", Value: claimSORD},
											{Key: "shar", Value: claimSHAR},
											{Key: "satr", Value: claimSATR},
										},
										},
									},
								},
								{
									Key: "$set", Value: bson.D{
										{Key: fmt.Sprintf("ordos.rewards.%s.uord", delegation.ValidatorAddress), Value: reward.UOrd},
										{Key: fmt.Sprintf("ordos.rewards.%s.scor", delegation.ValidatorAddress), Value: reward.SCOR},
										{Key: fmt.Sprintf("ordos.rewards.%s.sord", delegation.ValidatorAddress), Value: reward.SORD},
										{Key: fmt.Sprintf("ordos.rewards.%s.shar", delegation.ValidatorAddress), Value: reward.SHAR},
										{Key: fmt.Sprintf("ordos.rewards.%s.satr", delegation.ValidatorAddress), Value: reward.SATR},
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
										{Key: fmt.Sprintf("ordos.rewards.%s.shar", delegation.ValidatorAddress), Value: reward.SHAR},
										{Key: fmt.Sprintf("ordos.rewards.%s.satr", delegation.ValidatorAddress), Value: reward.SATR},
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
			}()

			if err := g.Wait(); err != nil {
				log.Panicln(err.Error())
			}

		}
		height += 1
	}
}
