package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	pool "socketeer.github.com/internal/network"
	"socketeer.github.com/internal/types/network"
)

type HostData struct {
	Host   string
	Port   string
	Method string
	Key    string
}

func main() {
	godotenv.Load()

	h := HostData{
		Host:   "192.168.1.147",
		Port:   "8080",
		Method: "tcp",
		Key:    "1234",
	}

	args := os.Args

	if len(args) > 1 && len(args[1]) > 0 {
		h.Host = args[1]
	}

	listener, err := net.Listen(h.Method, fmt.Sprintf("%s:%s", h.Host, h.Port))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listener.Close()

	connectionPool := pool.NewConnectionManager(10)

	connectionPool.AddMiddleware(func(c net.Conn) (net.Conn, error) {
		return c, nil
		p := network.Packet{}
		var buffer []byte
		c.Read(buffer)
		p.FromByteSlice(buffer)
		if p.Key != h.Key {
			return nil, fmt.Errorf("Failed to validate key.")
		}
		return c, nil
	})

	connectionPool.AssignTCPListener(listener)

	fmt.Printf("Listening to %s:%s.\n", h.Host, h.Port)
	go connectionPool.Run()

	connectionPool.ProcessConnections()
}
