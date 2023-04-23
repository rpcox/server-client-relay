## echo

Simple echo server

#### usage

    > echo -h

### example

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


