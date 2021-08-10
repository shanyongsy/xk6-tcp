package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"net"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/tcp", new(TCP))
}

type CallBack func(data []byte)

type TCP struct {

	// conn
	conn net.Conn

	// address and port, for example, 127.0.0.1:8000.
	connStr string

	// The last error from this struct.
	lastErr error

	// The function pointer used to receive the message needs to be passed in JS.
	onRevMsg CallBack
}

// Create new tcp.
// A new TCP link must be used, otherwise it will cause the same conn to communicate.
func (tcp *TCP) Create() *TCP {
	return new(TCP)
}

// To init all things.
func (tcp *TCP) Connect(addr string, onRevMsg CallBack) error {
	tcp.connStr = addr
	tcp.onRevMsg = onRevMsg
	tcp.conn, tcp.lastErr = net.Dial("tcp", tcp.connStr)
	if tcp.lastErr != nil {
		tcp.conn = nil
		return tcp.lastErr
	} else {
		go tcp.readConn()
	}

	return nil
}

// Send msg by this function.
func (tcp *TCP) Write(data string) error {

	if tcp.conn == nil {
		return errors.New("call Write function, but conn is nil.")
	}

	_, err := tcp.conn.Write([]byte(data))
	if err != nil {
		return err
	}

	return nil
}

// Send msg by this function.
func (tcp *TCP) WriteLn(data string) error {
	return tcp.Write(fmt.Sprintln(data))
}

// Get msg by this function.
func (tcp *TCP) readConn() {
	for {
		scanner := bufio.NewScanner(tcp.conn)

		if scanner != nil {
			for {
				ok := scanner.Scan()

				if !ok {
					fmt.Println("Reached EOF on server connection.")
					break
				} else {
					text := scanner.Bytes()
					fmt.Println(text)

					if tcp.onRevMsg != nil {
						tcp.onRevMsg(text)
					}
				}
			}
		}
	}
}
