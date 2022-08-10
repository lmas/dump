package vm

import "fmt"

var (
	ErrEOF           = fmt.Errorf("end of file")
	ErrEOC           = fmt.Errorf("end of cycles")
	ErrHalt          = fmt.Errorf("halted")
	ErrZeroDivide    = fmt.Errorf("divide by zero")
	ErrStackEmpty    = fmt.Errorf("stack empty")
	ErrStackOverflow = fmt.Errorf("stack overflow")
)

type ErrUnknownOpCode struct {
	Op int
}

func (err ErrUnknownOpCode) Error() string {
	return fmt.Sprintf("unknown op code: %d", err.Op)
}
