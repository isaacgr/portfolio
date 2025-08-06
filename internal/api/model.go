package api

type Error struct {
	Msg  string
	Code int
	Data string
}

func (e Error) Error() string {
	return e.Msg
}
