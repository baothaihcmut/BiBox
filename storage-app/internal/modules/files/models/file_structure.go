package models

type FileNode struct {
	Value    *File
	SubFiles []*FileNode
}

type FileStructure struct {
	Root []*FileNode
}
