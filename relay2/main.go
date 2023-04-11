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
	"sync"
	"syscall"
	"time"
)

type TCPDest struct {
	connect  string
	retry    int
	interval int
}

func TcpClient(d TCPDest) (net.Conn, error) {
	var conn net.Conn
	var err error

	for i := 1; i <= d.retry; i++ {
		conn, err = net.Dial("tcp", d.connect)
		if err == nil {
			break
		}

		if errors.Is(err, syscall.ECONNREFUSED) {
			log.Printf("attempt %d: ECONNREFUSED: %v\n", i, err)
			if i == d.retry {
				err1 := errors.New("connection attempts exhausted")
				return nil, err1
			}
		} else {
			return nil, err
		}

		time.Sleep(time.Duration(d.interval) * time.Second)
	}

	return conn, err
}

func Receiver(rx string) (net.Listener, [10]chan net.Conn) {
	var connChan [10]chan net.Conn
	for j := range connChan {
		connChan[j] = make(chan net.Conn, 10)
	}

	l, err := net.Listen("tcp", rx)
	if l == nil {
		log.Fatal(err)
	}

	i := 0

	go func() {
		for {
			client, err := l.Accept()
			if client == nil {
				log.Println(err)
				continue
			}

			fmt.Printf("%d: %v -> %v\n", i, client.RemoteAddr(), client.LocalAddr())
			connChan[10%i] <- client
			i++
		}
	}()

	return l, connChan
}

func Relay(connChan chan net.Conn, wg *sync.WaitGroup, d TCPDest) {
	client, err := TcpClient(d)
	if err != nil {
		log.Println("relay worker fail: ", err)
		wg.Done()
		return
	}

	for {
		select {
		case conn := <-connChan:
			b := bufio.NewReader(conn)
			line, err := b.ReadBytes('\n')
			if err != nil {
				log.Println("relay worker read error. continue: ", err)
				conn.Close()
				continue
			}

			log.Print("relay => ", string(line))
			_, err = client.Write(line)
			if err != nil {
				if errors.Is(err, syscall.EPIPE) {
					log.Printf("EPIPE: %v", err)
					client.Close()
					client, err = TcpClient(d)
					if err != nil {
						log.Println("relay worker epipe: unrecoveable: ", err)
						wg.Done()
						break
					}
				} else {
					log.Println("relay worker unknown: ", err)
					wg.Done()
					break
				}
			}
		}
	}
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
	listener, connChan := Receiver(rx)
	defer listener.Close()
	log.Println("reciever listening: ", rx)

	tx := *dst + ":" + strconv.Itoa(*dport)
	var d TCPDest
	d.connect = tx
	d.retry = *retry
	d.interval = *interval
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go Relay(connChan[i], &wg, d)
	}

	wg.Wait()
}
