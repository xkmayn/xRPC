package main

import (
	"fmt"
	xRPC "github.com/xkmayn/xrpc"
	"log"
	"net"
	"sync"
	"time"
)

func StartServer(addr chan string) {
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

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("xrpc req %d", i)
			var reply string
			if err := client.Call("XKM.Sum", args, &reply); err != nil {
				log.Fatal("call XKM.Sum err:", err)
			}
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
}
