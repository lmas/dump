package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/lmas/asm_game/opcodes"
)

const (
	COMMENT = "//"
	MACRO   = "#"
	LABEL   = ":"

	MACRO_DEFINE = "#DEFINE"

	MAX_SIZE int64 = 1000000 // 1MB
)

var errSkipToken = fmt.Errorf("skip")

type ErrLexer struct {
	Line int
	Msg  string
	Val  interface{}
}

func (err ErrLexer) Error() string {
	return fmt.Sprintf("line %d, %s: %v", err.Line, err.Msg, err.Val)
}

type Lexer struct {
	input        *bufio.Reader
	index        int
	line         int
	labels       map[string]int
	constants    map[string]int
	unknownJumps map[int]string
	lineno       map[int]int
}

func New(input io.Reader) *Lexer {
	l := &Lexer{
		input:        bufio.NewReader(io.LimitReader(input, MAX_SIZE)),
		line:         1,
		labels:       make(map[string]int),
		constants:    make(map[string]int),
		unknownJumps: make(map[int]string),
		lineno:       make(map[int]int),
	}
	return l
}

func (l *Lexer) Parse() ([]int, error) {
	var bc []int

	// First loop should find all constants and labels
	constants := make(map[string]int)
	labels := make(map[string]int)
	var instructions []string
	lines := -1
	s := bufio.NewScanner(l.input)
	for s.Scan() {
		l := s.Text()
		i := strings.Index(l, COMMENT)
		if i > -1 {
			l = l[:i]
		}
		l = strings.ToUpper(strings.TrimSpace(l))
		if len(l) < 1 {
			continue
		}

		if strings.HasPrefix(l, MACRO) {
			parts := strings.Split(l, " ")
			op, a, b := parts[0], "", ""
			if len(parts) > 1 {
				a = parts[1]
			}
			if len(parts) > 2 {
				b = parts[2]
			}
			switch op {
			case MACRO_DEFINE:
				if !unicode.IsLetter(rune(a[0])) {
					return nil, fmt.Errorf("invalid constant name: %s", l)
				}
				i, err := strconv.Atoi(b)
				if err != nil {
					return nil, fmt.Errorf("invalid constant value: %s", l)
				}
				constants[a] = i
			default:
				return nil, fmt.Errorf("unknown macro: %s", parts[0])
			}
			continue
		}

		i = strings.Index(l, LABEL)
		if i > -1 {
			// TODO: handle the rest of the non-empty label lines
			label := l[:i]
			labels[label] = lines + 1
			continue
		}

		lines++
		//fmt.Println(lines, l)
		instructions = append(instructions, l)
	}
	err := s.Err()
	if err != nil {
		return nil, err
	}
	fmt.Println(constants)
	fmt.Println(labels)

	// 2nd loop does the actual parsing, ignoring whitespace/comments
	for lineno, line := range instructions {
		parts := strings.Split(line, " ")
		token, a, b := parts[0], "", ""
		if len(parts) > 1 {
			a = parts[1]
		}
		if len(parts) > 2 {
			b = parts[2]
		}
		c, found := constants[b]
		if found {
			// TODO: try to skip this conversion altogether
			b = strconv.Itoa(c)
		}
		fmt.Println(lineno, token, a, b)

		// Parse the instruction
		op := opcodes.StringToOp(token)
		instruction := []int{op}
		switch op {
		case opcodes.NOOP, opcodes.HALT, opcodes.RET:
			// Do nothing

		case opcodes.HWQ, opcodes.HWI:
			i, err := strconv.Atoi(a)
			if err != nil {
				return nil, fmt.Errorf("invalid immidiate value for field A: %s", line)
			}
			instruction = append(instruction, i)

		case opcodes.MOV, opcodes.ADD, opcodes.SUB, opcodes.MUL, opcodes.DIV, opcodes.MOD:
			to := opcodes.StringToRegister(a)
			if to < 0 {
				return nil, fmt.Errorf("invalid register for field A: %s", line)
			}
			from := opcodes.StringToRegister(b)
			if from < 0 {
				return nil, fmt.Errorf("invalid register for field B: %s", line)
			}
			instruction = append(instruction, []int{to, from}...)

		case opcodes.MOVI, opcodes.ADDI, opcodes.SUBI, opcodes.MULI, opcodes.DIVI, opcodes.MODI, opcodes.TESTEQ,
			opcodes.TESTNE, opcodes.TESTGT, opcodes.TESTGE, opcodes.TESTLT, opcodes.TESTLE:
			r := opcodes.StringToRegister(a)
			if r < 0 {
				return nil, fmt.Errorf("invalid register for field A: %s", line)
			}
			i, err := strconv.Atoi(b)
			if err != nil {
				return nil, fmt.Errorf("invalid immidiate value for field B: %s", line)
			}
			instruction = append(instruction, []int{r, i}...)

		case opcodes.JUMP, opcodes.JPT, opcodes.JPF, opcodes.CALL:
			jump, found := labels[a]
			if !found {
				return nil, fmt.Errorf("unknown label: %s", line)
			}
			instruction = append(instruction, jump)

		default:
			return nil, l.err("unknown instruction", token)
		}
		bc = append(bc, instruction...)
		//l.lineno[l.index] = l.line
		//l.index += len(instruction)
	}
	return bc, nil
}

func (l *Lexer) Parse1() ([]int, error) {
	var bc []int
	for {
		token, err := l.next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Filter out unwanted stuff
		err = l.filterUnwated(token)
		if err != nil {
			if err == errSkipToken {
				continue
			}
			return nil, err
		}

		// Parse the instruction
		op := opcodes.StringToOp(token)
		instruction := []int{op}
		switch op {
		case opcodes.NOOP, opcodes.HALT, opcodes.RET:
			// Do nothing

		case opcodes.HWQ, opcodes.HWI:
			a, err := l.nextImmidiate()
			if err != nil {
				return nil, err
			}
			instruction = append(instruction, a)

		case opcodes.MOV, opcodes.ADD, opcodes.SUB, opcodes.MUL, opcodes.DIV, opcodes.MOD:
			a, err := l.nextRegister()
			if err != nil {
				return nil, err
			}
			b, err := l.nextRegister()
			if err != nil {
				return nil, err
			}
			instruction = append(instruction, []int{a, b}...)

		case opcodes.MOVI, opcodes.ADDI, opcodes.SUBI, opcodes.MULI, opcodes.DIVI, opcodes.MODI, opcodes.TESTEQ,
			opcodes.TESTNE, opcodes.TESTGT, opcodes.TESTGE, opcodes.TESTLT, opcodes.TESTLE:
			a, err := l.nextRegister()
			if err != nil {
				return nil, err
			}
			b, err := l.nextImmidiate()
			if err != nil {
				return nil, err
			}
			instruction = append(instruction, []int{a, b}...)

		case opcodes.JUMP, opcodes.JPT, opcodes.JPF, opcodes.CALL:
			label, err := l.next() // only whole words (not immidiate vals) allowed
			if err != nil {
				return nil, err
			}
			// Quick check if we already know the label (probably not)
			i, found := l.labels[label]
			if !found {
				i = -1
				l.unknownJumps[l.index] = label // gotta find it later
			}
			instruction = append(instruction, i)

		default:
			return nil, l.err("unknown instruction", token)
		}
		bc = append(bc, instruction...)
		l.lineno[l.index] = l.line
		l.index += len(instruction)
	}

	// Try fixing or detecting any broken jumps and labels
	var err error
	bc, err = l.checkJumps(bc)
	if err != nil {
		return nil, err
	}
	return bc, nil
}

func FormatInstructions(bc []int) string {
	var s strings.Builder
	for i := 0; i < len(bc); i++ {
		s.WriteString(fmt.Sprintf("#%d\t %s", i, opcodes.OpToString(bc[i])))
		switch bc[i] {
		case opcodes.NOOP, opcodes.HALT, opcodes.RET:
		case opcodes.HWQ, opcodes.HWI, opcodes.JUMP, opcodes.JPT, opcodes.JPF, opcodes.CALL,
			opcodes.TESTEQ, opcodes.TESTNE, opcodes.TESTGT, opcodes.TESTGE, opcodes.TESTLT, opcodes.TESTLE:
			i++
			s.WriteString(fmt.Sprintf("\t%v", bc[i]))
		case opcodes.MOV, opcodes.ADD, opcodes.SUB, opcodes.MUL, opcodes.DIV, opcodes.MOD,
			opcodes.MOVI, opcodes.ADDI, opcodes.SUBI, opcodes.MULI, opcodes.DIVI, opcodes.MODI:
			i++
			s.WriteString(fmt.Sprintf("\t%v", bc[i]))
			i++
			s.WriteString(fmt.Sprintf("\t%v", bc[i]))
		}
		s.WriteString("\n")
	}
	return s.String()
}
