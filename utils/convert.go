package utils

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/shopspring/decimal"
)

func accountDecode(account string) []byte {

	_, decoded, err := bech32.Decode(account)
	HandleErr("Decode Err", err)
	return decoded
}

func accountEncode(prefix string, account []byte) (address string) {
	encoded, err := bech32.Encode(prefix, []byte(account))
	HandleErr("Encode Err", err)

	return encoded
}

func MakeKey(account string) string {
	bytes := accountDecode(account)
	return accountEncode("atreides", bytes)
}

func MakeAddressPrefix(account, prefix string) string {
	bytes := accountDecode(account)
	return accountEncode(prefix, bytes)
}

func PrettyJson(data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "   ")
	PanicError(err)
	fmt.Println(string(bytes))
	// return bytes
}

func ChangeDeciaml(a string) decimal.Decimal {
	return decimal.RequireFromString(a)
}

func DecimalAddString(a decimal.Decimal, b string) decimal.Decimal {

	return a.Add(ChangeDeciaml(b))
}
