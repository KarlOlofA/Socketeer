package connectionPool

import (
	"fmt"
	"net"
	"sync"

	"socketeer.github.com/internal/types/network"
)

type TcpServer struct {
	TCPListener   net.Listener
	Connections   map[net.Conn]struct{}
	mu            sync.Mutex
	addChannel    chan net.Conn
	removeChannel chan net.Conn
	MiddleWare    []func(net.Conn) (net.Conn, error)
}

type TcpServerSettings struct {
	Host               string
	Port               string
	Method             string
	Key                string
	ConnectionPoolSize uint
}

func NewTcpServer(settings TcpServerSettings) *TcpServer {

	listener, err := net.Listen(settings.Method, fmt.Sprintf("%s:%s", settings.Host, settings.Port))

	if err != nil {
		return nil
	}
	defer listener.Close()

	ts := TcpServer{
		Connections:   make(map[net.Conn]struct{}, settings.ConnectionPoolSize),
		addChannel:    make(chan net.Conn, settings.ConnectionPoolSize),
		removeChannel: make(chan net.Conn, settings.ConnectionPoolSize),
	}
	ts.TCPListener = listener

	return &ts

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
			ts.mu.Lock()
			fmt.Printf("Added IP Address Channel: %v\n", conn.RemoteAddr().String())
			ts.Connections[conn] = struct{}{}
			ts.mu.Unlock()
		case conn := <-ts.removeChannel:
			ts.mu.Lock()
			fmt.Printf("Removed  IP Address Channel: %v\n", conn.RemoteAddr().String())
			delete(ts.Connections, conn)
			ts.mu.Unlock()
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

		if _, ok := ts.Connections[conn]; !ok {
			go ts.AddChannel(conn)
		}

		go func() {
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

func (ts *TcpServer) AddMiddleware(mw func(net.Conn) (net.Conn, error)) *TcpServer {
	ts.MiddleWare = append(ts.MiddleWare, mw)
	return ts
}

func (ts *TcpServer) ProcessMiddleware(conn net.Conn) (net.Conn, error) {

	for _, mwf := range ts.MiddleWare {
		_, err := mwf(conn)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("%v", err)))
			return nil, err
		}
	}

	return conn, nil
}

func (ts *TcpServer) denyPacketConn(conn net.Conn, reasoning string) {
	packet := []byte(reasoning)
	go ts.RemoveChannel(conn)
	conn.Write(packet)
}

func (ts *TcpServer) distributePacketConn(distConn net.Conn) {
	buffer := make([]byte, 1024)
	if _, err := distConn.Read(buffer); err != nil {
		return
	}

	p := network.Packet{}
	p.FromByteSlice(buffer)

	for conn := range ts.Connections {
		if conn.RemoteAddr().String() == distConn.RemoteAddr().String() {
			continue
		}

		conn.Write(buffer[:24+p.Length])
	}
}
