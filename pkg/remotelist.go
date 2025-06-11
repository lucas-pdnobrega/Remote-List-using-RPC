package remotelist

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type RemoteList struct {
	file_path string
	mu        sync.Mutex
	list      map[int][]int
	nextID    int
}

func (l *RemoteList) HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (l *RemoteList) startSnapshotRoutine(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			_ = l.SaveSnapshot()
		}
	}()
}

func (l *RemoteList) SaveSnapshot() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	data, err := json.MarshalIndent(l.list, "", "	")
	l.HandleErr(err)

	fmt.Printf("Saving snapshot file...\n")
	fileName := "data.json"
	path := fmt.Sprintf("%s/%s", l.file_path, fileName)
	err = os.WriteFile(path, data, 0644)
	l.HandleErr(err)

	fmt.Printf("Saved list to file: %s\n", path)
	return err
}

func (l *RemoteList) LoadSnapshot() error {
	fmt.Println("Loading snapshot file...")
	fileName := "data.json"
	path := fmt.Sprintf("%s/%s", l.file_path, fileName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		l.list = make(map[int][]int)
		return nil
	}

	data, err := os.ReadFile(path)
	l.HandleErr(err)

	l.mu.Lock()
	defer l.mu.Unlock()

	err = json.Unmarshal(data, &l.list)
	l.HandleErr(err)

	maxID := 0

	for id := range l.list {
		if id >= maxID {
			maxID = id + 1
		}
	}

	l.nextID = maxID

	fmt.Println("Successfully loaded list from file")
	fmt.Println(l.list)
	return nil
}

func (l *RemoteList) CreateList(_ *struct{}, reply *bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	detail := fmt.Sprintf("list with current id %d created", l.nextID)

	new_list := []int{}
	l.list[l.nextID] = new_list
	l.nextID++

	logOperation("CreateList", detail)
	fmt.Println(detail)
	fmt.Println(l.list)
	*reply = true
	return nil
}

func (l *RemoteList) RemoveList(Index int, reply *[]int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var detail string

	_, exists := l.list[Index]

	if !exists {
		detail = fmt.Sprintf("list with supplied id %d does not exist in records", Index)
		logOperation("RemoveList", detail)
		return errors.New(detail)
	}

	*reply = l.list[Index]
	delete(l.list, Index)

	detail = fmt.Sprintf("list with supplied id %d removed", Index)
	logOperation("RemoveList", detail)
	fmt.Println(detail)
	return nil
}

func (l *RemoteList) Get(args *GetArgs, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var detail string

	_, exists := l.list[args.ListID]

	if !exists {
		detail = fmt.Sprintf("list with supplied id %d does not exist in records", args.ListID)
		logOperation("Get", detail)
		return errors.New(detail)
	}

	if args.Index < 0 || args.Index >= len(l.list[args.ListID]) {
		detail = fmt.Sprintf("index %d out of list %d range", args.Index, args.ListID)
		logOperation("Get", detail)
		return errors.New(detail)
	}

	*reply = l.list[args.ListID][args.Index]

	detail = fmt.Sprintf("Get operation with supplied id %d", args.ListID)
	logOperation("Get", detail)
	fmt.Println(detail)

	return nil
}

func (l *RemoteList) Append(args *AppendArgs, reply *bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var detail string

	_, exists := l.list[args.ListID]

	if !exists {
		detail = fmt.Sprintf("list with supplied id %d does not exist in records", args.ListID)
		logOperation("Append", detail)
		return errors.New(detail)
	}

	sub_list := l.list[args.ListID]
	l.list[args.ListID] = append(sub_list, args.Value)

	detail = fmt.Sprintf("Append operation with supplied id %d and value %d", args.ListID, args.Value)
	logOperation("Append", detail)
	fmt.Println(detail)
	fmt.Println(l.list)

	*reply = true
	return nil
}

func (l *RemoteList) Remove(ListID int, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var detail string

	_, exists := l.list[ListID]

	if !exists {
		detail = fmt.Sprintf("list with supplied id %d does not exist in records", ListID)
		logOperation("Remove", detail)
		return errors.New(detail)
	}

	sub_list := l.list[ListID]

	if len(sub_list) > 0 {
		*reply = sub_list[len(sub_list)-1]
		l.list[ListID] = sub_list[:len(sub_list)-1]

		detail = fmt.Sprintf("Remove operation with supplied id %d successful", ListID)
		logOperation("Remove", detail)
		fmt.Println(detail)
		fmt.Println(l.list)

	} else {
		detail = fmt.Sprintf("Remove operation with supplied id %d is impossible because list is empty", ListID)
		logOperation("Remove", detail)
		return errors.New(detail)
	}
	return nil
}

func (l *RemoteList) Size(ListID int, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var detail string

	_, exists := l.list[ListID]

	if !exists {
		detail = fmt.Sprintf("list with supplied id %d does not exist in records", ListID)
		logOperation("Size", detail)
		return errors.New(detail)
	}

	*reply = len(l.list[ListID])
	detail = fmt.Sprintf("Size operation on list with supplied id %d successful", ListID)
	logOperation("Size", detail)
	fmt.Println(detail)

	return nil
}

func NewRemoteList() *RemoteList {
	l := &RemoteList{
		file_path: "./src",
		list:      make(map[int][]int),
		nextID:    0,
	}
	if err := l.LoadSnapshot(); err != nil {
		fmt.Printf("Failed to load snapshot: %v", err)
	}
	l.startSnapshotRoutine(5 * time.Minute)
	return l
}
