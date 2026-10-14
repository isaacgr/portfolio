package webserver

import "io/fs"

type File struct {
	path  string
	entry fs.DirEntry
}

type Error struct {
	Msg  string
	Code int
	Data string
}

func (e Error) Error() string {
	return e.Msg
}
