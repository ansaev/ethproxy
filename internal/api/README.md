# ethproxy
this service proxy requests to https://developers.cloudflare.com/distributed-web/ethereum-gateway

Also it requests to blocks in cache for blockCachingTime (2 hour for now). 
It is not cache requests to latest 6 blocks.

## deployment
1. copy file ".env.example" as ".env" and edit configuration
2. install golang dependencies "$ go mod vendor"
3. run redis locally and edit .env due to it's settigns
4. run application by "$ go run ethproxy/cmd/api" (when you are in the directory with source code)
