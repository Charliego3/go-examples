package main

import (
	"github.com/whimthen/temp/meituan/apis"
	"github.com/whimthen/temp/meituan/repos"
	"github.com/whimthen/temp/times"
	"log"
)

func main() {
	milliseconds := apis.Milliseconds()
	log.Println("服务器时间:", milliseconds, times.Parse(milliseconds))

	repos.FetchUser()

	// resp, err := apis.LoginV5("18929387993", "yy12347890", "6519763", true)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	//
	// log.Printf("%+v\n", resp)
}
