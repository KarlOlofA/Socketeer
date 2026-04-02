package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/joho/godotenv"
	auth "socketeer.github.com/internal/auth"
)

const (
	HOST   = "localhost"
	PORT   = "8080"
	METHOD = "tcp"
)

type Connection struct {
	user       auth.User
	connection net.Conn
}

var connections map[string]Connection

func main() {
	godotenv.Load()
	connections = make(map[string]Connection)
	ln, err := net.Listen(METHOD, fmt.Sprintf("%s:%s", HOST, PORT))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ln.Close()
	fmt.Printf("Listening to port %s.\n", PORT)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Print("TCP accept failed.\n")
			return
		}
		go handleConnection(conn)
	}
}

type Data struct {
	Msg    string    `json:"msg"`
	ApiKey string    `json:"apiKey"`
	Brush  BrushData `json:"brushdata"`
}

type BrushData struct {
	Size  float32 `json:"size"`
	Color struct {
		R float32 `json:"r"`
		G float32 `json:"g"`
		B float32 `json:"b"`
		A float32 `json:"a"`
	} `json:"color"`
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Failed to connect %s\n", conn.RemoteAddr())
			return
		}
		var d Data
		if err := json.Unmarshal(buffer[:n], &d); err != nil {

			responseStr := "Failed to parse message to json."
			_, err = conn.Write([]byte(responseStr))
			return
		}
		address := conn.RemoteAddr().String()
		connection, ok := connections[address]
		if !ok {
			fmt.Printf("Create user at %s\n", address)
			connections[address] = Connection{
				user: auth.User{
					IpAddress: address,
					IsAuthed:  false,
				},
				connection: conn,
			}
		}

		if !connection.user.IsAuthed {
			_, err := connection.user.ValidateApiKey(d.ApiKey)
			if err != nil {
				fmt.Printf("Failed to authenticate user at %s\n", address)
				return
			}
			connection.user.IsAuthed = true
		}

		distributeMessage(connection, d)

		fmt.Printf("%v : Conn -> \nMsg: %s\nBrush size: %v\nBrush color: %v\n", conn.RemoteAddr(), d.Msg, d.Brush.Size, d.Brush.Color)
		time := time.Now().Format(time.ANSIC)
		responseStr := fmt.Sprintf("Valid client message recieved at %v", time)
		_, err = conn.Write([]byte(responseStr))
	}
}

func distributeMessage(distConn Connection, data Data) {
	for _, conn := range connections {
		if conn.user.IpAddress == distConn.user.IpAddress {
			continue
		}

		json, err := json.Marshal(data)

		if err != nil {
			continue
		}

		conn.connection.Write(json)
	}
}

func closeConnection(conn net.Conn) {
	if conn == nil {
		return
	}

	conn.Close()
}
