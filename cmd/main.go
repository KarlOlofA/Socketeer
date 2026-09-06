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

func main() {
	godotenv.Load()

	args := os.Args

	settings := pool.TcpServerSettings{
		Host:               "192.168.1.147",
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

	server.AddMiddleware(func(c net.Conn) (net.Conn, error) {

		p := network.Packet{}
		var buffer []byte = make([]byte, 1024)
		_, err := c.Read(buffer)
		if err != nil {
			return nil, fmt.Errorf("Failed to read packet from buffer: %v", err)
		}
		err = p.FromByteSlice(buffer)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse buffer to packet struct: %v", err)
		}
		fmt.Printf("Key -> %s\n", p.Key)
		if p.Key != settings.Key {
			return nil, fmt.Errorf("Failed to validate key")
		}
		return c, nil
	})

	fmt.Printf("Listening to %s:%s.\n", settings.Host, settings.Port)
	go server.Run()

	server.ProcessConnections()
}
