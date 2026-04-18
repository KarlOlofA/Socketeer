package main

import (
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	auth "socketeer.github.com/internal/auth"
	types "socketeer.github.com/internal/types"
)

const (
	HOST   = "194.47.156.73"
	PORT   = "8080"
	METHOD = "tcp"
)

type Connection struct {
	user       auth.User
	connection net.Conn
}

var connections map[string]Connection
var conns map[string]net.Conn

func main() {
	godotenv.Load()

	connections = make(map[string]Connection)

	listener, err := net.Listen(METHOD, fmt.Sprintf("%s:%s", HOST, PORT))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening to port %s.\n", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Print("TCP accept failed.\n")
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			closeConnection(conn, fmt.Sprintf("Failed to connect %s\n", conn.RemoteAddr()))
			break
		}
		/*
			var packets []types.Packet
			var packet types.Packet
			packet.BuildFromByteSlice(buffer[:n])
			packets = append(packets, packet)
		*/
		address := conn.RemoteAddr().String()
		if _, ok := conns[address]; !ok {
			conns[address] = conn
		}
		distributePacketConn(conn, buffer[:n])

		continue
		/*
			for _, packet := range packets {

							if !ok {
								fmt.Printf("Create user at %s\n", address)
								connections[address] = Connection{
									user: auth.User{
										IpAddress: address,
										IsAuthed:  false,
									},
									connection: conn,
								}
								connection = connections[address]
							}
							packet.Key = "1234"
							if !connection.user.IsAuthed {
								_, err := connection.user.ValidateApiKey(packet.Key)
								if err != nil {
									closeConnection(conn, fmt.Sprintf("Failed to authenticate user at %s\n", address))
									return
								}
								connection.user.IsAuthed = true
							}
							if _, ok := connections[connection.connection.RemoteAddr().String()]; !ok {
								go distributeWelcome(conn)
								continue
							}*/
		//go distributePacket(&connection, packet)
	}

}

func distributePacketConn(distConn net.Conn, packet []byte) {
	for _, conn := range conns {
		if conn.RemoteAddr().String() == distConn.RemoteAddr().String() {
			continue
		}

		conn.Write(packet)
	}
}
func distributePacket(distConn *Connection, packet types.Packet) {
	for _, conn := range connections {
		if conn.user.IpAddress == distConn.user.IpAddress {
			continue
		}

		conn.connection.Write(packet.Data)
	}
}

func distributeWelcome(distConn net.Conn) {
	if len(connections) <= 0 {
		connections[distConn.RemoteAddr().String()] = Connection{
			connection: distConn,
		}
	} else if _, ok := connections[distConn.RemoteAddr().String()]; !ok {
		connections[distConn.RemoteAddr().String()] = Connection{
			connection: distConn,
		}
	}

	for _, conn := range connections {
		if conn.user.IpAddress == distConn.RemoteAddr().String() {
			conn.connection.Write(fmt.Appendf(nil, "You (%s) connected to the server\n", distConn.RemoteAddr().String()))
			continue
		}
		conn.connection.Write(fmt.Appendf(nil, "%s connected to the session\n", distConn.RemoteAddr().String()))
	}
}

func closeConnection(conn net.Conn, msg string) {
	if conn == nil {
		return
	}

	fmt.Printf("Closed connection -> %s\n", msg)

	go conn.Write([]byte(msg))

	conn.Close()
}
