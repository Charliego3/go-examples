package main

import (
	"fmt"
	"syscall/js"
	"time"
)

func fib(i int) int {
	return i * 2
}

func fibFunc(this js.Value, args []js.Value) interface{} {
	callback := args[len(args)-1]
	go func() {
		time.Sleep(300 * time.Millisecond)
		v := fib(args[0].Int())
		println("计算结果:", v)
		callback.Invoke(v)
	}()

	js.Global().Get("ans").Set("innerHTML", "Waiting 1s...")
	return nil
}

func page(this js.Value, args []js.Value) interface{} {
	pageNum := args[0].String()

	js.Global().Get("content").Set("innerHTML", fmt.Sprintf(`this is %s Page

<input id="num" type="number" />
<!-- <button id="btn" onclick="ans.innerHTML=fibFunc(num.value * 1)">Click</button> -->
<button id="btn" onclick="fibFunc(num.value * 1, (v) => ans.innerHTML = v)">Click</button>
<p id="ans">1</p>`, pageNum))
	return nil
}

func main() {
	done := make(chan int, 0)
	// 注册全局函数
	js.Global().Set("fibFunc", js.FuncOf(fibFunc))
	js.Global().Set("page", js.FuncOf(page))
	<-done
}
