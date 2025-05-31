package main

import (
	"fmt"
	"remotelist/pkg"
	"sync"
)

type AppendArgs struct {
	ListID int
	Value  int
}

type GetArgs struct {
	ListID int
	Index  int
}

type Void struct{}

func main() {
	client := remotelist.NewClientStub()
	var wg sync.WaitGroup

	fmt.Println("Inicializando listas...")
	client.CreateList()
	client.CreateList()

	fmt.Println("Iniciando testes concorrentes...")

	// Teste concorrente: múltiplos appends simultâneos
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			client.Append(0, v)
		}((i + 1) * 10)
	}

	// Teste concorrente: múltiplos gets simultâneos (depois de um pequeno delay)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			val := client.Get(0, index)
			fmt.Printf("[GET] Valor na posição %d: %d\n", index, val)
		}(i)
	}

	// Teste concorrente: múltiplos Size
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			size := client.Size(0)
			fmt.Println("[SIZE] Tamanho atual da lista 0:", size)
		}()
	}

	wg.Wait()
	fmt.Println("\nTestes concorrentes finalizados.")

	client.CreateLogFile()
}
