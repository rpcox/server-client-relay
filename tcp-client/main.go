package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
)

func TcpClient(dst string) net.Conn {
	/*var dialer net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := dialer.DialContext(ctx, "tcp", dst)*/
	conn, err := net.Dial("tcp", dst)

	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) {
			log.Fatal("ECONNREFUSED: ", err)
		} else {
			log.Fatal(err)
		}
	}

	return conn
}

func main() {
	var (
		serverDst = flag.String("dst-server", "127.0.0.1", "Destination hostname / IP address")
		portDst   = flag.Int("dst-port", 5000, "Destination port")
	)

	dst := *serverDst + ":" + strconv.Itoa(*portDst)
	client := TcpClient(dst)
	log.Println("Sending to ", dst)

	for {
		_, err := client.Write([]byte("TEST\n"))
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				log.Printf("EPIPE: %v", err)
			} else {
				log.Println(err)
			}
			os.Exit(1)
		}
		time.Sleep(1 * time.Second)
	}
}

