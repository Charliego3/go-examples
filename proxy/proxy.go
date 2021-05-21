package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

//go:embed proxy.pac
var pac []byte

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Panic(err)
	}

	go func() {
		var one sync.Once
		http.HandleFunc("/pac", func(w http.ResponseWriter, r *http.Request) {
			proxyAutoConfig := bytes.ReplaceAll(pac, []byte("loopbackAddress"), []byte("192.168.1.20"))
			w.Write(proxyAutoConfig)
			one.Do(func() {
				fmt.Printf("\n========== Proxy Auto-Config ==========\n%s\n=======================================\n", string(proxyAutoConfig))
			})
		})
		err := http.ListenAndServe(":8082", http.DefaultServeMux)
		if err != nil {
			log.Println(err)
			return
		}
	}()

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
	defer func() {
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
	//log.Printf("Bytes: %s\n", b[:])
	var method, host, address string
	indexByte := bytes.IndexByte(b[:], '\n')
	if indexByte == -1 {
		return
	}
	_, err = fmt.Sscanf(string(b[:indexByte]), "%s%s", &method, &host)
	if err != nil {
		log.Println("Scan host error:", err)
		return
	}
	log.Printf("METHOD: %q, HOST: %q, ADDRESS: %q\n", method, host, address)
	if strings.Index(host, "http://") == 0 || strings.Index(host, "https://") == 0 {

		// 解析POST请求参数
		buffer := bytes.NewBuffer(b[:])
		index := 0
		var body bytes.Buffer
		for {
			line, err := buffer.ReadBytes('\n')
			if err != nil && err == io.EOF {
				break
			}

			index += len(line)
			if string(line) == "\r\n" {
				//log.Printf("Body: %s", b[index+1:])
				body.Write(b[index+1:])
				break
			}
		}

		request, err := http.NewRequest(method, host, &body)
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

	//log.Printf("Client Conn: %+v, Address: %s", client, address)

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
