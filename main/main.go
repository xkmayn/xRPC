package main

import (
	"context"
	xRPC "github.com/xkmayn/xrpc"
	"github.com/xkmayn/xrpc/xclient"
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

func (xk XK) Sleep(args Args, reply *int) error {
	time.Sleep(time.Second * time.Duration(args.Num1))
	*reply = args.Num1 + args.Num2
	return nil
}

func StartServer(addr chan string) {
	var xk XK
	server := xRPC.NewServer()
	if err := server.Register(&xk); err != nil {
		log.Fatal("register error", err)
	}
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}

	log.Println("start rpc server on :", l.Addr())
	// server.HandleHTTP()
	addr <- l.Addr().String()
	server.Accept(l)
	// _ = http.Serve(l, nil)
}

func xk(xc *xclient.XClient, ctx context.Context, typ, serviceMethod string, args *Args) {
	var reply int
	var err error
	switch typ {
	case "call":
		err = xc.Call(ctx, serviceMethod, args, &reply)
	case "broadcast":
		err = xc.Broadcast(ctx, serviceMethod, args, &reply)
	}
	if err != nil {
		log.Printf("%s %s error: %v", typ, serviceMethod, err.Error())
	} else {
		log.Printf("%s %s success: %d + %d = %d", typ, serviceMethod, args.Num1, args.Num2, reply)
	}
}

//func call(addr chan string) {
//	client, _ := xRPC.XDial("http@" + <-addr)
//	defer func() { _ = client.Close() }()
//
//	time.Sleep(time.Second)
//
//	wg := new(sync.WaitGroup)
//
//	for i := 1; i <= 5; i++ {
//		wg.Add(1)
//
//		go func(i int) {
//			defer wg.Done()
//			args := &Args{i, i * i}
//			ctx, _ := context.WithTimeout(context.Background(), time.Second)
//			var reply int
//			if err := client.Call(ctx, "XK.Sum", args, &reply); err != nil {
//				log.Fatal("call XK.Sum err:", err)
//			}
//			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
//		}(i)
//	}
//	wg.Wait()
//}

func call(addr1, addr2 string) {
	d := xclient.NewMultiServersDiscovery([]string{"tcp@" + addr1, "tcp@" + addr2})
	xc := xclient.NewXClient(d, xclient.RandomSelect, nil)
	defer func() {
		_ = xc.Close()
	}()
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			xk(xc, context.Background(), "call", "XK.Sum", &Args{i, i * i})
		}(i)
	}
	wg.Wait()
}

func broadcast(addr1, addr2 string) {
	d := xclient.NewMultiServersDiscovery([]string{"tcp@" + addr1, "tcp@" + addr2})
	xc := xclient.NewXClient(d, xclient.RandomSelect, nil)
	defer func() {
		_ = xc.Close()
	}()
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			xk(xc, context.Background(), "broadcast", "XK.Sum", &Args{i, i * i})
			ctx, _ := context.WithCancel(context.Background())
			xk(xc, ctx, "broadcast", "XK.Sleep", &Args{i, i * i})
		}(i)
	}
	wg.Wait()
}

func main() {
	log.SetFlags(0)
	ch1 := make(chan string)
	ch2 := make(chan string)
	go StartServer(ch1)
	go StartServer(ch2)

	addr1 := <-ch1
	addr2 := <-ch2

	time.Sleep(time.Second * 2)

	call(addr1, addr2)
	// fmt.Println("broadcast")
	// TCP粘包
	broadcast(addr1, addr2)
}
