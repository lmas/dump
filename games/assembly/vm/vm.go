package vm

import (
	"fmt"
	"log"
	"time"

	"github.com/lmas/asm_game/opcodes"
)

type Conf struct {
	Frequency int // instructions / sec
	MaxCycles int
	MaxStack  int
	Hardware  []HWModule
	Logger    *log.Logger
}

type VirtualMachine struct {
	Conf Conf

	Memory    []int // Instruction memory
	Cycles    int
	PC        int // ProgramCounter
	LastOp    int
	Stack     []int
	Registers map[int]int
}

func New(c Conf) *VirtualMachine {
	v := &VirtualMachine{
		Conf: c,
	}
	v.Reset(nil)
	return v
}

func (v *VirtualMachine) Reset(bc []int) {
	if v.Conf.MaxStack < 1 {
		v.Conf.MaxStack = 100
	}

	v.Memory = bc
	v.Cycles = 0
	v.PC = -1
	v.LastOp = -1
	v.Stack = []int{}
	v.Registers = map[int]int{
		opcodes.RA: 0,
		opcodes.RB: 0,
		opcodes.RC: 0,
		opcodes.RT: 0,
		opcodes.RX: 0,
		opcodes.RY: 0,
	}
}

func (v *VirtualMachine) Run() error {
	var err error
	lastOp := -1
	skipTicks := 1000 / float64(v.Conf.Frequency)
	lastTick := now()
	for {
		err = v.Step()
		if lastOp != v.LastOp {
			lastOp = v.LastOp
			v.logLastOp()
		}
		if err != nil {
			break
		}
		// Clamp execution time (ops/sec) to the specified frequency
		c := now()
		skip := skipTicks - (c - lastTick)
		lastTick = c
		if skip > 0 {
			time.Sleep(time.Duration(skip) * time.Millisecond)
		}
	}
	return err
}

func (v *VirtualMachine) SetRegister(reg int, val int) {
	v.Registers[reg] = opcodes.ClampImmidiate(val)
}

//------------------------------------------------------------------------------

//func (v *VirtualMachine) log(s string, args ...interface{}) {
//if v.Conf.Logger != nil {
//v.Conf.Logger.Printf(s+"\n", args...)
//}
//}

// now returns the current unix time in milliseconds
func now() float64 {
	now := time.Now().UnixNano()
	return float64(now) / 1000 / 1000
}

func (v *VirtualMachine) logLastOp() {
	if v.Conf.Logger == nil {
		return
	}
	pc := v.LastOp
	op := v.Memory[pc]
	a, b := "", ""
	switch op {
	case opcodes.NOOP, opcodes.HALT, opcodes.RET:
		// pass
	case opcodes.HWQ, opcodes.HWI, opcodes.JUMP, opcodes.JPT, opcodes.JPF, opcodes.CALL:
		a = fmt.Sprint(v.Memory[pc+1])
	case opcodes.MOV, opcodes.ADD, opcodes.SUB, opcodes.MUL, opcodes.DIV, opcodes.MOD:
		a = fmt.Sprint(opcodes.RegisterToString(v.Memory[pc+1]))
		b = fmt.Sprint(opcodes.RegisterToString(v.Memory[pc+2]))
	case opcodes.MOVI, opcodes.ADDI, opcodes.SUBI, opcodes.MULI, opcodes.DIVI, opcodes.MODI, opcodes.TESTEQ,
		opcodes.TESTNE, opcodes.TESTGT, opcodes.TESTGE, opcodes.TESTLT, opcodes.TESTLE:
		a = fmt.Sprint(opcodes.RegisterToString(v.Memory[pc+1]))
		b = fmt.Sprint(v.Memory[pc+2])
	}
	v.Conf.Logger.Printf("%s\t %s\t %s\t CYC:%d\t PC:%d\t SC:%d\t RA:%d\t RB:%d\t RC:%d\t RT:%d\t RX:%d\t RY:%d\t\n",
		opcodes.OpToString(op),
		a,
		b,
		v.Cycles,
		pc,
		len(v.Stack),
		v.Registers[opcodes.RA],
		v.Registers[opcodes.RB],
		v.Registers[opcodes.RC],
		v.Registers[opcodes.RT],
		v.Registers[opcodes.RX],
		v.Registers[opcodes.RY],
	)
}

func (v *VirtualMachine) nextOperand() int {
	v.PC++
	return v.Memory[v.PC]
}

func (v *VirtualMachine) Step() error {
	v.Cycles++
	if v.Conf.MaxCycles > 0 && v.Cycles >= v.Conf.MaxCycles {
		return ErrEOC
	}
	v.PC++
	if v.PC >= len(v.Memory) {
		return ErrEOF
	}
	v.LastOp = v.PC
	op := v.Memory[v.PC]

	switch op {
	case opcodes.NOOP:
		// Do nothing

	case opcodes.HALT:
		return ErrHalt

	case opcodes.HWQ, opcodes.HWI:
		a := v.nextOperand()
		var hw HWModule
		for _, h := range v.Conf.Hardware {
			if h.ID() == a {
				hw = h
				break
			}
		}
		if hw == nil {
			v.Registers[opcodes.RT] = 0
			return nil
		}
		switch op {
		case opcodes.HWQ:
			v.SetRegister(opcodes.RX, hw.ID())
			v.Registers[opcodes.RT] = 1
		case opcodes.HWI:
			if err := hw.Interrupt(v); err != nil {
				return err
			}
			v.Registers[opcodes.RT] = 1
		}

	case opcodes.MOV, opcodes.ADD, opcodes.SUB, opcodes.MUL, opcodes.DIV, opcodes.MOD:
		a, b := v.nextOperand(), v.nextOperand()
		switch op {
		case opcodes.MOV:
			v.Registers[a] = v.Registers[b]
		case opcodes.ADD:
			v.Registers[a] += v.Registers[b]
		case opcodes.SUB:
			v.Registers[a] -= v.Registers[b]
		case opcodes.MUL:
			v.Registers[a] *= v.Registers[b]
		case opcodes.DIV:
			if v.Registers[b] == 0 {
				return ErrZeroDivide
			}
			v.Registers[a] /= v.Registers[b]
		case opcodes.MOD:
			if v.Registers[b] == 0 {
				return ErrZeroDivide
			}
			v.Registers[a] %= v.Registers[b]
		}
		v.SetRegister(a, v.Registers[a]) // clamps the value

	case opcodes.MOVI, opcodes.ADDI, opcodes.SUBI, opcodes.MULI, opcodes.DIVI, opcodes.MODI:
		a, b := v.nextOperand(), v.nextOperand()
		switch op {
		case opcodes.MOVI:
			v.Registers[a] = b
		case opcodes.ADDI:
			v.Registers[a] += b
		case opcodes.SUBI:
			v.Registers[a] -= b
		case opcodes.MULI:
			v.Registers[a] *= b
		case opcodes.DIVI:
			if b == 0 {
				return ErrZeroDivide
			}
			v.Registers[a] /= b
		case opcodes.MODI:
			if b == 0 {
				return ErrZeroDivide
			}
			v.Registers[a] %= b
		}
		v.SetRegister(a, v.Registers[a]) // clamps the value

	case opcodes.JUMP, opcodes.JPT, opcodes.JPF, opcodes.CALL:
		a := v.nextOperand()
		switch op {
		case opcodes.JUMP:
			v.PC = a - 1
		case opcodes.JPT:
			if v.Registers[opcodes.RT] > 0 {
				v.PC = a - 1
			}
		case opcodes.JPF:
			if v.Registers[opcodes.RT] < 1 {
				v.PC = a - 1
			}
		case opcodes.CALL:
			if len(v.Stack)+1 >= v.Conf.MaxStack {
				return ErrStackOverflow
			}
			v.Stack = append(v.Stack, v.PC)
			v.PC = a - 1
		}

	case opcodes.RET:
		if len(v.Stack) < 1 {
			return ErrStackEmpty
		}
		v.PC = v.Stack[0]
		v.Stack = v.Stack[1:]

	case opcodes.TESTEQ, opcodes.TESTNE, opcodes.TESTGT, opcodes.TESTGE, opcodes.TESTLT, opcodes.TESTLE:
		a, b, t := v.nextOperand(), v.nextOperand(), false
		switch op {
		case opcodes.TESTEQ:
			t = v.Registers[a] == b
		case opcodes.TESTNE:
			t = v.Registers[a] != b
		case opcodes.TESTGT:
			t = v.Registers[a] > b
		case opcodes.TESTGE:
			t = v.Registers[a] >= b
		case opcodes.TESTLT:
			t = v.Registers[a] < b
		case opcodes.TESTLE:
			t = v.Registers[a] <= b
		}
		if t == true {
			v.Registers[opcodes.RT] = 1
		} else {
			v.Registers[opcodes.RT] = 0
		}

	default:
		return ErrUnknownOpCode{op}
	}
	return nil
}
