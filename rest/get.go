package rest

import (
	"fmt"

	"github.com/DonggyuLim/Alliance-Rank/db"
	"github.com/DonggyuLim/Alliance-Rank/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

type ToTalResponse struct {
	Address string `json:"address"`
	UAtr    string `json:"uatr"`
	UCor    string `json:"ucor"`
	UHar    string `json:"uhar"`
	UOrd    string `json:"uord"`
	SCOR    string `json:"scor"`
	SORD    string `json:"sord"`
	SHAR    string `json:"shar"`
	SATR    string `json:"satr"`
	Total   string `json:"total"`
}

func Root(c *gin.Context) {
	fmt.Println("ROOT")
	list, err := db.Find("", "", "total.total", 100)
	// fmt.Println(list)
	if err != nil {
		fmt.Println(err)
		c.String(404, err.Error())
		return
	}
	var res []ToTalResponse
	for _, el := range list {
		total := ToTalResponse{
			Address: el.Address,
			UAtr:    fmt.Sprintf("%v", el.Total.UAtr),
			UCor:    fmt.Sprintf("%v", el.Total.UCor),
			UHar:    fmt.Sprintf("%v", el.Total.UHar),
			UOrd:    fmt.Sprintf("%v", el.Total.UOrd),
			SCOR:    fmt.Sprintf("%v", el.Total.SCOR),
			SORD:    fmt.Sprintf("%v", el.Total.SORD),
			SATR:    fmt.Sprintf("%v", el.Total.SATR),
			SHAR:    fmt.Sprintf("%v", el.Total.SHAR),
			Total:   fmt.Sprintf("%v", el.Total.Total),
		}
		res = append(res, total)
	}

	c.JSON(200, res)
}

type UAtrResponse struct {
	Address string `json:"address"`
	UAtr    string `json:"amount"`
}

func UatrRank(c *gin.Context) {
	list, err := db.Find("", "", "total.uatr", 100)
	var res []UAtrResponse
	for _, el := range list {
		atr := UAtrResponse{
			Address: el.Atreides.Address,
			UAtr:    fmt.Sprintf("%v", el.Total.UAtr),
		}
		res = append(res, atr)
	}
	if err != nil {
		fmt.Println(err)
		c.String(404, err.Error())
		return
	}
	c.JSON(200, res)
}

type UharResponse struct {
	Address string `json:"address"`
	UHar    string `json:"amount"`
}

func UHarRank(c *gin.Context) {
	list, err := db.Find("", "", "total.uhar", 100)
	var res []UharResponse
	for _, el := range list {
		uhar := UharResponse{
			Address: el.Harkonnen.Address,
			UHar:    fmt.Sprintf("%v", el.Total.UHar),
		}
		res = append(res, uhar)
	}
	if err != nil {
		fmt.Println(err)
		c.String(404, err.Error())
		return
	}
	c.JSON(200, res)
}

type UCorResponse struct {
	Address string `json:"address"`
	UCor    string `json:"amount"`
}

func UCorRank(c *gin.Context) {
	list, err := db.Find("", "", "total.ucor", 100)
	var res []UCorResponse
	for _, el := range list {
		ucor := UCorResponse{
			Address: el.Corrino.Address,
			UCor:    fmt.Sprintf("%v", el.Total.UCor),
		}
		res = append(res, ucor)
	}
	if err != nil {
		fmt.Println(err)
		c.String(404, err.Error())
		return
	}
	c.JSON(200, res)
}

type UOrdResponse struct {
	Address string `json:"address"`
	UOrd    string `json:"amount"`
}

func UOrdRank(c *gin.Context) {
	list, err := db.Find("", "", "total.uord", 100)
	var res []UOrdResponse
	for _, el := range list {
		uord := UOrdResponse{
			Address: el.Ordos.Address,
			UOrd:    fmt.Sprintf("%v", el.Total.UOrd),
		}
		res = append(res, uord)
	}
	if err != nil {
		fmt.Println(err)
		c.String(404, err.Error())
		return
	}
	c.JSON(200, res)
}

type ScorResponse struct {
	Address string `json:"address"`
	SCor    string `json:"amount"`
}

func SCorRank(c *gin.Context) {
	list, err := db.Find("", "", "total.scor", 100)
	var res []ScorResponse
	for _, el := range list {
		scor := ScorResponse{
			Address: el.Address,
			SCor:    fmt.Sprintf("%v", el.Total.SCOR),
		}
		res = append(res, scor)
	}
	if err != nil {
		fmt.Println(err)
		c.String(404, err.Error())
		return
	}
	c.JSON(200, res)
}

type SOrdResponse struct {
	Address string `json:"address"`
	Sord    string `json:"amount"`
}

func SOrdRank(c *gin.Context) {
	list, err := db.Find("", "", "total.sord", 100)
	var res []SOrdResponse
	for _, el := range list {
		sord := SOrdResponse{
			Address: el.Address,
			Sord:    fmt.Sprintf("%v", el.Total.SORD),
		}
		res = append(res, sord)
	}
	if err != nil {
		fmt.Println(err)
		c.String(404, err.Error())
		return
	}
	c.JSON(200, res)
}

type SAtrResponse struct {
	Address string `json:"address"`
	SATR    string `json:"amount"`
}

func SAtrRank(c *gin.Context) {
	list, err := db.Find("", "", "total.satr", 100)
	var res []SAtrResponse
	for _, el := range list {
		satr := SAtrResponse{
			Address: el.Address,
			SATR:    fmt.Sprintf("%v", el.Total.SATR),
		}
		res = append(res, satr)
	}
	if err != nil {
		fmt.Println(err)
		c.String(404, err.Error())
		return
	}
	c.JSON(200, res)
}

type SHARResponse struct {
	Address string `json:"address"`
	SHAR    string `json:"amount"`
}

func SHarRank(c *gin.Context) {
	list, err := db.Find("", "", "total.shar", 100)
	var res []SHARResponse
	for _, el := range list {
		uord := SHARResponse{
			Address: el.Address,
			SHAR:    fmt.Sprintf("%v", el.Total.SHAR),
		}
		res = append(res, uord)
	}
	if err != nil {
		fmt.Println(err)
		c.String(404, err.Error())
		return
	}
	c.JSON(200, res)
}

type AccountdResponse struct {
	Address string `json:"address"`
	UAtr    string `json:"uatr"`
	UHar    string `json:"uhar"`
	UCor    string `json:"ucor"`
	UOrd    string `json:"uord"`
	SCor    string `json:"scor"`
	SOrd    string `json:"sord"`
	SHar    string `json:"shar"`
	SAtr    string `json:"satr"`
	Total   string `json:"total"`
}

func UserReward(c *gin.Context) {

	address := c.Param("address")

	key := utils.MakeKey(address)
	filter := bson.D{{Key: "address", Value: key}}

	a, ok := db.FindOne(filter)
	switch ok {
	case nil:
		myReward := AccountdResponse{
			Address: a.Address,
			UAtr:    fmt.Sprintf("%v", a.Total.UAtr),
			UHar:    fmt.Sprintf("%v", a.Total.UHar),
			UCor:    fmt.Sprintf("%v", a.Total.UCor),
			UOrd:    fmt.Sprintf("%v", a.Total.UOrd),
			SCor:    fmt.Sprintf("%v", a.Total.SCOR),
			SOrd:    fmt.Sprintf("%v", a.Total.SORD),
			SAtr:    fmt.Sprintf("%v", a.Total.SATR),
			SHar:    fmt.Sprintf("%v", a.Total.SHAR),
			Total:   fmt.Sprintf("%v", a.Total.Total),
		}
		c.JSON(200, myReward)
	case mongo.ErrNoDocuments:
		// fmt.Println(address)
		c.String(404, "Invalid Address")
	}
}
