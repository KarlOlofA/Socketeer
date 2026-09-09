package connectionPool

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/KarlOlofA/socketeer/internal/types/network"
)

type TcpServer struct {
	TCPListener   net.Listener
	Connections   sync.Map
	mu            sync.Mutex
	addChannel    chan net.Conn
	removeChannel chan net.Conn
	Middleware    []func(net.Conn) (net.Conn, error)
}

type TcpServerSettings struct {
	Host               string
	Port               string
	Method             string
	Key                string
	ConnectionPoolSize uint
}

type Middleware func(net.Conn) (net.Conn, error)

func NewTcpServer(settings TcpServerSettings) *TcpServer {

	listener, err := net.Listen(settings.Method, fmt.Sprintf("%s:%s", settings.Host, settings.Port))

	if err != nil {
		return nil
	}

	ts := TcpServer{
		Connections:   sync.Map{},
		addChannel:    make(chan net.Conn, settings.ConnectionPoolSize),
		removeChannel: make(chan net.Conn, settings.ConnectionPoolSize),
	}
	ts.TCPListener = listener

	return &ts

}

func (ts *TcpServer) Close() error {
	return ts.TCPListener.Close()
}

func (ts *TcpServer) AssignTCPListener(listener net.Listener) {
	ts.TCPListener = listener
}

func (ts *TcpServer) AddChannel(conn net.Conn) {
	ts.addChannel <- conn
}
func (ts *TcpServer) RemoveChannel(conn net.Conn) {
	ts.removeChannel <- conn
}

func (ts *TcpServer) Run() {
	go func() {
		for {
			select {
			case conn := <-ts.addChannel:
				fmt.Printf("Added IP Address Channel: %v\n", conn.RemoteAddr().String())
				ts.Connections.Store(conn.RemoteAddr().String(), conn)
			case conn := <-ts.removeChannel:
				fmt.Printf("Removed  IP Address Channel: %v\n", conn.RemoteAddr().String())
				ts.Connections.Delete(conn.RemoteAddr().String())
			}
		}
	}()
	ts.ProcessConnections()
}

func (ts *TcpServer) ProcessConnections() {
	for {
		conn, err := ts.TCPListener.Accept()
		if err != nil {
			fmt.Print("TCP accept failed.\n")
			continue
		}

		_, exists := ts.Connections.Load(conn.RemoteAddr().String())
		if !exists {
			go ts.AddChannel(conn)
		}

		go func() {
			ts.distributePacketConn(conn)
		}()
	}
}

func (ts *TcpServer) denyPacketConn(conn net.Conn, reasoning string) {
	if conn == nil {
		fmt.Printf("Connection is nil: %s\n", reasoning)
		return
	}
	fmt.Printf("Closing connection: %v\n", conn.RemoteAddr())
	conn.Close()
	ts.RemoveChannel(conn)
}

func (ts *TcpServer) distributePacketConn(distConn net.Conn) {
	var size uint32
	err := binary.Read(distConn, binary.BigEndian, &size)
	if err != nil {
		// Handle error
	}

	var buffer []byte = make([]byte, size)
	_, err = distConn.Read(buffer)
	if err != nil {
		ts.denyPacketConn(distConn, fmt.Sprintf("Failed to parse message for distribution: %v", err))
		return
	}

	p := network.Packet{}
	p.FromByteSlice(buffer)

	fmt.Printf("%v | %v | %v | %v\n", p.Key, p.User, p.Length, string(p.Data))

	ts.Connections.Range(func(key, value any) bool {

		conn := value.(net.Conn)
		ip := key.(string)

		if ip == distConn.RemoteAddr().String() {
			return true
		}

		conn.Write(buffer)
		return true
	})

}
