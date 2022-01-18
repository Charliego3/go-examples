package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/kataras/golog"
	"io"
	"net"
)

func main() {
	port := flag.String("port", ":8888", "--port")
	flag.Parse()

	golog.Errorf("Port: %s", *port)
	listen, err := net.Listen("tcp", *port)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			golog.Error("Accept error:", err)
			continue
		}

		go handlerConn(conn)
	}
}

func handlerConn(client net.Conn) {
	defer func(client net.Conn) {
		err := client.Close()
		if err != nil {
			golog.Error("connection close error", err)
		}

		if err := recover(); err != nil {
			golog.Error(err)
		}
	}(client)

	var b [1024]byte
	n, err := client.Read(b[:])
	_ = n
	if err != nil {
		golog.Errorf("[%s]Read Error: %+v", client.RemoteAddr().String(), err)
		return
	}
	//log.Printf("Bytes: %s\n", b[:])
	var method, host string
	indexByte := bytes.IndexByte(b[:], '\n')
	if indexByte == -1 {
		return
	}
	_, err = fmt.Sscanf(string(b[:indexByte]), "%s%s", &method, &host)
	if err != nil {
		golog.Error("[%s]Scan host error: %+v", client.RemoteAddr().String(), err)
		return
	}

	golog.Errorf("%+v, Addr:%+v, method:%s, host:%s", client, client.LocalAddr(), method, host)

	//获得了请求的host和port，就开始拨号吧
	server, err := net.Dial("tcp", "tradingbot.100-130.net:48620")
	if err != nil {
		golog.Error(err)
		return
	}
	defer server.Close()
	if method == "CONNECT" {
		_, _ = fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		_, _ = server.Write(b[:n])
	}
	//进行转发
	go io.Copy(server, client)
	_, _ = io.Copy(client, server)
}
