package account

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/DonggyuLim/Alliance-Rank/utils"
	"github.com/cosmos/cosmos-sdk/types"
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

type Account struct {
	Address   string `json:"address"`
	Atreides  Chain  `json:"atreides"`
	Harkonnen Chain  `json:"harkonnen"`
	Corrino   Chain  `json:"corrino"`
	Ordos     Chain  `json:"ordos"`
	Total     Total  `json:"total"`
}
type Chain struct {
	Address string            `json:"address"`
	Rewards map[string]Reward `json:"rewards"` //key = validator Address
	Total   ChainTotal        `json:"total"`
}

type Reward struct {
	UAtr  int `json:"uatr"`
	UHar  int `json:"uhar"`
	UOrd  int `json:"uord"`
	UCor  int `json:"ucor"`
	SCOR  int `json:"scor"`
	SORD  int `json:"sord"`
	SHAR  int `json:"shar"`
	SATR  int `json:"satr"`
	Claim `json:"claim"`
}

func (r *Reward) Add(c []types.Coin) {
	for _, re := range c {
		switch re.Denom {
		case sATR:
			r.SATR += int(re.Amount.Int64())
		case sHAR:
			r.SHAR += int(re.Amount.Int64())
		case sCOR:
			r.SCOR += int(re.Amount.Int64())
		case sORD:
			r.SORD += int(re.Amount.Int64())
		case uatr:
			r.UAtr += int(re.Amount.Int64())
		case uhar:
			r.UHar += int(re.Amount.Int64())
		case ucor:
			r.UCor += int(re.Amount.Int64())
		case uord:
			r.UOrd += int(re.Amount.Int64())
		}
	}
}

type Claim struct {
	UAtr int `json:"uatr"`
	UCor int `json:"ucor"`
	UHar int `json:"uhar"`
	UOrd int `json:"uord"`
	SCOR int `json:"scor"`
	SORD int `json:"sord"`
	SHAR int `json:"shar"`
	SATR int `json:"satr"`
}

func (c *Claim) Add(coin []types.Coin) {
	for _, re := range coin {
		switch re.Denom {
		case sATR:
			c.SATR += int(re.Amount.Int64())
		case sHAR:
			c.SHAR += int(re.Amount.Int64())
		case sCOR:
			c.SCOR += int(re.Amount.Int64())
		case sORD:
			c.SORD += int(re.Amount.Int64())
		case uatr:
			c.UAtr += int(re.Amount.Int64())
		case uhar:
			c.UHar += int(re.Amount.Int64())
		case ucor:
			c.UCor += int(re.Amount.Int64())
		case uord:
			c.UOrd += int(re.Amount.Int64())
		}
	}
}

type Total struct {
	UAtr  int ` json:"uatr"`
	UCor  int ` json:"ucor"`
	UHar  int ` json:"uhar"`
	UOrd  int ` json:"uord"`
	SCOR  int ` json:"scor"`
	SORD  int ` json:"sord"`
	SHAR  int `json:"shar"`
	SATR  int `json:"satr"`
	Total int `json:"total"`
}

type ChainTotal struct {
	UAtr int `json:"uatr"`
	UCor int `json:"ucor"`
	UHar int `json:"uhar"`
	UOrd int `json:"uord"`
	SCOR int `json:"scor"`
	SORD int `json:"sord"`
	SHAR int `json:"shar"`
	SATR int `json:"satr"`
	// Total uint `json:"total"`
}

func (a *Account) SetAccount(address, validator string, reward Reward, chainCode int) {

	m1 := make(map[string]Reward)
	m2 := make(map[string]Reward)
	m3 := make(map[string]Reward)
	m4 := make(map[string]Reward)
	a.Address = address
	a.Atreides.Rewards = m1
	a.Harkonnen.Rewards = m2
	a.Corrino.Rewards = m3
	a.Ordos.Rewards = m4

	a.Atreides.Address = utils.MakeAddressPrefix(address, "atreides")
	a.Harkonnen.Address = utils.MakeAddressPrefix(address, "harkonnen")
	a.Corrino.Address = utils.MakeAddressPrefix(address, "corrino")
	a.Ordos.Address = utils.MakeAddressPrefix(address, "ordos")

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

func (a *Account) CalculateTotal() {

	ct := ChainTotal{
		UAtr: 0,
		UHar: 0,
		UCor: 0,
		UOrd: 0,
		SCOR: 0,
		SORD: 0,
		SHAR: 0,
		SATR: 0,
	}

	for _, v := range a.Atreides.Rewards {
		ct.UAtr += v.UAtr
		ct.UAtr += v.Claim.UAtr
		ct.SCOR += v.SCOR
		ct.SCOR += v.Claim.SCOR
		ct.SORD += v.SORD
		ct.SORD += v.Claim.SORD
		ct.SHAR += v.SHAR
		ct.SHAR += v.Claim.SHAR
		ct.SATR += v.SATR
		ct.SATR += v.Claim.SATR
	}

	a.Atreides.Total = ct

	for _, v := range a.Harkonnen.Rewards {
		ct.UHar += v.UHar
		ct.UHar += v.Claim.UHar
		ct.SCOR += v.SCOR
		ct.SCOR += v.Claim.SCOR
		ct.SORD += v.SORD
		ct.SORD += v.Claim.SORD
		ct.SHAR += v.SHAR
		ct.SHAR += v.Claim.SHAR
		ct.SATR += v.SATR
		ct.SATR += v.Claim.SATR
	}
	//claim reward +

	a.Harkonnen.Total = ct
	// a.Total = a.Total+ a.Harkonnen.Total.NativeTotal)+ a.Harkonnen.Total.SCOR)+ a.Harkonnen.Total.SORD)
	for _, v := range a.Corrino.Rewards {
		ct.UCor += v.UCor
		ct.UCor += v.Claim.UCor
		ct.SCOR += v.SCOR
		ct.SCOR += v.Claim.SCOR
		ct.SORD += v.SORD
		ct.SORD += v.Claim.SORD
		ct.SHAR += v.SHAR
		ct.SHAR += v.Claim.SHAR
		ct.SATR += v.SATR
		ct.SATR += v.Claim.SATR
	}
	//claim reward +

	a.Corrino.Total = ct

	for _, v := range a.Ordos.Rewards {
		ct.UOrd += v.UOrd
		ct.UOrd += v.Claim.UOrd
		ct.SCOR += v.SCOR
		ct.SCOR += v.Claim.SCOR
		ct.SORD += v.SORD
		ct.SORD += v.Claim.SORD
		ct.SHAR += v.SHAR
		ct.SHAR += v.Claim.SHAR
		ct.SATR += v.SATR
		ct.SATR += v.Claim.SATR
	}
	//claim reward +

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
	a.Total.SHAR = a.Atreides.Total.SHAR + a.Harkonnen.Total.SHAR + a.Corrino.Total.SHAR + a.Ordos.Total.SHAR
	a.Total.SATR = a.Atreides.Total.SATR + a.Harkonnen.Total.SATR + a.Corrino.Total.SATR + a.Ordos.Total.SATR
	a.Total.Total = a.Total.UAtr + a.Total.UHar + a.Total.UCor + a.Total.UOrd + a.Total.SCOR + a.Total.SORD + a.Total.SHAR + a.Total.SATR
}

func (r Reward) EncodeJson() string {
	bytes, err := json.MarshalIndent(r, "", "   ")
	utils.PanicError(err)
	return string(bytes)
}

// func (r Reward) GetReward(endpint string, chainCode int) {
// 	client := req.C.R()

// }
