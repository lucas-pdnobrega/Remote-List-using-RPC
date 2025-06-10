package remotelist


type AppendArgs struct {
	ListID int
	Value  int
}

type GetArgs struct {
	ListID int
	Index  int
}

type Void struct{}
