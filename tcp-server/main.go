package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func PrintLine(conn net.Conn) {
	b := bufio.NewReader(conn)
	for {
		line, err := b.ReadBytes('\n') 
		if err != nil {
			log.Println("break ", err)
 			break
 		}

		fmt.Fprintf(os.Stdout, "%s", string(line))
	}
}

func Server(l net.Listener, rx string) chan net.Conn {
	channel := make(chan net.Conn)
	i := 0

	go func() {
		for {
			conn, err := l.Accept()
			if conn == nil {
				log.Println(err)
				continue
			}
			i++
			fmt.Printf("%d: %v <-> %v\n", i, conn.LocalAddr(), conn.RemoteAddr())
			channel <- conn
		}
	}()

	return channel
}

func main() {
	var (
		bind = flag.String("bind-ip", "127.0.0.1", "Hostname / IP address to bind to")
		port = flag.Int("port", 5000, "Port to bind to")
	)
	flag.Parse()

	rx := *bind + ":" + strconv.Itoa(*port)

	l, err := net.Listen("tcp", rx)
	if l == nil {
		log.Fatal(err)
	}
	defer l.Close()
	conn := Server(l, rx)
	log.Println("Listening on " + rx)

	for {
		go PrintLine(<-conn)
	}
}

