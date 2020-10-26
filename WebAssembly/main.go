package main

import (
	"fmt"
	"syscall/js"
	"time"
)

func fib(i int) int {
	if i == 0 || i == 1 {
		return 1
	}
	return fib(i-1) + fib(i-2)
}

func fibFunc(this js.Value, args []js.Value) interface{} {
	callback := args[len(args)-1]
	go func() {
		time.Sleep(3 * time.Second)
		v := fib(args[0].Int())
		fmt.Println(v)
		callback.Invoke(v)
	}()

	js.Global().Get("ans").Set("innerHTML", "Waiting 3s...")
	return nil
}

func main() {
	done := make(chan int, 0)
	// 注册全局函数
	js.Global().Set("fibFunc", js.FuncOf(fibFunc))
	<-done
}
