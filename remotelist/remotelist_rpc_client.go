package main

import (
	"fmt"
	"net/rpc"
)

var (
	ok bool
	val int
	list []int
)

type AppendArgs struct {
	ListID int
	Value  int
}

type GetArgs struct {
	ListID int
	Index  int
}

func main() {
	client, err := rpc.Dial("tcp", ":5000")
	if err != nil {
		fmt.Print("dialing:", err)
	}

	var reply bool

	fmt.Println("Inicializando duas listas...")
	err = client.Call("RemoteList.CreateList", struct{}{}, &reply)
	err = client.Call("RemoteList.CreateList", struct{}{}, &reply)


	fmt.Println("Realizando appends...")
	err = client.Call("RemoteList.Append", &AppendArgs{0, 10}, &reply)
	err = client.Call("RemoteList.Append", &AppendArgs{0, 20}, &reply)
	err = client.Call("RemoteList.Append", &AppendArgs{1, 30}, &reply)
	err = client.Call("RemoteList.Append", &AppendArgs{1, 40}, &reply)
	err = client.Call("RemoteList.Append", &AppendArgs{1, 50}, &reply)


	fmt.Println("Resgatando valores...")
	_ = client.Call("RemoteList.Get", &GetArgs{0, 0}, &val)
	fmt.Println("Get[0][0]:", val)
	_ = client.Call("RemoteList.Get", &GetArgs{0, 1}, &val)
	fmt.Println("Get[0][1]:", val)
	_ = client.Call("RemoteList.Get", &GetArgs{1, 0}, &val)
	fmt.Println("Get[1][0]:", val)


	fmt.Println("\nTestando tamanho...")
	_ = client.Call("RemoteList.Size", 0, &val)
	fmt.Println("Tamanho da lista 0:", val)
	_ = client.Call("RemoteList.Size", 1, &val)
	fmt.Println("Tamanho da lista 1:", val)


	fmt.Println("\nRetirando último elemento da lista 0...")
	_ = client.Call("RemoteList.Remove", 0, &val)
	fmt.Println("Valor:", val)
	fmt.Println("\nRetirando último elemento da lista 1...")
	_ = client.Call("RemoteList.RemoveList", 1, &list)
	fmt.Println("Valor:", list)


	fmt.Println("\nSalvando...")
	_ = client.Call("RemoteList.CreateLogFile", struct{}{}, &ok)


	fmt.Println("\nLendo...")
	_ = client.Call("RemoteList.ReadLogFile", struct{}{}, &ok)
}
