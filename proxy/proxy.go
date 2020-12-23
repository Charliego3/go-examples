package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Panic(err)
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}

		go handleClientRequest(client)
	}
}

func handleClientRequest(client net.Conn) {
	go func() {
		if err := recover(); err != nil {
			log.Println("Handle client request error:", err)
		}
	}()
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	var method, host, address string
	_, err = fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	if err != nil {
		log.Println("Scan host error:", err)
		return
	}
	log.Printf("METHOD: %s, HOST: %s, ADDRESS: %sl\n", method, host, address)
	if strings.Index(host, "http://") == 0 || strings.Index(host, "https://") == 0 {
		request, err := http.NewRequest(method, host, nil)
		// resp, err := http.Get(host)
		if err != nil {
			log.Println("NewRequest error:", err)
			return
		}
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Printf("Http.Do(%s) error: %s", method, err)
			return
		}
		io.Copy(client, resp.Body)
		return
	}
	if !strings.Contains(host, ":") {
		hostPortURL, err := url.Parse(host)
		if err != nil {
			log.Println(err)
			return
		}

		if hostPortURL.Opaque == "443" { //https访问
			address = hostPortURL.Scheme + ":443"
		} else { //http访问
			if strings.Index(hostPortURL.Host, ":") == -1 { //host不带端口， 默认80
				address = hostPortURL.Host + ":80"
			} else {
				address = hostPortURL.Host
			}
		}
	} else {
		address = host
	}

	log.Printf("Client Conn: %+v, Address: %s", client, address)

	//获得了请求的host和port，就开始拨号吧
	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	if method == "CONNECT" {
		_, _ = fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		_, _ = server.Write(b[:n])
	}
	//进行转发
	go io.Copy(server, client)
	_, _ = io.Copy(client, server)
}
