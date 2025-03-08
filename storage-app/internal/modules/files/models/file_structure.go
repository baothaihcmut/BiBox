package models

type FileNode struct {
	Value    *File
	SubFiles []*File
}

type FileStructure struct {
	Root []*FileNode
}
