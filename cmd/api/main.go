package main

import (
	"ethproxy/internal/api"
	"ethproxy/internal/ethfinder"
	"ethproxy/internal/txfinder"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	ListenAddressName     = "SERVER_ADDRESS"
	ETHGatewayAddressName = "ETH_GATEWAY_ADDRESS"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}
	ethBlockService := ethfinder.New(os.Getenv(ETHGatewayAddressName), http.DefaultClient)

	txService := txfinder.New(ethBlockService)

	apiService := api.New(os.Getenv(ListenAddressName), txService)
	err := apiService.ListenAndServe()
	if err != nil {
		log.Fatal("api error: ", err)
	}

}
