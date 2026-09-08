package connectionPool

import (
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
	for {
		select {
		case conn := <-ts.addChannel:
			fmt.Printf("Added IP Address Channel: %v\n", conn.RemoteAddr().String())
			ts.Connections.Store(conn, struct{}{})
		case conn := <-ts.removeChannel:
			fmt.Printf("Removed  IP Address Channel: %v\n", conn.RemoteAddr().String())
			ts.Connections.Delete(conn)
		}

	}
}

func (ts *TcpServer) ProcessConnections() {
	for {
		conn, err := ts.TCPListener.Accept()
		if err != nil {
			fmt.Print("TCP accept failed.\n")
			continue
		}

		_, exists := ts.Connections.Load(conn)
		if !exists {
			go ts.AddChannel(conn)
		}

		func() {
			defer conn.Close()
			for {
				conn, err := ts.ProcessMiddleware(conn)
				if err != nil {
					go ts.denyPacketConn(conn, fmt.Sprintf("Middleware Error: %v", err))
					break
				} else if conn == nil {
					go ts.denyPacketConn(conn, "Middleware failed to return a connection")
					break
				}
				ts.distributePacketConn(conn)
			}
		}()
	}
}

func (ts *TcpServer) AddMiddleware(mw Middleware) *TcpServer {
	ts.Middleware = append(ts.Middleware, mw)
	return ts
}

func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(conn net.Conn) (net.Conn, error) {
		for _, mw := range middlewares {
			var err error
			conn, err = mw(conn)
			if err != nil {
				return nil, err
			}
		}
		return conn, nil
	}
}

func (ts *TcpServer) ProcessMiddleware(conn net.Conn) (net.Conn, error) {
	for _, mw := range ts.Middleware {
		var err error
		conn, err = mw(conn)
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}

func (ts *TcpServer) denyPacketConn(conn net.Conn, reasoning string) {
	if conn == nil {
		fmt.Printf("Connection is nil: %s\n", reasoning)
		return
	}
	fmt.Printf("Closing connection: %v\n", conn.RemoteAddr())
	packet := []byte(reasoning)
	if _, err := conn.Write(packet); err != nil {
		fmt.Printf("%v\n", err)
	}
	conn.Close()
	go ts.RemoveChannel(conn)
}

func (ts *TcpServer) distributePacketConn(distConn net.Conn) {
	fmt.Printf("Wa?\n")
	var buffer []byte = make([]byte, 1024)
	_, err := distConn.Read(buffer)
	if err != nil {
		fmt.Printf("Failed to parse message for distribution: %v\n", err)
		return
	}

	p := network.Packet{}
	p.FromByteSlice(buffer)

	fmt.Printf("%v | %v | %v | %v\n", p.Key, p.User, p.Length, string(p.Data))

	ts.Connections.Range(func(key, value any) bool {

		conn := key.(net.Conn)

		if conn.RemoteAddr().String() == distConn.RemoteAddr().String() {
			return true
		}

		conn.Write(buffer[:24+p.Length])
		return true
	})

}
