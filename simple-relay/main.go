// A messaging relay
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"syscall"
	"time"
)

func Relay(conn net.Conn, client net.Conn, retry, interval int, tx string) {
	b := bufio.NewReader(conn)

	for {
		line, err := b.ReadBytes('\n')
		if err != nil {
			log.Println("break ", err)
			break
		}

		log.Print("relay => ", string(line))
		_, err = client.Write(line)
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				log.Printf("EPIPE: %v", err)
				client.Close()
				client, err = TcpClient(tx, retry, interval)
				if err != nil {
					log.Println(err)
					break
				}
			} else {
				log.Println(err)
			}
		}
	}
}

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

func Receiver(l net.Listener, rx string) chan net.Conn {
	connChan := make(chan net.Conn)
	i := 0

	go func() {
		for {
			client, err := l.Accept()
			if client == nil {
				log.Println(err)
				continue
			}

			i++
			fmt.Printf("%d: %v -> %v\n", i, client.RemoteAddr(), client.LocalAddr())
			connChan <- client
		}
	}()

	return connChan
}

func main() {
	var (
		dst      = flag.String("dst", "127.0.0.1", "Destination IP or name")
		dport    = flag.Int("dport", 6000, "Destination port")
		src      = flag.String("src", "127.0.0.1", "Souce IP or name")
		sport    = flag.Int("sport", 5000, "Source port")
		retry    = flag.Int("retry", 3, "Number of times TCP client should attempt to connect")
		interval = flag.Int("interval", 2, "Number of seconds between retrys")
	)
	flag.Parse()

	rx := *src + ":" + strconv.Itoa(*sport)
	tx := *dst + ":" + strconv.Itoa(*dport)

	listener, err := net.Listen("tcp", rx)
	if listener == nil {
		log.Fatal(err)
	}
	defer listener.Close()
	connChan := Receiver(listener, rx)
	log.Println("reciever listening: ", rx)

	client, err := TcpClient(tx, *retry, *interval)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	log.Println("tcp-client connected: ", tx)

	for {
		select {
		case conn := <-connChan:
			go Relay(conn, client, *retry, *interval, tx)
		}
	}
}
