package main

import (
	"crypto/rc4"
	"flag"
	"fmt"
	"net"
)

func handle_read(rconn net.Conn, wconn net.Conn, password string) {
	defer rconn.Close()
	defer wconn.Close()

	var rbuf = make([]byte, 1500)
	var wbuf = make([]byte, 1500)
	var key []byte = []byte(password)
	var cipher, _ = rc4.NewCipher(key)
	for {
		if n, err := rconn.Read(rbuf); err == nil {
			cipher.XORKeyStream(wbuf[:n], rbuf[:n])
			wconn.Write(wbuf[:n])
		} else {
			return
		}
	}
}

func handle(conn net.Conn, password string, forward string) {
	var forconn, err = net.Dial("tcp", forward)
	if err != nil {
		return
	}
	go handle_read(conn, forconn, password)
	go handle_read(forconn, conn, password)
}

func main() {
	var port = flag.Uint("l", 80, "port")
	var forward = flag.String("f", "127.0.0.1:8080", "forward")
	var password = flag.String("p", "password", "password")
	flag.Parse()
	var lis, _ = net.Listen("tcp", fmt.Sprintf(":%v", *port))
	defer lis.Close()
	for {
		var conn, _ = lis.Accept()
		go handle(conn, *password, *forward)
	}
}
