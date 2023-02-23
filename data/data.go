package data

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/DonggyuLim/Alliance-Rank/account"
	"github.com/DonggyuLim/Alliance-Rank/db"
	"github.com/DonggyuLim/Alliance-Rank/utils"
	"go.mongodb.org/mongo-driver/bson"
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
	// go MakeReward(w, ATREIDES)
	// go MakeReward(w, Harkonnen)
	go MakeReward(w, CORRINO)
	// go MakeReward(w, ORDOS)
	go MakeTotal(w)
	wg.Wait()
}

func MakeTotal(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		time.Sleep(time.Second * 300)
		fmt.Println("Total!")
		accountList, err := db.Find("", "", "total.total", 100000000)
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
	height := 100961
	for {
		lastBlock := GetLastBlock(chainCode)
		if height >= GetLastBlock(chainCode) {
			height = lastBlock
			time.Sleep(2 * time.Second)
		}

		delegationsData, err := GetDelegations(height, chainCode)
		delegations := delegationsData.Deligations
		if len(delegations) == 0 || err != nil {
			fmt.Printf("chain : %v height: %v lastBlock: %v Not Delegate \n", chainCode, height, lastBlock)
			height += 1
			continue
		} else {
			fmt.Printf("chain: %v  height: %v delecount: %v  lastblock:%v \n", chainCode, height, len(delegations), lastBlock)
		}

		w := &sync.WaitGroup{}
		w.Add(len(delegations))
		fmt.Println("w :", w)
		fmt.Println("delegation loop start!")
		for i := 0; i <= len(delegations)-1; i++ {
			delegation := delegations[i].Delegation
			fmt.Println("delegation : ", delegation)
			go func() {
				fmt.Println("Go Routine Start!")
				resReward, err := GetRewards(
					chainCode,
					height,
					delegation.DelegatorAddress,
					delegation.ValidatorAddress,
					delegation.Denom,
				)
				if err != nil || len(resReward) == 0 {
					fmt.Println("Not Reward!")
				}
				reward := account.Reward{
					LastHeight: uint(height),
					UAtr:       0,
					UHar:       0,
					UCor:       0,
					UOrd:       0,
					SCOR:       0,
					SORD:       0,
				}
				fmt.Println("reward loop start!")
				for _, re := range resReward {
					switch re.Denom {
					case sCOR:
						amount, err := strconv.Atoi(re.Amount)
						utils.PanicError(err)
						reward.SCOR = uint(amount)
					case sORD:
						amount, err := strconv.Atoi(re.Amount)
						utils.PanicError(err)
						reward.SORD = uint(amount)
					case uatr:
						amount, err := strconv.Atoi(re.Amount)
						utils.PanicError(err)
						reward.UAtr = uint(amount)
					case uhar:
						amount, err := strconv.Atoi(re.Amount)
						utils.PanicError(err)
						reward.UHar = uint(amount)
					case ucor:
						amount, err := strconv.Atoi(re.Amount)
						utils.PanicError(err)
						reward.UCor = uint(amount)
					case uord:
						amount, err := strconv.Atoi(re.Amount)
						utils.PanicError(err)
						reward.UOrd = uint(amount)
					}
				}
				filter := bson.D{{Key: "address", Value: delegation.DelegatorAddress}}
				a, ok := db.FindOne(filter)
				fmt.Println(ok)
				switch ok {
				case nil:
					fmt.Println("Exsits!!!")
					switch chainCode {
					case 0:
						o := a.Atreides.Rewards[delegation.ValidatorAddress]
						if o.UAtr > reward.UAtr {
							fmt.Println("Claim")
							claimAtr := (o.UAtr - reward.UAtr) + a.Atreides.Claim.UAtr
							claimSCOR := (o.SCOR - reward.SCOR) + a.Atreides.Claim.SCOR
							claimSORD := (o.SORD - reward.SORD) + a.Atreides.Claim.SORD
							claimUpdate := bson.D{{Key: "$set", Value: bson.D{{Key: "atreides", Value: bson.M{
								"uatr": claimAtr,
								"scor": claimSCOR,
								"sord": claimSORD,
							}}}}}
							db.UpdateOne(filter, claimUpdate)
							update := bson.D{{
								Key: "$set", Value: bson.D{
									{Key: "atreides", Value: bson.D{
										{Key: "reward", Value: bson.D{
											{Key: delegation.ValidatorAddress, Value: reward},
										}},
									},
									},
								},
							},
							}
							db.UpdateOne(filter, update)
						} else {
							fmt.Println("Update Reward!")
							update := bson.D{{
								Key: "$set", Value: bson.D{
									{Key: "atreides", Value: bson.D{
										{Key: "reward", Value: bson.D{
											{Key: delegation.ValidatorAddress, Value: reward},
										}},
									},
									},
								},
							},
							}
							db.UpdateOne(filter, update)
						}

					case 1:
						o := a.Harkonnen.Rewards[delegation.ValidatorAddress]
						if o.UHar > reward.UHar {
							fmt.Println("Claim")
							claimhar := (o.UHar - reward.UHar) + a.Harkonnen.Claim.UHar
							claimSCOR := (o.SCOR - reward.SCOR) + a.Harkonnen.Claim.SCOR
							claimSORD := (o.SORD - reward.SORD) + a.Harkonnen.Claim.SORD
							claimUpdate := bson.D{{Key: "$set", Value: bson.D{{Key: "harkonnen", Value: bson.M{
								"uHar": claimhar,
								"scor": claimSCOR,
								"sord": claimSORD,
							}}}}}
							db.UpdateOne(filter, claimUpdate)
							update := bson.D{{
								Key: "$set", Value: bson.D{
									{Key: "harkonnen", Value: bson.D{
										{Key: "reward", Value: bson.D{
											{Key: delegation.ValidatorAddress, Value: reward},
										}},
									},
									},
								},
							},
							}
							db.UpdateOne(filter, update)
						} else {
							fmt.Println("Update Reward!")
							update := bson.D{{
								Key: "$set", Value: bson.D{
									{Key: "harkonnen", Value: bson.D{
										{Key: "reward", Value: bson.D{
											{Key: delegation.ValidatorAddress, Value: reward},
										}},
									},
									},
								},
							},
							}
							db.UpdateOne(filter, update)
						}

					case 2:

						o := a.Corrino.Rewards[delegation.ValidatorAddress]
						if o.UCor > reward.UCor {
							fmt.Println("Claim")
							claimCor := (o.UCor - reward.UCor) + a.Corrino.Claim.UCor
							claimSCOR := (o.SCOR - reward.SCOR) + a.Corrino.Claim.SCOR
							claimSORD := (o.SORD - reward.SORD) + a.Corrino.Claim.SORD
							claimUpdate := bson.D{{Key: "$set", Value: bson.D{{Key: "corrino", Value: bson.M{
								"ucor": claimCor,
								"scor": claimSCOR,
								"sord": claimSORD,
							}}}}}
							db.UpdateOne(filter, claimUpdate)
							update := bson.D{{
								Key: "$set", Value: bson.D{
									{Key: "corrino", Value: bson.D{
										{Key: "reward", Value: bson.D{
											{Key: delegation.ValidatorAddress, Value: reward},
										}},
									},
									},
								},
							},
							}
							db.UpdateOne(filter, update)
						} else {
							fmt.Println("Update Reward!")
							update := bson.D{{
								Key: "$set", Value: bson.D{
									{Key: "corrino", Value: bson.D{
										{Key: "reward", Value: bson.D{
											{Key: delegation.ValidatorAddress, Value: reward},
										}},
									},
									},
								},
							},
							}
							db.UpdateOne(filter, update)
						}
					case 3:
						o := a.Ordos.Rewards[delegation.ValidatorAddress]
						if o.UOrd > reward.UOrd {
							fmt.Println("Claim")
							claimCor := (o.UOrd - reward.UOrd) + a.Ordos.Claim.UOrd
							claimSCOR := (o.SCOR - reward.SCOR) + a.Ordos.Claim.SCOR
							claimSORD := (o.SORD - reward.SORD) + a.Ordos.Claim.SORD
							claimUpdate := bson.D{{Key: "$set", Value: bson.D{{Key: "ordos", Value: bson.M{
								"uord": claimCor,
								"scor": claimSCOR,
								"sord": claimSORD,
							}}}}}
							db.UpdateOne(filter, claimUpdate)

							update := bson.D{{
								Key: "$set", Value: bson.D{
									{Key: "ordos", Value: bson.D{
										{Key: "reward", Value: bson.D{
											{Key: delegation.ValidatorAddress, Value: reward},
										}},
									},
									},
								},
							},
							}
							db.UpdateOne(filter, update)
						} else {
							fmt.Println("Update Reward!")
							update := bson.D{{
								Key: "$set", Value: bson.D{
									{Key: "ordos", Value: bson.D{
										{Key: "reward", Value: bson.D{
											{Key: delegation.ValidatorAddress, Value: reward},
										}},
									},
									},
								},
							},
							}
							db.UpdateOne(filter, update)
						}
					}
				default:
					fmt.Println("New Account!")
					a.SetAccount(delegation.DelegatorAddress, delegation.ValidatorAddress, reward, chainCode)
					db.Insert(a)
				}

			}()
			w.Done()
		}
		w.Wait()
		height += 1
	}

}

// func Controller(a, h, c, o <-chan DataChan, wg *sync.WaitGroup) {

// 	for {
// 		select {
// 		case d := <-a:
// 			fmt.Println("atreides")
// 			c := account.Chain{}
// 			ok := db.FindChain(d.Address, ATREIDES, &c)
// 			switch ok {
// 			case nil:
// 				c.UpdateClaimAndReward(
// 					d.Address,
// 					d.ValidatorAddress,
// 					d.Reward,
// 					ATREIDES,
// 				)
// 				c.UpdateUndelegate(0, int(d.Reward.LastHeight))
// 				filter := bson.D{{Key: "address", Value: utils.MakeAddress(d.Address)}}
// 				update := bson.D{{Key: "$set", Value: bson.D{{Key: "atreides", Value: c}}}}
// 				db.UpdateOne(filter, update)

// 			case mongo.ErrNoDocuments:
// 				account := account.Account{}
// 				account.SetAccount(d.Address, d.ValidatorAddress, d.Reward, ATREIDES)
// 				db.Insert(account)
// 			}

// 		case d := <-h:
// 			fmt.Println("harkonnen")
// 			c := account.Chain{}
// 			ok := db.FindChain(d.Address, Harkonnen, &c)
// 			switch ok {
// 			case nil:
// 				c.UpdateUndelegate(Harkonnen, int(d.Reward.LastHeight))
// 				c.UpdateClaimAndReward(
// 					d.Address,
// 					d.ValidatorAddress,
// 					d.Reward,
// 					Harkonnen,
// 				)
// 				filter := bson.D{{Key: "address", Value: utils.MakeAddress(d.Address)}}
// 				update := bson.D{{Key: "$set", Value: bson.D{{Key: "harkonnen", Value: c}}}}
// 				db.UpdateOne(filter, update)

// 			case mongo.ErrNoDocuments:
// 				account := account.Account{}
// 				account.SetAccount(d.Address, d.ValidatorAddress, d.Reward, Harkonnen)
// 				db.Insert(account)
// 			}

// 		case d := <-c:
// 			fmt.Println("corrino")
// 			c := account.Chain{}
// 			ok := db.FindChain(d.Address, CORRINO, &c)
// 			switch ok {
// 			case nil:
// 				c.UpdateUndelegate(CORRINO, int(d.Reward.LastHeight))
// 				c.UpdateClaimAndReward(
// 					d.Address,
// 					d.ValidatorAddress,
// 					d.Reward,
// 					CORRINO,
// 				)
// 				filter := bson.D{{Key: "address", Value: utils.MakeAddress(d.Address)}}
// 				update := bson.D{{Key: "$set", Value: bson.D{{Key: "corrino", Value: c}}}}
// 				db.UpdateOne(filter, update)

// 			case mongo.ErrNoDocuments:
// 				account := account.Account{}
// 				account.SetAccount(d.Address, d.ValidatorAddress, d.Reward, CORRINO)
// 				db.Insert(account)
// 			}

// 		case d := <-o:
// 			fmt.Println("ordos")
// 			c := account.Chain{}
// 			ok := db.FindChain(d.Address, ORDOS, &c)
// 			switch ok {
// 			case nil:
// 				c.UpdateUndelegate(ORDOS, int(d.Reward.LastHeight))
// 				c.UpdateClaimAndReward(
// 					d.Address,
// 					d.ValidatorAddress,
// 					d.Reward,
// 					ORDOS,
// 				)
// 				filter := bson.D{{Key: "address", Value: utils.MakeAddress(d.Address)}}
// 				update := bson.D{{Key: "$set", Value: bson.D{{Key: "ordos", Value: c}}}}
// 				db.UpdateOne(filter, update)

// 			case mongo.ErrNoDocuments:
// 				account := account.Account{}
// 				account.SetAccount(d.Address, d.ValidatorAddress, d.Reward, ORDOS)
// 				db.Insert(account)
// 			}
// 		}
// 	}
// }
