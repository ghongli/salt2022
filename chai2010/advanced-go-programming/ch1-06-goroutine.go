package advanced_go_programming_book

import (
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/gatefs"
)

func aboutgatefs() {
	gatefs.New(vfs.OS("/path"), make(chan bool, 8))
}
