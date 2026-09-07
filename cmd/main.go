package main

import (
	"fmt"
	"log"
	"os"

	pool "github.com/KarlOlofA/socketeer/internal/network"
	authMw "github.com/KarlOlofA/socketeer/internal/network/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	args := os.Args

	settings := pool.TcpServerSettings{
		Host:               "192.168.1.41",
		Port:               "8080",
		Method:             "tcp",
		Key:                "1234",
		ConnectionPoolSize: 10,
	}
	if len(args) > 1 && len(args[1]) > 0 {
		settings.Host = args[1]
	}

	server := pool.NewTcpServer(settings)
	if server == nil {
		log.Fatal("Failed to construct tcp server.")
		return
	}
	defer server.Close()

	auth := authMw.AuthMiddleware(settings.Key)

	mwChain := pool.MiddlewareChain(auth)
	server.AddMiddleware(mwChain)

	fmt.Printf("Listening to %s:%s.\n", settings.Host, settings.Port)
	go server.Run()

	server.ProcessConnections()
}
