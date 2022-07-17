package main

import (
	"github.com/IoanaAdrian/HubA/src/api"
	"github.com/IoanaAdrian/HubA/src/api/services"
	"net/http"
)

func main() {
	services.Init()
	panic(http.ListenAndServe(":8080", http.HandlerFunc(api.Serve)))
}
