package main

import (
	xRPC "github.com/xkmayn/xrpc"
	"log"
	"net"
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
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on :", l.Addr())
	addr <- l.Addr().String()
	xRPC.Accept(l)
}

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go StartServer(addr)

	client, _ := xRPC.Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)

	wg := new(sync.WaitGroup)

	for i := 1; i <= 5; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			args := &Args{i, i * i}
			var reply int
			if err := client.Call("XK.Sum", args, &reply); err != nil {
				log.Fatal("call XK.Sum err:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}
