package account

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/DonggyuLim/Alliance-Rank/utils"
)

type Account struct {
	Address   string `bson:"address" json:"address"`
	Atreides  Chain  `bson:"atreides" json:"atreides"`
	Harkonnen Chain  `bson:"harkonnen" json:"harkonnen"`
	Corrino   Chain  `bson:"corrino" json:"corrino"`
	Ordos     Chain  `bson:"ordos" json:"ordos"`
	Total     Total  `bson:"total" json:"total"`
}
type Chain struct {
	Address string            `bson:"address" json:"address"`
	Rewards map[string]Reward `bson:"reward" json:"rewards"` //key = validator Address
	Claim   Claim             `bson:"claim" json:"claim"`
	Total   ChainTotal        `bson:"total" json:"total"`
}

type Reward struct {
	LastHeight uint `bson:"last_height" json:"last_height"`
	UAtr       uint `bson:"uatr" json:"uatr"`
	UHar       uint `bson:"uhar" json:"uhar"`
	UOrd       uint `bson:"uord" json:"uord"`
	UCor       uint `bson:"ucor" json:"ucor"`
	SCOR       uint `bson:"scor" json:"scor"`
	SORD       uint `bson:"sord" json:"sord"`
}

type Claim struct {
	UAtr uint `json:"uatr"`
	UCor uint `json:"ucor"`
	UHar uint `json:"uhar"`
	UOrd uint `json:"uord"`
	SCOR uint `json:"scor"`
	SORD uint `json:"sord"`
}
type Total struct {
	UAtr  uint `json:"uatr"`
	UCor  uint `json:"ucor"`
	UHar  uint `json:"uhar"`
	UOrd  uint `json:"uord"`
	SCOR  uint `json:"scor"`
	SORD  uint `json:"sord"`
	Total uint `json:"total"`
}

type ChainTotal struct {
	UAtr uint `json:"uatr"`
	UCor uint `json:"ucor"`
	UHar uint `json:"uhar"`
	UOrd uint `json:"uord"`
	SCOR uint `json:"scor"`
	SORD uint `json:"sord"`
	// Total uint `json:"total"`
}

func (a *Account) SetAccount(address, validator string, reward Reward, chainCode int) {

	m1 := make(map[string]Reward)
	m2 := make(map[string]Reward)
	m3 := make(map[string]Reward)
	m4 := make(map[string]Reward)
	a.Address = utils.MakeAddress(address)
	a.Atreides.Rewards = m1
	a.Harkonnen.Rewards = m2
	a.Corrino.Rewards = m3
	a.Ordos.Rewards = m4
	switch chainCode {
	case 0:
		a.Atreides.Rewards[validator] = reward
	case 1:
		a.Harkonnen.Rewards[validator] = reward
	case 2:
		a.Corrino.Rewards[validator] = reward
	case 3:
		a.Ordos.Rewards[validator] = reward
	}

}

func (a Account) EncodeByte() []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	utils.PanicError(encoder.Encode(a))
	return aBuffer.Bytes()
}

func (a *Account) FromBytes(data []byte) {
	encoder := gob.NewDecoder(bytes.NewReader(data))
	utils.PanicError(encoder.Decode(&a))
}

func (c *Chain) UpdateClaimAndReward(
	delegator,
	validator string,
	r Reward,
	chainCode int) {

	switch chainCode {
	case 0:
		c.Address = delegator
		origin := c.Rewards[validator]
		if origin.UAtr > r.UAtr {
			claim := origin.UAtr - r.UAtr
			c.Claim.UAtr =
				c.Claim.UAtr + claim

		}
		c.Rewards[validator] = r
	case 1:
		c.Address = delegator
		origin := c.Rewards[validator]
		if origin.UHar > r.UHar {
			claim := origin.UHar - r.UHar
			c.Claim.UHar =
				c.Claim.UHar + claim
		}
		c.Rewards[validator] = r
	case 2:
		c.Address = delegator
		origin := c.Rewards[validator]
		if origin.UCor > r.UCor {
			claim := origin.UCor - r.UCor
			c.Claim.UCor =
				c.Claim.UCor + claim
		}

		c.Rewards[validator] = r
	case 3:
		c.Address = delegator
		origin := c.Rewards[validator]
		if origin.UOrd > r.UOrd {
			claim := origin.UOrd - r.UOrd
			c.Claim.UOrd =
				c.Claim.UOrd + claim
		}
		c.Rewards[validator] = r
	}
}

func (c *Chain) UpdateUndelegate(chainCode, height int) {
	deleteKey := []string{}
	h := uint(height)
	switch chainCode {
	case 0:
		for k, v := range c.Rewards {
			if v.LastHeight < h {
				c.Claim.UAtr =
					c.Claim.UAtr + v.UAtr
				c.Claim.SCOR =
					c.Claim.SCOR + v.SCOR
				c.Claim.SORD =
					c.Claim.SORD + v.SORD
				deleteKey = append(deleteKey, k)
			}
		}

	case 1:
		for k, v := range c.Rewards {
			if v.LastHeight < h {
				c.Claim.UHar =
					c.Claim.UHar + v.UHar
				c.Claim.SCOR =
					c.Claim.SCOR + v.SCOR
				c.Claim.SORD =
					c.Claim.SORD + v.SORD
				deleteKey = append(deleteKey, k)
			}
		}
	case 2:
		for k, v := range c.Rewards {
			if v.LastHeight < h {
				c.Claim.UCor =
					c.Claim.UCor + v.UCor
				c.Claim.SCOR =
					c.Claim.SCOR + v.SCOR
				c.Claim.SORD =
					c.Claim.SORD + v.SORD
				deleteKey = append(deleteKey, k)
			}
		}
	case 3:
		for k, v := range c.Rewards {
			if v.LastHeight < h {
				c.Claim.UOrd =
					c.Claim.UOrd + v.UOrd
				c.Claim.SCOR =
					c.Claim.SCOR + v.SCOR
				c.Claim.SORD =
					c.Claim.SORD + v.SORD
				deleteKey = append(deleteKey, k)
			}
		}
	}
	for _, key := range deleteKey {
		delete(c.Rewards, key)
	}
}

func (a *Account) CalculateTotal() {

	ct := ChainTotal{
		UAtr: 0,
		UHar: 0,
		UCor: 0,
		UOrd: 0,
		SCOR: 0,
		SORD: 0,
	}

	for _, v := range a.Atreides.Rewards {
		ct.UAtr += v.UAtr
		ct.SCOR += v.SCOR
		ct.SORD += v.SORD
	}
	//claim reward +
	ct.UAtr += a.Atreides.Claim.UAtr
	ct.SCOR += a.Atreides.Claim.SCOR
	ct.SORD += a.Atreides.Claim.SORD
	a.Atreides.Total = ct

	for _, v := range a.Harkonnen.Rewards {
		ct.UHar += v.UHar
		ct.SCOR += v.SCOR
		ct.SORD += v.SORD
	}
	//claim reward +
	ct.UHar += a.Harkonnen.Claim.UHar
	ct.SCOR += a.Harkonnen.Claim.SCOR
	ct.SORD += a.Harkonnen.Claim.SORD

	a.Harkonnen.Total = ct
	// a.Total = a.Total+ a.Harkonnen.Total.NativeTotal)+ a.Harkonnen.Total.SCOR)+ a.Harkonnen.Total.SORD)
	for _, v := range a.Corrino.Rewards {
		ct.UCor += v.UCor
		ct.SCOR += v.SCOR
		ct.SORD += v.SORD
	}
	//claim reward +
	ct.UCor += a.Corrino.Claim.UCor
	ct.SCOR += a.Corrino.Claim.SCOR
	ct.SORD += a.Corrino.Claim.SORD

	a.Corrino.Total = ct

	for _, v := range a.Ordos.Rewards {
		ct.UOrd += v.UOrd

		ct.SCOR += v.SCOR

		ct.SORD += v.SORD

	}
	//claim reward +
	ct.UOrd += a.Ordos.Claim.UOrd
	ct.SCOR += a.Ordos.Claim.SCOR
	ct.SORD += a.Ordos.Claim.SORD

	a.Ordos.Total = ct

	// a.Total = a.Total+ a.Ordos.Total.NativeTotal)+ a.Ordos.Total.SCOR)+ a.Ordos.Total.SORD)

	a.Total = Total{}
	//calculate NativeTotal

	a.Total.UAtr = a.Atreides.Total.UAtr
	a.Total.UHar = a.Harkonnen.Total.UHar
	a.Total.UCor = a.Corrino.Total.UCor
	a.Total.UOrd = a.Ordos.Total.UOrd
	//calculate SCOR Total
	a.Total.SCOR = a.Atreides.Total.SCOR + a.Harkonnen.Total.SCOR + a.Corrino.Total.SCOR + a.Ordos.Total.SCOR
	///calculate SORD Total
	a.Total.SORD = a.Atreides.Total.SORD + a.Harkonnen.Total.SORD + a.Corrino.Total.SORD + a.Ordos.Total.SORD
	a.Total.Total = a.Total.UAtr + a.Total.UHar + a.Total.UCor + a.Total.UOrd + a.Total.SCOR + a.Total.SORD
}

func (r Reward) EncodeJson() string {
	bytes, err := json.MarshalIndent(r, "", "   ")
	utils.PanicError(err)
	return string(bytes)
}

// func (r Reward) GetReward(endpint string, chainCode int) {
// 	client := req.C.R()

// }
