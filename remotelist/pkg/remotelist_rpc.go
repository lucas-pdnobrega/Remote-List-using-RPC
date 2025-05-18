package remotelist

import (
	"errors"
	"fmt"
	"sync"
)

type RemoteList struct {
	mu   sync.Mutex
	list []int
	size uint32
}

func (l *RemoteList) Get(value int, index int, reply *int) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if index < len(l.list) {
		return l.list[index], nil
	} else {
		return -1, errors.New("index out of list range")
	}
}

func (l *RemoteList) Append(value int, reply *bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.list = append(l.list, value)
	fmt.Println(l.list)
	l.size++
	*reply = true
	return nil
}

func (l *RemoteList) Remove(arg int, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.list) > 0 {
		*reply = l.list[len(l.list)-1]
		l.list = l.list[:len(l.list)-1]
		fmt.Println(l.list)
	} else {
		return errors.New("empty list")
	}
	return nil
}

func (l *RemoteList) Size(arg int, reply *int) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.list)
}

func NewRemoteList() *RemoteList {
	return new(RemoteList)
}
