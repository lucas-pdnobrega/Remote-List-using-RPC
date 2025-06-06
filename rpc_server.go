package main

import (
	"fmt"
	"net"
	"net/rpc"
	"remotelist/pkg"
)

func main() {
	list := remotelist.NewRemoteList()
	rpcs := rpc.NewServer()
	err := rpcs.RegisterName("RemoteList", list)
	if err != nil {
		fmt.Println("Erro ao registrar RemoteList:", err)
	}

	l, e := net.Listen("tcp", "[localhost]:5000")
	defer l.Close()
	if e != nil {
		fmt.Println("listen error:", e)
	}
	for {
		conn, err := l.Accept()
		if err == nil {
			go rpcs.ServeConn(conn)
		} else {
			break
		}
	}
}
