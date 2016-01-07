package dl

type Library interface {
	Close() error

	Symbol(name string) (uintptr, error)
}

type Mode int

type Error struct {
	Message string
}

const (
	Lazy   Mode = 1 << 0
	Now    Mode = 1 << 1
	Global Mode = 1 << 2
	Local  Mode = 1 << 3
)

func Open(path string, mode Mode) (Library, error) {
	return open(path, mode)
}

func (err *Error) Error() string {
	return err.Message
}
