package main

import (
	"encoding/xml"
	"github.com/fatih/color"
	"io/ioutil"
)

type Server struct {
	Service struct {
		Connector []struct {
			Port string `xml:"port,attr"`
		} `xml:"Connector"`
	} `xml:"Service"`
}

func xmlUnmarshal() {
	content, err := ioutil.ReadFile("/Users/nzlong/dev/apache-tomcat-8.5.39/conf/server.xml")
	if err != nil {
		color.Red(err.Error())
		return
	}

	server := &Server{}
	_ = xml.Unmarshal(content, server)

	println(server)
}
