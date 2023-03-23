package main

import (
	"fmt"

	"github.com/keybase/go-keychain"
)

func main() {
    item := keychain.NewItem()
    item.SetSecClass(keychain.SecClassGenericPassword)
    item.SetLabel("apple")
    item.SetReturnData(true)

    resules, err := keychain.QueryItem(item)
    if err != nil {
        panic(err)
    }

    fmt.Printf("res.len = %d\n", len(resules))
    for _, acc := range resules {
        fmt.Printf("Account: %+v\n", acc)
    }
    fmt.Println("query end....")
}
