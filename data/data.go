package data

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/DonggyuLim/Alliance-Rank/account"
	"github.com/DonggyuLim/Alliance-Rank/client"
	"github.com/DonggyuLim/Alliance-Rank/db"
	"github.com/DonggyuLim/Alliance-Rank/utils"
	"golang.org/x/sync/errgroup"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	w.Add(4)
	go MakeReward(w, ATREIDES)
	go MakeReward(w, Harkonnen)
	go MakeReward(w, CORRINO)
	// go MakeReward(w, ORDOS)
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

	c := client.QueryClient(chainCode)
	height := GetHeight(chainCode)
	for {
		start := time.Now()
		lastheight := GetLastBlock(chainCode)
		if height > lastheight {
			height = lastheight
		}
		res := GetDelegation(height, chainCode).Deligations

		if len(res) == 0 {
			fmt.Printf("chain : %v height : %v  time : %v Not Dele!!\n", chainCode, height, start)
			height += 100
			continue
		}
		fmt.Printf("chain : %v height : %v  time : %v Start!!\n", chainCode, height, start)
		g, _ := errgroup.WithContext(context.Background())
		for i := 0; i < len(res); i++ {
			d := res[i].Delegation
			g.Go(func() error {

				rw := account.Reward{}
				claim := account.Claim{}
				reward, err := GetRewardHeight(c, d.DelegatorAddress, d.ValidatorAddress, d.Denom, height)
				if err != nil || len(reward.Rewards) == 0 {
					return nil
				}

				rw.Add(reward.Rewards)
				lastClaimHeight, _ := strconv.Atoi(d.LastRewardClaimHeight)
				// utils.PrettyJson(reward.Rewards)
				c := GetClaim(c, d.DelegatorAddress, d.ValidatorAddress, lastClaimHeight)
				// fmt.Println("Claim Count", len(c))

				claim.Add(c)
				rw.Claim = claim
				filter := bson.D{
					{Key: "address", Value: utils.MakeKey(d.DelegatorAddress)},
				}
				a, ok := db.FindOne(filter)
				switch ok {
				case mongo.ErrNoDocuments:
					// fmt.Println("New Account!", d.DelegatorAddress)
					key := utils.MakeKey(d.DelegatorAddress)
					a.SetAccount(key, d.ValidatorAddress, rw, chainCode)
					db.Insert(a)
				case nil:
					switch chainCode {
					case 0:
						// fmt.Printf("atreides Update!\n address : %s\n validator : %s \n height:%v \n denom:%v \n  ", d.DelegatorAddress, d.ValidatorAddress, height, d.Denom)
						// utils.PrettyJson(rw)
						update := bson.D{
							{
								Key: "$set", Value: bson.D{
									{Key: fmt.Sprintf("atreides.rewards.%s.uatr", d.ValidatorAddress), Value: rw.UAtr},
									{Key: fmt.Sprintf("atreides.rewards.%s.scor", d.ValidatorAddress), Value: rw.SCOR},
									{Key: fmt.Sprintf("atreides.rewards.%s.sord", d.ValidatorAddress), Value: rw.SORD},
									{Key: fmt.Sprintf("atreides.rewards.%s.shar", d.ValidatorAddress), Value: rw.SHAR},
									{Key: fmt.Sprintf("atreides.rewards.%s.satr", d.ValidatorAddress), Value: rw.SATR},
									{Key: fmt.Sprintf("atreides.rewards.%s.claim", d.ValidatorAddress), Value: rw.Claim},
								},
							},
						}
						db.UpdateOne(filter, update)
					case 1:
						// fmt.Printf("harkonnen Update!\n address : %s\n validator : %s \n height:%v \n  denom:%v \n ", d.DelegatorAddress, d.ValidatorAddress, height, d.Denom)
						// utils.PrettyJson(rw)
						update := bson.D{
							{
								Key: "$set", Value: bson.D{
									{Key: fmt.Sprintf("harkonnen.rewards.%s.uhar", d.ValidatorAddress), Value: rw.UHar},
									{Key: fmt.Sprintf("harkonnen.rewards.%s.scor", d.ValidatorAddress), Value: rw.SCOR},
									{Key: fmt.Sprintf("harkonnen.rewards.%s.sord", d.ValidatorAddress), Value: rw.SORD},
									{Key: fmt.Sprintf("harkonnen.rewards.%s.shar", d.ValidatorAddress), Value: rw.SHAR},
									{Key: fmt.Sprintf("harkonnen.rewards.%s.satr", d.ValidatorAddress), Value: rw.SATR},
									{Key: fmt.Sprintf("harkonnen.rewards.%s.claim", d.ValidatorAddress), Value: rw.Claim},
								},
							},
						}
						db.UpdateOne(filter, update)
					case 2:
						// fmt.Printf("corrino Update!\n address : %s\n validator : %s \n height:%v \n   denom:%v \n ", d.DelegatorAddress, d.ValidatorAddress, height, d.Denom)
						// utils.PrettyJson(rw)
						update := bson.D{
							{
								Key: "$set", Value: bson.D{
									{Key: fmt.Sprintf("corrino.rewards.%s.ucor", d.ValidatorAddress), Value: rw.UCor},
									{Key: fmt.Sprintf("corrino.rewards.%s.scor", d.ValidatorAddress), Value: rw.SCOR},
									{Key: fmt.Sprintf("corrino.rewards.%s.sord", d.ValidatorAddress), Value: rw.SORD},
									{Key: fmt.Sprintf("corrino.rewards.%s.shar", d.ValidatorAddress), Value: rw.SHAR},
									{Key: fmt.Sprintf("corrino.rewards.%s.satr", d.ValidatorAddress), Value: rw.SATR},
									{Key: fmt.Sprintf("corrino.rewards.%s.claim", d.ValidatorAddress), Value: rw.Claim},
								},
							},
						}
						db.UpdateOne(filter, update)
					case 3:
						// fmt.Printf("ordos Update!\n address : %s\n validator : %s \n height:%v \n   denom:%v \n ", d.DelegatorAddress, d.ValidatorAddress, height, d.Denom)
						// utils.PrettyJson(rw)
						update := bson.D{
							{
								Key: "$set", Value: bson.D{
									{Key: fmt.Sprintf("ordos.rewards.%s.uord", d.ValidatorAddress), Value: rw.UOrd},
									{Key: fmt.Sprintf("ordos.rewards.%s.scor", d.ValidatorAddress), Value: rw.SCOR},
									{Key: fmt.Sprintf("ordos.rewards.%s.sord", d.ValidatorAddress), Value: rw.SORD},
									{Key: fmt.Sprintf("ordos.rewards.%s.shar", d.ValidatorAddress), Value: rw.SHAR},
									{Key: fmt.Sprintf("ordos.rewards.%s.satr", d.ValidatorAddress), Value: rw.SATR},
									{Key: fmt.Sprintf("ordos.rewards.%s.claim", d.ValidatorAddress), Value: rw.Claim},
								},
							},
						}
						db.UpdateOne(filter, update)
					}
				}
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			log.Fatal(err)
		}

		end := time.Since(start)
		if end > (time.Second * 90) {
			height += 150
		} else if end > (time.Second * 60) {
			height += 100
		} else if end > (time.Second * 30) {
			height += 50
		} else {
			height += 10
		}

		fmt.Printf("chain : %v height : %v  time : %v End!\n", chainCode, height, end)
	}
}
