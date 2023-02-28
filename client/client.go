package client

import (
	"sync"

	alliancemoduletypes "github.com/terra-money/alliance/x/alliance/types"
	"google.golang.org/grpc"
)

var once sync.Once
var client alliancemoduletypes.QueryClient

func QueryClient() alliancemoduletypes.QueryClient {
	once.Do(func() {
		conn, _ := grpc.Dial("localhost:19090",
			grpc.WithInsecure(),
			grpc.WithBlock(),
		)
		client = alliancemoduletypes.NewQueryClient(conn)
	})

	return client
}
