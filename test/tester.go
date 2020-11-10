package main

import "strings"

func main() {
	//CompleteNginx("172.16.100.123", "vip")

	//xmlUnmarshal()

	//getIncludeExists()

	//goQuery()

	s := "{\n  \"newCoinWhiteList\": \"110486,111022,110303,359784,359797,359783\",\n  \"rwOpenTime\": \"2017-11-29 12:00:00\",\n  \"transOpenTime\": \"2017-11-29 18:00:00\",\n  \"limitTransArea\": \"\"\n}dev@cloud-test-107-1:/home/appl/zbqc/conf$ "
	//s2 := s[:len("dev@cloud-test-107-1:/home/appl/zbqc/conf$")]
	s2 := strings.ReplaceAll(s, "dev@cloud-test-107-1:/home/appl/zbqc/conf$", "")
	println(s2)
}
