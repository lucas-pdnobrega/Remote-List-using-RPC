package main

import (
	"fmt"
	"remotelist/pkg"
)

var (
	ok bool
	val int
	list []int
)

func main() {
	client := remotelist.NewClientStub()

	fmt.Println("Inicializando duas listas...")
	client.CreateList()
	client.CreateList()


	fmt.Println("Realizando appends...")
	client.Append(0, 10)
	client.Append(0, 20)
	client.Append(1, 30)
	client.Append(1, 40)
	client.Append(1, 50)


	fmt.Println("Resgatando valores...")
	val = client.Get(0, 0)
	fmt.Println("Get[0][0]:", val)
	val = client.Get(0, 1)
	fmt.Println("Get[0][1]:", val)
	val = client.Get(1, 0)
	fmt.Println("Get[1][0]:", val)


	fmt.Println("\nTestando tamanho...")
	val = client.Size(0)
	fmt.Println("Tamanho da lista 0:", val)
	val = client.Size(1)
	fmt.Println("Tamanho da lista 1:", val)


	fmt.Println("\nRetirando último elemento da lista 0...")
	val = client.Remove(0)
	fmt.Println("Valor:", val)
	fmt.Println("\nRetirando último elemento da lista 1...")
	val = client.Remove(1)
	fmt.Println("Valor:", val)
}
