package files

type FileItem struct {
	IsDir bool
	Name  string
	Path  string
	Size  int64
	Mtime int64
}
