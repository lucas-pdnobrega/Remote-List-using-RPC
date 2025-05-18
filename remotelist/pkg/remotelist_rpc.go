package remotelist

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"encoding/json"
)

type RemoteList struct {
	file_path string
	mu   sync.Mutex
	list [][]int
	size uint32
}

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateLogFile(l *RemoteList) {
	fmt.Printf("Creating log file...")
	fileName := "log.txt"
	path := fmt.Sprintf("%s/%s", l.file_path, fileName)
	
	l.mu.Lock()
	defer l.mu.Unlock()

	data, err := json.MarshalIndent(l.list, "", "	")
	HandleErr(err)

	err = os.WriteFile(path, data, 0644)
	HandleErr(err)

	fmt.Printf("Saved list to file: %s\n", path)
}

func ReadLogFile(l *RemoteList) {
	fmt.Printf("Creating log file...")
	fileName := "log.txt"
	path := fmt.Sprintf("%s/%s", l.file_path, fileName)
	
	data, err := os.ReadFile(path)
	HandleErr(err)

	l.mu.Lock()
	defer l.mu.Unlock()

	err = json.Unmarshal(data, &l.list)
	HandleErr(err)

	l.size = uint32(len(l.list))

	fmt.Println("Successfully loaded list from file:")
	fmt.Println(l.list)
}

func (l *RemoteList) CreateList( reply *bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	new_list := []int{}
	l.list = append(l.list, new_list)
	l.size++

	fmt.Println(l.list)
	*reply = true
	return nil
}

func (l *RemoteList) RemoveList( index int, reply *[]int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if index < 0 || index >= len(l.list) {
		return errors.New("index out of range")
	}

	if l.size > 0 {
		*reply = l.list[index]
		l.list = append(l.list[:index], l.list[index + 1:]...)
		l.size--
		fmt.Println(l.list)
	} else {
		return errors.New("empty list")
	}
	return nil
}


func (l *RemoteList) Get(list_id int, item_index int, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if list_id < 0 || list_id >= len(l.list) {
		return fmt.Errorf("index %d out of list range", list_id)
	}

	if item_index < 0 || item_index >= len(l.list[list_id]) {
		return fmt.Errorf("index %d out of list %d range", item_index, list_id)
	}

	*reply = l.list[list_id][item_index]

	return nil
}

func (l *RemoteList) Append(list_id int, v int, reply *bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if list_id < 0 || list_id >= len(l.list) {
		return fmt.Errorf("index %d out of list range", list_id)
	}

	sub_list := l.list[list_id]

	l.list[list_id] = append(sub_list, v)
	fmt.Println(l.list)
	*reply = true
	return nil
}

func (l *RemoteList) Remove(list_id int, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if list_id < 0 || list_id >= len(l.list) {
		return fmt.Errorf("index %d out of list range", list_id)
	}

	sub_list := l.list[list_id]

	if len(sub_list) > 0 {
		*reply = sub_list[len(sub_list)-1]
		l.list[list_id] = sub_list[:len(sub_list)-1]
		fmt.Println(l.list)
	} else {
		return errors.New("empty list")
	}
	return nil
}

func (l *RemoteList) Size(list_id int, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if list_id < 0 || list_id >= len(l.list) {
		return fmt.Errorf("index %d out of list range", list_id)
	}

	*reply = len(l.list[list_id])

	return nil
}

func NewRemoteList() *RemoteList {
	l := &RemoteList{
		file_path: "../src/",
	}
	return l
}
