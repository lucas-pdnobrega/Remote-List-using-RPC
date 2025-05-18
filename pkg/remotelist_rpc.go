package remotelist

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"encoding/json"
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

func (l *RemoteList) CreateLogFile(_ *struct{}, reply *bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	fmt.Printf("Creating log file...")
	fileName := "log.json"
	path := fmt.Sprintf("%s/%s", l.file_path, fileName)	

	data, err := json.MarshalIndent(l.list, "", "	")
	HandleErr(err)

	err = os.WriteFile(path, data, 0644)
	HandleErr(err)

	fmt.Printf("Saved list to file: %s\n", path)
	*reply = true
	return nil
}

func (l *RemoteList) ReadLogFile() {
	fmt.Println("Creating log file...")
	fileName := "log.json"
	path := fmt.Sprintf("%s/%s", l.file_path, fileName)
	
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		HandleErr(err)

		l.mu.Lock()
		defer l.mu.Unlock()

		err = json.Unmarshal(data, &l.list)
		HandleErr(err)

		l.size = uint32(len(l.list))

		fmt.Println("Successfully loaded list from file:")
		fmt.Println(l.list)
	} else if errors.Is(err, os.ErrNotExist) {
		data, err := json.MarshalIndent(l.list, "", "	")
		HandleErr(err)

		err = os.WriteFile(path, data, 0644)
		HandleErr(err)
		fmt.Printf("Created list on file: %s\n", path)
	}

}

func (l *RemoteList) CreateList(_ *struct{}, reply *bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	new_list := []int{}
	l.list = append(l.list, new_list)
	l.size++

	fmt.Println(l.list)
	*reply = true
	return nil
}

func (l *RemoteList) RemoveList( Index int, reply *[]int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if Index < 0 || Index >= len(l.list) {
		return errors.New("index out of range")
	}

	if l.size > 0 {
		*reply = l.list[Index]
		l.list = append(l.list[:Index], l.list[Index + 1:]...)
		l.size--
		fmt.Println(l.list)
	} else {
		return errors.New("empty list")
	}
	return nil
}


func (l *RemoteList) Get(args *GetArgs, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if args.ListID < 0 || args.ListID >= len(l.list) {
		return fmt.Errorf("index %d out of list range", args.ListID)
	}

	if args.Index < 0 || args.Index >= len(l.list[args.ListID]) {
		return fmt.Errorf("index %d out of list %d range", args.Index, args.ListID)
	}

	*reply = l.list[args.ListID][args.Index]

	return nil
}

func (l *RemoteList) Append(args *AppendArgs, reply *bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if args.ListID < 0 || args.ListID >= len(l.list) {
		return fmt.Errorf("index %d out of list range", args.ListID)
	}

	sub_list := l.list[args.ListID]

	l.list[args.ListID] = append(sub_list, args.Value)
	fmt.Println(l.list)
	*reply = true
	return nil
}

func (l *RemoteList) Remove(ListID int, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ListID < 0 || ListID >= len(l.list) {
		return fmt.Errorf("index %d out of list range", ListID)
	}

	sub_list := l.list[ListID]

	if len(sub_list) > 0 {
		*reply = sub_list[len(sub_list)-1]
		l.list[ListID] = sub_list[:len(sub_list)-1]
		fmt.Println(l.list)
	} else {
		return errors.New("empty list")
	}
	return nil
}

func (l *RemoteList) Size(ListID int, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ListID < 0 || ListID >= len(l.list) {
		return fmt.Errorf("index %d out of list range", ListID)
	}

	*reply = len(l.list[ListID])

	return nil
}

func NewRemoteList() *RemoteList {
	l := &RemoteList{
		file_path: "./src",
	}
	l.ReadLogFile()
	return l
}
