/*
 *  relay.go
 *
 *  A simple relay server
 */

package main

import (
  "bufio"
  "flag"
  "fmt"
  "net"
)


var dstip, dstport, port string

func init() {
  flag.StringVar(&dstport, "dstip",   "127.0.0.1", "Relay destination IP")
  flag.StringVar(&dstport, "dstport", "3456",      "Relay destination port")
  flag.StringVar(&port,    "port",    "2345",      "Server bind port")
}


func relay(message []byte) string {
  status := "failed"
  c, err := net.Dial("tcp", dstip + ":" + dstport)
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

    status := relay(line);
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
  flag.Parse()

  server, err := net.Listen("tcp", ":" + port)
  if server == nil {
    panic(err)
  }

  client := Server(server)

  for {
    go handler(<-client)
  }
}


