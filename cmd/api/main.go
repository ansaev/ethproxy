package main

import (
	"ethproxy/internal/api"
	"ethproxy/internal/ethfinder"
	"ethproxy/internal/txfinder"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	envListenAddress     = "SERVER_ADDRESS"
	envETHGatewayAddress = "ETH_GATEWAY_ADDRESS"
	envETHTimeout        = "ETH_TIMEOUT" // in seconds

	defaultTimeout int64 = 30
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// parse and set timeout setting for eth requests
	var timeout int64
	timeoutStr := os.Getenv(envETHTimeout)
	if timeoutStr == "" {
		timeout = defaultTimeout
	} else {
		var errParsetimeout error
		timeout, errParsetimeout = strconv.ParseInt(timeoutStr, 10, 64)
		if errParsetimeout != nil {
			log.Fatalf("Invalid timeout value: \"%s\", unable to parse: %s\n",
				timeoutStr, errParsetimeout.Error())
		}
	}

	ethBlockService := ethfinder.New(os.Getenv(envETHGatewayAddress),
		&http.Client{Timeout: time.Duration(timeout) * time.Second})

	txService := txfinder.New(ethBlockService)

	apiService := api.New(os.Getenv(envListenAddress), txService)
	err := apiService.ListenAndServe()
	if err != nil {
		log.Fatal("api error: ", err)
	}
}
