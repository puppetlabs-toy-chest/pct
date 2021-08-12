package mock

import (
	"io/fs"
	"time"
)

type FileInfo struct {
	FName  string
	FSize  int64
	FMode  fs.FileMode
	FTime  time.Time
	FIsDir bool
	FSys   interface{}
}

func (f *FileInfo) Name() string {
	return f.FName
}

func (f *FileInfo) Size() int64 {
	return f.FSize
}

func (f *FileInfo) Mode() fs.FileMode {
	return f.FMode
}

func (f *FileInfo) ModTime() time.Time {
	return f.FTime
}

func (f *FileInfo) IsDir() bool {
	return f.FIsDir
}

func (f *FileInfo) Sys() interface{} {
	return f.FSys
}
