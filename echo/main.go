// A simple echo server
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

func echo(client net.Conn) {
	b := bufio.NewReader(client)

	for {
		line, err := b.ReadBytes('\n') // read it
		if err != nil {
			break
		}

		client.Write(line) // send it
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

var port, help string

func init() {
	flag.StringVar(&port, "port", "2345", "Server bind port")
}

func main() {
	flag.Parse()

	server, err := net.Listen("tcp", ":"+port)
	if server == nil {
		panic(err)
	}

	client := Server(server)

	for {
		go echo(<-client)
	}
}
