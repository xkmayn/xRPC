package main

import (
	"encoding/json"
	"fmt"
	xRPC "github.com/xkmayn/xrpc"
	"github.com/xkmayn/xrpc/codec"
	"log"
	"net"
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
	addr := make(chan string)
	go StartServer(addr)

	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(1)

	_ = json.NewEncoder(conn).Encode(xRPC.DefaultOption)
	cc := codec.NewGobCodec(conn)

	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "XKM:Sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("xrpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)

		log.Println("reply:", reply)
	}
}
