package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	pool "socketeer.github.com/internal/network"
	network "socketeer.github.com/internal/types/network"
)

const (
	HOST   = "localhost"
	PORT   = "8080"
	METHOD = "tcp"
)

var connections map[string]network.Connection
var conns map[string]net.Conn

func main() {
	godotenv.Load()

	listener, err := net.Listen(METHOD, fmt.Sprintf("%s:%s", HOST, PORT))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listener.Close()

	connectionPool := pool.NewConnectionManager(10).AddMiddleware(func(net.Conn) error {
		fmt.Printf("Womp\n")
		return nil
	}).AddMiddleware(func(net.Conn) error {
		fmt.Printf("Gomp\n")
		return nil
	}).AddMiddleware(func(net.Conn) error {
		fmt.Printf("erromp\n")
		return errors.New("Fail")
	}).AddMiddleware(func(net.Conn) error {
		fmt.Printf("Finality\n")
		return nil
	})

	connectionPool.AssignTCPListener(listener)

	fmt.Printf("Listening to port %s.\n", PORT)
	go connectionPool.Run()

	connectionPool.ProcessConnections()
}
