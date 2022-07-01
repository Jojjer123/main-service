package main

import (
	northboundinterface "main-service/pkg/northbound-interface"
	store "main-service/pkg/store-wrapper"
)

func main() {
	store.CreateStores()
	northboundinterface.StartServer()
}
