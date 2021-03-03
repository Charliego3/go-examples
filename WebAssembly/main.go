package main

import (
	"strconv"
	"syscall/js"
	"time"
)

var (
	//go:embed static/1.html
	h1 string
	//go:embed static/2.html
	h2 string
)

func fib(i int) int {
	return i * 2
}

func fibFunc(this js.Value, args []js.Value) interface{} {
	callback := args[len(args)-1]
	go func() {
		time.Sleep(200 * time.Millisecond)
		v := fib(args[0].Int())
		println("计算结果:", v)

		//// Set up a connection to the server.
		//conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
		//if err != nil {
		//	log.Fatalf("did not connect: %v", err)
		//}
		//defer conn.Close()
		//c := pb.NewGreeterClient(conn)
		//
		//// Contact the server and print out its response.
		//name := args[0].String()
		//if len(os.Args) > 1 {
		//	name = os.Args[1]
		//}
		////ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		////defer cancel()
		//
		//ctx := context.Background()
		//
		//r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
		//if err != nil {
		//	log.Fatalf("could not greet: %v", err)
		//}
		//log.Printf("Greeting: %s", r.GetMessage())
		//
		//callback.Invoke(strconv.Itoa(v) + " -> " + r.GetMessage())

		callback.Invoke(strconv.Itoa(v))
	}()

	js.Global().Get("ans").Set("innerHTML", "Waiting 1s...")
	return nil
}

func page(this js.Value, args []js.Value) interface{} {
	pageNum := args[0].Int()

	if pageNum == 1 {
		js.Global().Get("content").Set("innerHTML", h1)
	} else {
		js.Global().Get("content").Set("innerHTML", h2)
	}

	//	js.Global().Get("content").Set("innerHTML", fmt.Sprintf(`this is %s Page
	//
	//<input id="num" type="number" />
	//<!-- <button id="btn" onclick="ans.innerHTML=fibFunc(num.value * 1)">Click</button> -->
	//<button id="btn" onclick="fibFunc(num.value * 1, (v) => ans.innerHTML = v)">Click</button>
	//<p id="ans">1</p>`, pageNum))
	return nil
}

func main() {
	done := make(chan int, 0)
	// 注册全局函数
	js.Global().Set("fibFunc", js.FuncOf(fibFunc))
	js.Global().Set("page", js.FuncOf(page))
	<-done
}
