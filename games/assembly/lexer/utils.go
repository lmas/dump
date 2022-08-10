package lexer

import (
	"io"
	"strconv"
	"strings"

	"github.com/lmas/asm_game/opcodes"
)

func isSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\v', '\f', '\r': // Partly copied from the bufio pkg
		return true
	}
	return false
}

func (l *Lexer) err(msg string, val interface{}) error {
	return ErrLexer{l.line, msg, val}
}

func (l *Lexer) skip() error {
	for {
		b, err := l.input.ReadByte()
		if err != nil {
			return err
		}
		if !isSpace(b) {
			return l.input.UnreadByte()
		}
	}
	return nil
}

// Next returns the next word or single newline
func (l *Lexer) next() (string, error) {
	// Skip leading whitespace
	err := l.skip()
	if err != nil {
		return "", err
	}
	// Keep reading until next space
	var s strings.Builder
	for {
		b, err := l.input.ReadByte()
		if err != nil {
			// While fuzzing, found out we would miss the last word and
			// run into EOF if there wasn't any following char (like whitespace)
			if err == io.EOF {
				break
			}
			return "", err
		}
		if isSpace(b) {
			break
		}
		if b == '\n' {
			if s.Len() > 0 {
				if err := l.input.UnreadByte(); err != nil {
					return "", err
				}
				break
			}
			s.WriteRune('\n')
			l.line++
			break
		}
		s.WriteByte(b) // Docs says this will always return a nil error
	}
	return strings.ToUpper(s.String()), nil
}

// SkipLine will keep skipping tokens until the next newline
func (l *Lexer) skipLine() error {
	for {
		token, err := l.next()
		if err != nil {
			return err
		}
		if token == "\n" {
			break
		}
	}
	return nil
}

func (l *Lexer) filterUnwated(token string) error {
	switch {
	case token == "\n":
		return errSkipToken

	case strings.HasPrefix(token, COMMENT):
		// keep reading until next newline
		if err := l.skipLine(); err != nil {
			return err
		}
		return errSkipToken

	case strings.HasPrefix(token, MACRO):
		switch token {
		case MACRO_DEFINE:
			k, err := l.next()
			if err != nil {
				return err
			}
			v, err := l.nextImmidiate()
			if err != nil {
				return err
			}
			l.constants[k] = v
			return errSkipToken
		default:
			return l.err("unknown macro", token)
		}

	case strings.HasSuffix(token, LABEL):
		label := strings.TrimSuffix(token, LABEL)
		if _, found := l.labels[label]; found {
			return l.err("duplicate label", label)
		}
		l.labels[label] = l.index
		return errSkipToken

	default:
		return nil
	}
}

func (l *Lexer) nextRegister() (int, error) {
	t, err := l.next()
	if err != nil {
		return -1, err
	}
	reg := opcodes.StringToRegister(t)
	if reg < 0 {
		return -1, l.err("invalid register", t)
	}
	return reg, nil
}

func (l *Lexer) nextImmidiate() (int, error) {
	t, err := l.next()
	if err != nil {
		return -1, err
	}
	i, err := strconv.Atoi(t)
	if err == nil {
		return opcodes.ClampImmidiate(i), nil
	}
	i, ok := l.constants[t]
	if ok {
		return i, nil
	}
	return -1, l.err("invalid immidiate", t)
}

func (l *Lexer) checkJumps(bc []int) ([]int, error) {
	badRange := func(i int) bool {
		return i < 0 || i > len(bc)-1
	}

	// Find any labels that's "out of range" so to speak
	for label, jump := range l.labels {
		if badRange(jump) {
			// in theory, if a label has a bad range it's probably
			// at the very end of the instruction list so there's
			// no need to update the lineno. here
			return nil, l.err("invalid label", label)
		}
	}

	// Do a check on all unknown jump labels and see if we have found them yet
	for i, label := range l.unknownJumps {
		found := false
		for k, v := range l.labels {
			if k == label {
				bc[i+1] = v
				found = true
				break
			}
		}
		if !found {
			l.line = l.lineno[i]
			return nil, l.err("invalid label", label)
		}
	}

	// Make sure all jumps are in valid ranges
	for i := 0; i < len(bc)-1; i++ {
		switch bc[i] {
		case opcodes.JUMP, opcodes.JPT, opcodes.JPF, opcodes.CALL:
			// Found a panic(outofrange) here with the fuzzer :D
			j := bc[i+1]
			if badRange(j) {
				l.line = l.lineno[i]
				return nil, l.err("jump out of range", j)
			}
		}
	}

	return bc, nil
}
