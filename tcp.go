package tcp

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/tcp", new(TCP))
}

// TCP is the k6 tcp extension.
type TCP struct{}

// Client is the TCP Client wrapper.
type Client struct {

	// conn
	conn net.Conn

	// address and port, for example, 127.0.0.1:8000.
	connStr string

	// The last error from this struct.
	lastErr error

	// The function pointer used to receive the message needs to be passed in JS.
	onRevMsg CallBack
}

// XClient represents the Client constructor (i.e. `new tcp.Client()`) and
// returns a new TCP client object.
func (r *TCP) XClient(ctxPtr *context.Context) interface{} {
	rt := common.GetRuntime(*ctxPtr)
	return common.Bind(rt, &Client{}, ctxPtr)
}

// Create new tcp.
// A new TCP link must be used, otherwise it will cause the same conn to communicate.
//func (tcp *Client) Create() *Client {
//	return new(Client)
//}

// Send the received message to JS. Com through this function.
type CallBack func(data []byte)

// To init all things.
func (client *Client) Connect(addr string, onRevMsg CallBack) error {
	client.connStr = addr
	client.onRevMsg = onRevMsg
	client.conn, client.lastErr = net.Dial("tcp", client.connStr)
	if client.lastErr != nil {
		client.conn = nil
		return client.lastErr
	} else {
		go client.readConn()
	}

	return nil
}

// Send msg by this function.
func (client *Client) WriteStr(data string) error {

	if client.conn == nil {
		return errors.New("call Write function, but conn is nil.")
	}

	_, err := client.conn.Write([]byte(data))
	if err != nil {
		return err
	}

	return nil
}

// Send msg by this function.
func (client *Client) WriteStrLn(data string) error {
	return client.WriteStr(fmt.Sprintln(data))
}

// Get msg by this function.
func (client *Client) readConn() {
	defer func() {
		log.Println("readConn Err")
	}()

	for {
		scanner := bufio.NewScanner(client.conn)

		if scanner != nil {
			for {
				ok := scanner.Scan()

				if !ok {
					fmt.Println("Reached EOF on server connection.")
					break
				} else {
					data := scanner.Bytes()

					if client.onRevMsg != nil {
						client.onRevMsg(data)
					}
				}
			}
		}
	}
}
