package client

import (
	"sync"

	alliancemoduletypes "github.com/terra-money/alliance/x/alliance/types"
	"google.golang.org/grpc"
)

var once sync.Once

func QueryClient(chainCode int) alliancemoduletypes.QueryClient {
	var endpoint string
	switch chainCode {
	case 0:
		endpoint = "localhost:9090"
	case 1:
		endpoint = "localhost:19090"
	case 2:
		endpoint = "localhost:29090"
	case 3:
		endpoint = "localhost:39090"
	}
	conn, _ := grpc.Dial(endpoint,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	client := alliancemoduletypes.NewQueryClient(conn)

	return client
}
