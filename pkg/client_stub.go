package remotelist

import (
	"fmt"
	"net/rpc"
)


type ClientStub struct {
	client *rpc.Client
	reply bool
	err error
}

func NewClientStub() *ClientStub{
	client, err := rpc.Dial("tcp", ":5000")
	if err != nil {
		fmt.Print("dialing:", err)
	}
	c := &ClientStub{
		client: client,
	}
	return c
}

func (c *ClientStub) HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (c *ClientStub) CreateList() bool {
	c.err = c.client.Call("RemoteList.CreateList", &Void{}, &c.reply)
	fmt.Println(c.reply)
	c.reply = true
	return c.reply
}

func (c *ClientStub) RemoveList(Index int) []int {
	var val []int
	_ = c.client.Call("RemoteList.Remove", 0, &val)
	fmt.Println("Valor:", val)
	return val
}

func (c *ClientStub) Get(ListID int, Index int) int{
	var val int
	_ = c.client.Call("RemoteList.Get", &GetArgs{ListID, Index}, &val)
	fmt.Println("Get:", val)
	return val
}

func (c *ClientStub) Append(ListID int, Value int) bool {
	_ = c.client.Call("RemoteList.Append", &AppendArgs{ListID, Value}, &c.reply)
	fmt.Println(c.reply)
	c.reply = true
	return c.reply
}

func (c *ClientStub) Remove(ListID int) int {
	var val int
	_ = c.client.Call("RemoteList.Remove", ListID, &val)
	fmt.Println("Valor:", val)
	return val
}

func (c *ClientStub) Size(ListID int) int {
	var val int
	_ = c.client.Call("RemoteList.Size", ListID, &val)
	fmt.Println("Valor:", val)
	return val
}
