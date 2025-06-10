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

	fmt.Printf("Saving snapshot file...")
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

	fmt.Println("Successfully loaded list from file:")
	fmt.Println(l.list)
	return nil
}

func (l *RemoteList) CreateList(_ *struct{}, reply *bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	new_list := []int{}
	l.list[l.nextID] = new_list
	l.nextID++

	fmt.Println(l.list)
	*reply = true
	return nil
}

func (l *RemoteList) RemoveList(Index int, reply *[]int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	_, exists := l.list[Index]

	if !exists {
		return errors.New("list with supplied id does not exist in records")
	}

	*reply = l.list[Index]
	delete(l.list, Index)

	fmt.Println(l.list)
	return nil
}

func (l *RemoteList) Get(args *GetArgs, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	_, exists := l.list[args.ListID]

	if !exists {
		return errors.New("list with supplied id does not exist in records")
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

	_, exists := l.list[args.ListID]

	if !exists {
		return errors.New("list with supplied id does not exist in records")
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

	_, exists := l.list[ListID]

	if !exists {
		return errors.New("list with supplied id does not exist in records")
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

	_, exists := l.list[ListID]

	if !exists {
		return errors.New("list with supplied id does not exist in records")
	}

	*reply = len(l.list[ListID])

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
