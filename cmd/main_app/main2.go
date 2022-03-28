package main

import (
	"rut955_openwrt/src/wialonClient"
	"sync"
)

type AppConfigPaths struct {
	WialonClientConfig string
	DataServerCConfig  string
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	client := wialonClient.NewWialonClient()


	wg.Wait()
}
