package db

type Sort struct {
	By    string
	IsAsc bool
}

type Range struct {
	Limit  uint64
	Offset uint64
}

type Option struct {
	Filters map[string][]string
	Sort    *Sort
	Range   *Range
}
