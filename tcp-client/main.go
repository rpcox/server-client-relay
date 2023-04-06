// TCP client
package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"strconv"
	"syscall"
	"time"
)

func TcpClient(tx string, retry, interval int) (net.Conn, error) {

	var conn net.Conn
	var err error

	for i := 1; i <= retry; i++ {
		conn, err = net.Dial("tcp", tx)
		if err == nil {
			break
		}

		if errors.Is(err, syscall.ECONNREFUSED) {
			log.Printf("attempt %d: ECONNREFUSED: %v\n", i, err)
			if i == retry {
				err1 := errors.New("connection attempts exhausted")
				return nil, err1
			}
		} else {
			return nil, err
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}

	return conn, err
}

func main() {
	var (
		dst      = flag.String("dst", "127.0.0.1", "Destination hostname or IP address")
		dport    = flag.Int("dport", 6000, "Destination port")
		msg      = flag.String("msg", "TCP-Client TEST message", "Repeating message to send")
		retry    = flag.Int("retry", 3, "Number of times TCP client should attempt to connect")
		interval = flag.Int("interval", 2, "Number of seconds between retrys")
	)

	flag.Parse()

	tx := *dst + ":" + strconv.Itoa(*dport)
	client, err := TcpClient(tx, *retry, *interval)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Sending to ", tx)

	message := []byte(*msg)
	message = append(message, '\n')

	for {
		_, err := client.Write(message)
		if err != nil {
			// server closed the socket
			if errors.Is(err, syscall.EPIPE) {
				log.Printf("EPIPE: %v", err)
				client.Close()
				client, err = TcpClient(tx, *retry, *interval)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
