package fuzz

import (
	"bytes"
	"io"

	"github.com/lmas/asm_game/lexer"
	"github.com/lmas/asm_game/vm"
)

func Fuzz(data []byte) int {
	// In its basic form the Fuzz function just parses the input, and go-fuzz
	// ensures that it does not panic, crash the program, allocate insane amount
	// of memory nor hang.
	// The function must return:
	//	 1 if the fuzzer should increase priority of given input during
	//	   subsequent fuzzing (for example, the input is lexically correct
	//	   and was parsed successfully)
	//	-1 if the input must not be added to corpus even if gives new coverage
	//	 0 otherwise; other values are reserved for future use.
	// To communicate application-level bugs Fuzz function should panic (os.Exit(1)
	// will work too, but panic message contains more info). Note that Fuzz
	// function should not output to stdout/stderr, it will slow down fuzzing
	// and nobody will see the output anyway.
	// The exception is printing info about a bug just before panicking.

	b := bytes.NewReader(data)
	l := lexer.New(b)
	bc, err := l.Parse()

	switch err.(type) {
	case nil:
		//return 1
	case lexer.ErrLexer:
		return 0
	default:
		if err == io.EOF {
			break
		}
		panic(err)
	}

	///////////////////////////////////////////////////////////////////////

	v := vm.New(
		vm.Conf{
			Frequency: -1, // Prevents capping execution time
			MaxCycles: 100,
			MaxStack:  100,
			//Logger:    log.New(os.Stderr, "", 0),
		},
	)
	v.Reset(bc)
	err = v.Run()

	switch err {
	case nil:
		return 1
	case vm.ErrStackEmpty, vm.ErrStackOverflow, vm.ErrZeroDivide:
		return 0
	case vm.ErrEOF, vm.ErrHalt, vm.ErrEOC:
		return 0
	default:
		switch err.(type) {
		case vm.ErrUnknownOpCode:
			return 0
		}
		panic(err)
	}

}
