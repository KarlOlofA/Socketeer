package connectionPool

import (
	"fmt"
	"net"
	"sync"

	"socketeer.github.com/internal/types/network"
)

type ConnectionManager struct {
	TCPListener   net.Listener
	Connections   map[net.Conn]struct{}
	mu            sync.Mutex
	addChannel    chan net.Conn
	removeChannel chan net.Conn
	MiddleWare    []func(net.Conn) (net.Conn, error)
}

func NewConnectionManager(size int) *ConnectionManager {
	return &ConnectionManager{
		Connections:   make(map[net.Conn]struct{}, size),
		addChannel:    make(chan net.Conn, size),
		removeChannel: make(chan net.Conn, size),
	}

}

func (c *ConnectionManager) AssignTCPListener(listener net.Listener) {
	c.TCPListener = listener
}

func (c *ConnectionManager) AddChannel(conn net.Conn) {
	c.addChannel <- conn
}
func (c *ConnectionManager) RemoveChannel(conn net.Conn) {
	c.removeChannel <- conn
}

func (c *ConnectionManager) Run() {
	for {
		select {
		case conn := <-c.addChannel:
			c.mu.Lock()
			fmt.Printf("Added IP Address Channel: %v\n", conn.RemoteAddr().String())
			c.Connections[conn] = struct{}{}
			c.mu.Unlock()
		case conn := <-c.removeChannel:
			c.mu.Lock()
			fmt.Printf("Removed  IP Address Channel: %v\n", conn.RemoteAddr().String())
			delete(c.Connections, conn)
			c.mu.Unlock()
		}

	}
}

func (c *ConnectionManager) ProcessConnections() {
	for {
		conn, err := c.TCPListener.Accept()
		if err != nil {
			fmt.Print("TCP accept failed.\n")
			continue
		}

		if _, ok := c.Connections[conn]; !ok {
			go c.AddChannel(conn)
		}

		go func() {
			defer conn.Close()
			for {

				conn, err := c.ProcessMiddleware(conn)
				if err != nil {
					go c.denyPacketConn(conn, fmt.Sprintf("Middleware Error: %v", err))
					break
				} else if conn == nil {
					go c.denyPacketConn(conn, "Middleware failed to return a connection")
					break
				}

				c.distributePacketConn(conn)
			}
		}()

	}
}

func (p *ConnectionManager) AddMiddleware(mw func(net.Conn) (net.Conn, error)) *ConnectionManager {
	p.MiddleWare = append(p.MiddleWare, mw)
	return p
}

func (c *ConnectionManager) ProcessMiddleware(conn net.Conn) (net.Conn, error) {

	for _, mwf := range c.MiddleWare {
		_, err := mwf(conn)
		if err != nil {
			//conn.Write([]byte(fmt.Sprintf("%v", err)))
			return nil, err
		}
	}

	return conn, nil
}

func (c *ConnectionManager) denyPacketConn(conn net.Conn, reasoning string) {
	packet := []byte(reasoning)
	go c.RemoveChannel(conn)
	conn.Write(packet)
}

func (c *ConnectionManager) distributePacketConn(distConn net.Conn) {
	buffer := make([]byte, 1024)
	if _, err := distConn.Read(buffer); err != nil {
		return
	}

	p := network.Packet{}
	p.FromByteSlice(buffer)
	fmt.Printf("%v\n", p.Key)
	fmt.Printf("%v\n", p.Length)
	fmt.Printf("%v\n", p.User)
	fmt.Printf("%v\n", string(p.Data))

	for conn := range c.Connections {
		if conn.RemoteAddr().String() == distConn.RemoteAddr().String() {
			continue
		}

		conn.Write(buffer[:24+p.Length])
	}
}
