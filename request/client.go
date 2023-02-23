package request

import "github.com/imroc/req/v3"

func GetClient() *req.Request {
	client := req.C().R()
	return client
}
