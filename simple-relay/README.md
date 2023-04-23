## relay

Simple relay server

#### usage

    > relay -h

### example

                     SERVER                 |            CLIENT           |        DESTINATION

     > go run echo.go                       |                             | > nc -lk 3456
     1: 127.0.0.1:2345 <-> 127.0.0.1:40422  | > telnet localhost 2345     |
                                            | Trying 127.0.0.1...         |
                                            | Connected to localhost.     |
                                            | Escape character is '^]'.   |
                                            | relay my message            | relay my messag
                                            | relay success               |
                                            | do it again                 | do it again
                                            | relay success               | ^C
                                            | far side is down            |
                                            | relay failed                |
                                            | ^]                          |
                                            | telnet> Connection closed.  |
                                            | >                           |

