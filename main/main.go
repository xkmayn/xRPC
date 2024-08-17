package main

import (
	"context"
	xRPC "github.com/xkmayn/xrpc"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type XK int

type Args struct{ Num1, Num2 int }

func (xk XK) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func StartServer(addr chan string) {
	var xk XK
	if err := xRPC.Register(&xk); err != nil {
		log.Fatal("register error", err)
	}
	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on :", l.Addr())
	xRPC.HandleHTTP()
	addr <- l.Addr().String()
	_ = http.Serve(l, nil)
}

func call(addr chan string) {
	client, _ := xRPC.XDial("http@" + <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)

	wg := new(sync.WaitGroup)

	for i := 1; i <= 5; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			args := &Args{i, i * i}
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			var reply int
			if err := client.Call(ctx, "XK.Sum", args, &reply); err != nil {
				log.Fatal("call XK.Sum err:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go call(addr)
	StartServer(addr)
}
