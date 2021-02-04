## echo-relay

Simple echo and relay servers

#### Make the executables

    >  go build echo.go && go build relay.go

#### Usage

    > echo -h
    > relay -h

### echo

                     SERVER                 |            CLIENT           |

     > go run echo.go                       |
     1: 127.0.0.1:2345 <-> 127.0.0.1:39606  | > telnet localhost 2345
     2: 127.0.0.1:2345 <-> 127.0.0.1:39608  | Trying 127.0.0.1...
                                            | Connected to localhost.
                                            | Escape character is '^]'.
                                            | echo
                                            | echo
                                            | ^]
                                            | telnet> Connection closed.
                                            | >
                                            | > telnet localhost 2345
                                            | Trying 127.0.0.1...
                                            | Connected to localhost.
                                            | Escape character is '^]'.
                                            | echo
                                            | echo
                                            | ^]
                                            | telnet> Connection closed.
                                            | >


### relay

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

