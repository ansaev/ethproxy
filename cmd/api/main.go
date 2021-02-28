package main

import (
	"ethproxy/internal/api"
	"ethproxy/internal/cachefinder"
	"ethproxy/internal/ethfinder"
	"ethproxy/internal/txfinder"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"

	"github.com/joho/godotenv"
)

const (
	envListenAddress     = "SERVER_ADDRESS"
	envETHGatewayAddress = "ETH_GATEWAY_ADDRESS"
	envETHTimeout        = "ETH_TIMEOUT" // in seconds
	envRedisDB           = "REDIS_DB"
	envRedisHost         = "REDIS_HOST"
	envRedisPassword     = "REDIS_PASSWORD"

	defaultEthTimeout int64 = 30
	redisDialTimeout        = 10 * time.Second
	redisReadTimeout        = 30 * time.Second
	redisWriteTimeout       = 30 * time.Second
	redisPoolSize           = 10
	redisPoolTimeout        = 10 * time.Second

	EthereumName = "Ethereum"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// parse and set timeout setting for eth requests
	var timeout int64
	timeoutStr := os.Getenv(envETHTimeout)
	if timeoutStr == "" {
		timeout = defaultEthTimeout
	} else {
		var errParsetimeout error
		timeout, errParsetimeout = strconv.ParseInt(timeoutStr, 10, 64)
		if errParsetimeout != nil {
			log.Fatalf("Invalid timeout value: \"%s\", unable to parse: %s\n",
				timeoutStr, errParsetimeout.Error())
		}
	}

	redisDB, err := strconv.Atoi(os.Getenv(envRedisDB))
	if err != nil {
		log.Fatalf("unable to parse redis Db number: %v\n", err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:         os.Getenv(envRedisHost),
		Password:     os.Getenv(envRedisPassword),
		DB:           redisDB,
		DialTimeout:  redisDialTimeout,
		ReadTimeout:  redisReadTimeout,
		WriteTimeout: redisWriteTimeout,
		PoolSize:     redisPoolSize,
		PoolTimeout:  redisPoolTimeout,
	})

	ethBlockService := ethfinder.New(os.Getenv(envETHGatewayAddress),
		&http.Client{Timeout: time.Duration(timeout) * time.Second})

	cachefinder.New(ethBlockService, redisClient, EthereumName)

	txService := txfinder.New(ethBlockService)

	apiService := api.New(os.Getenv(envListenAddress), txService)
	errServe := apiService.ListenAndServe()
	if errServe != nil {
		log.Fatalf("api server error: %v\n", errServe)
	}
}
