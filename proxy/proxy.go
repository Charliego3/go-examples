package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/kataras/golog"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed proxy.pac
var pac string

func main() {
	golog.SetLevel("debug")
	golog.SetTimeFormat("2006-01-02 15:04:05.000")

	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		golog.Error(err)
		return
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		golog.Error("获取本机IP地址失败")
		return
	}

	firstAddr := ""

	for _, addr := range addrs {
		if addr, ok := addr.(net.Addr); ok {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil && strings.HasPrefix(ipnet.IP.String(), "192.168") {
					if firstAddr == "" {
						firstAddr = ipnet.IP.String()
					}
					golog.Infof("http://%s:8082/pac", ipnet.IP.String())
				}
			}
		}
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		golog.Error("获取Home失败", err)
		return
	}

	conf := filepath.Join(dir, ".vpn-proxy")
	file, err := os.OpenFile(conf, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		golog.Error("读取配置文件失败", err)
		return
	}

	var exclude []string
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			golog.Error("读取配置文件失败: ", err)
			return
		}

		if err == io.EOF {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		exclude = append(exclude, fmt.Sprintf("shExpMatch(host, \"%s\")", line))
	}

	file.Close()

	if len(exclude) == 0 {
		golog.Warnf("将需要排除的域名或IP配置在`%s`中, 按行分割")
	}

	go func() {
		var one sync.Once
		http.HandleFunc("/pac", func(w http.ResponseWriter, r *http.Request) {
			proxyAutoConfig := fmt.Sprintf(pac, strings.Join(exclude, " || \n\t"), firstAddr, firstAddr)
			w.Write([]byte(proxyAutoConfig))
			one.Do(func() {
				fmt.Printf("========== Proxy Auto-Config ==========\n%s\n=======================================\n", proxyAutoConfig)
			})
		})
		err := http.ListenAndServe(":8082", http.DefaultServeMux)
		if err != nil {
			golog.Error(err)
			return
		}
	}()

	for {
		client, err := l.Accept()
		if err != nil {
			golog.Error("Accept error:", err)
			continue
		}

		go handleClientRequest(client)
	}
}

func handleClientRequest(client net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			golog.Errorf("Handle client[%s] request error: %+v", client.RemoteAddr().String(), err)
		}
	}()
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		golog.Errorf("[%s]Read Error: %+v", client.RemoteAddr().String(), err)
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
		golog.Error("[%s]Scan host error: %+v", client.RemoteAddr().String(), err)
		return
	}
	golog.Infof("METHOD: %*q, ADDRESS: %q, HOST: %q", 9, method, address, host)
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
			golog.Error("NewRequest error:", err)
			return
		}
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			golog.Error("Http.Do(%s) error: %s", method, err)
			return
		}
		io.Copy(client, resp.Body)
		return
	}
	if !strings.Contains(host, ":") {
		hostPortURL, err := url.Parse(host)
		if err != nil {
			golog.Error(err)
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
		golog.Error(err)
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
