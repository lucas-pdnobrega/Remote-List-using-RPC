package remotelist

type AppendArgs struct {
	ListID int
	Value  int
}

type GetArgs struct {
	ListID int
	Index  int
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Operation string `json:"operation"`
	Details   string `json:"details"`
}

type Void struct{}
