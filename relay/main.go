// A messaging relay
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strconv"
)

func relay(message []byte) string {
	status := "failed"
	c, err := net.Dial("tcp", dstip+":"+dstport)
	if err != nil {
		return status
	}

	c.Write(message)
	if err != nil {
		return status
	}

	status = "success"
	c.Close()
	return status
}

func handler(client net.Conn) {
	b := bufio.NewReader(client)

	for {
		line, err := b.ReadBytes('\n')
		if err != nil {
			break
		}

		status := relay(line)
		b := []byte("relay " + status + "\n")
		client.Write(b)
	}
}

func Server(listener net.Listener) chan net.Conn {
	channel := make(chan net.Conn)
	i := 0

	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				fmt.Printf("Accept() fail : ", err)
				continue
			}

			i++
			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			channel <- client
		}
	}()

	return channel
}

func main() {
	var (
		serverDst = flag.String("dst", "127.0.0.1", "Destination IP or name")
		portDst   = flag.Int("dport", 5000, "Destination port")
		serverSrc = flag.String("src", 8000, "Souce IP or name")
		portSrc   = flag.Int("sport", 5000, "Source port")
	)
	flag.Parse()

	rx := *serverSrc + ":" + strconv.Itoa(*portSrc)
	tx := *serverDst + ":" + strconv.Itoa(*portDst)

	dst := DestinationConnect(tx)
	defer dst.Close()

	listener := SourceConnect(rx)
	defer listener.Close()

	server, err := net.Listen("tcp", ":"+port)
	if server == nil {
		panic(err)
	}

	client := Server(server)

	for {
		go handler(<-client)
	}
}
