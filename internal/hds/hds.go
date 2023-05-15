package hds

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import (
	"bytes"
	"encoding/gob"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

var (

	// mu is
	// mutex for thread
	// safe add / rmv ip
	// from baseNodesIps
	mu = &sync.Mutex{}

	// baseNodesIps is
	// map of connected base nodes
	baseNodesIps = map[string]bool{}

	// wsConnUpgraded is upgrader
	// from HTTP(s) to web sockets connection
	wsConnUpgrader = websocket.Upgrader{}
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// if base node
// disconnects from root node - add ip
func addToBaseNodesIps(ip string) {

	mu.Lock()
	defer mu.Unlock()

	baseNodesIps[ip] = true

}

// ------------------------------------------------------------------------ //

// if base node
// disconnects from root node - rmv ip
func rmFromBaseNodeIps(ip string) {

	mu.Lock()
	defer mu.Unlock()

	delete(baseNodesIps, ip)

}

// ------------------------------------------------------------------------ //

// parseIp is function
// for separating remote ip address from port
func parseIp(remoteAddr string) string {
	return strings.Split(remoteAddr, ":")[0]
}

// ------------------------------------------------------------------------ //

// excludeBaseNodeIp is function
// that returns array of connected base
// node ips with the exception of connected node ip
func excludeBaseNodeIp(ip string) []string {

	var (

		// tempBaseNodesIps is connected
		// connected base nodes except for
		// the ip address of the node that requests
		tempBaseNodesIps []string
	)

	mu.Lock()
	defer mu.Unlock()

	for k, _ := range baseNodesIps {
		if k != ip {
			tempBaseNodesIps = append(tempBaseNodesIps, k)
		}
	}

	return tempBaseNodesIps

}

// ------------------------------------------------------------------------ //

// nodeBaseReader is function
// if base node disconnects - we remove it from
// baseNodesIps and close web sockets connection
func nodeBaseReader(wsConn *websocket.Conn, c chan<- bool, ip string) {

	for {

		if _, _, err := wsConn.ReadMessage(); err != nil {
			rmFromBaseNodeIps(ip)
			wsConn.Close()
			c <- true
			return
		}

	}

}

// ------------------------------------------------------------------------ //

// encodeComd is function for
// converting command in string view to
// command in bytes array for sending to base node
func encodeComd(comd string) []byte {
	var bytes [20]byte

	for i, c := range comd {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

// ------------------------------------------------------------------------ //

// encodeData is function for
// converting data passed as interface
// in bytes array for sending to base node
func encodeData(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// ------------------------------------------------------------------------ //

// nodeBaseWriter is function
// sends to connected base nodes array of available base nodes
func nodeBaseWriter(wsConn *websocket.Conn, c <-chan bool, ip string) {

	for {

		select {
		case <-c:
			{
				return
			}
		default:
			{
				wsConn.WriteMessage(
					websocket.BinaryMessage,
					append(encodeComd("nodeIps"), encodeData(excludeBaseNodeIp(ip))...),
				)
				time.Sleep(time.Second * 5)
			}
		}

	}

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// BaseNodeGetIp is handler
// for "/sh-ip" web server route
//
// For defining blockchain base node in internet network,
// it have to know own public IP address for adding base nodes
func BaseNodeGetIp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(parseIp(r.RemoteAddr)))
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// BaseNodeTrack is handler
// for "/track" web server route
//
// After starting the base node it
// connects to root node in blockchain
// This handler is the way for defining something like DNS for blockchain
func BaseNodeTrack(w http.ResponseWriter, r *http.Request) {

	var (

		// ip is
		// ip address of base node
		ip = parseIp(r.RemoteAddr)

		// ch is
		// channel for synchronization ws reader and writes
		ch = make(chan bool)
	)

	// upgrade HTTP(s)
	// connection to web sockets connection
	conn, connErr := wsConnUpgrader.Upgrade(w, r, nil)
	if connErr != nil {
		return
	}

	addToBaseNodesIps(ip)

	go nodeBaseWriter(conn, ch, ip)
	go nodeBaseReader(conn, ch, ip)

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
