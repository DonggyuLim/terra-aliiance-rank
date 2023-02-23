package main

import (
	"sync"

	"github.com/DonggyuLim/Alliance-Rank/data"
	"github.com/DonggyuLim/Alliance-Rank/db"
	"github.com/DonggyuLim/Alliance-Rank/rest"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	db.Connect()
	defer db.Close()

	go data.Main(wg)

	go rest.Start(wg)
	wg.Wait()
}
