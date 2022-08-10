package main

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/lmas/asm_game/lexer"
	"github.com/lmas/asm_game/opcodes"
	"github.com/lmas/asm_game/vm"
)

func main() {
	input := `// a test comment

// test instructions with weird white space around it
noop	    
	    noop

	// constants
	#define CONST -11111
	#define CO 11111
	movi a CONST
	movi b co

	// labels
	mov:
	a: b: C:

	// hardware interrupts
	hwq 0 // RNG
	hwq 1 // Printer
	hwi 0
	hwi 1

	// math
	movi a 1
	mov b a
	addi a 1
	add b a
	subi a 1
	sub b a
	muli a 2
	mul b a
	divi a 2
	div b a
	modi a 2
	mod b a

	// test instructions
	movi a 1
	teq a 0
	tne a 0
	tgt a 0
	tge a 0
	tlt a 0
	tle a 0

	// test registers
	movi a 1
	movi b 1
	movi c 1
	movi t 1
	movi x 1
	movi y 1

	// jumps
	jump true
	noop
true:
	movi t 1
	jpt false
	noop
false:
	movi t 0
	jpf call
	noop
call:
	call func
	noop
	//jump end
	jump end
	//noop

func:
	noop
	ret
	noop

end:
	//divi t 0
	//jump 200
	halt // all done, end program
	`
	// TODO: STILL CAUSING EOF WITH NO NEWLINES AT THE END OF INPUT!!!

	//input = `HWQ 19`
	//input = `call 7 call 0 call 0 call 6`

	logger := log.New(os.Stderr, "", 0)
	b := strings.NewReader(input)
	l := lexer.New(b)

	start := time.Now()
	bc, err := l.Parse()
	dur := time.Since(start)
	if err != nil {
		panic(err)
	}
	logger.Printf("Instruction list:\n%s", lexer.FormatInstructions(bc))
	logger.Printf("Parse duration: %s\n\n", dur)

	v := vm.New(
		vm.Conf{
			Frequency: 100,
			MaxCycles: 100,
			MaxStack:  100,
			Hardware: []vm.HWModule{
				&hwRNG{},
				//&hwDebugger{},
			},
			Logger: logger,
		},
	)
	v.Reset(bc)

	logger.Println("Running...")
	start = time.Now()
	err = v.Run()
	if err != nil {
		if err != vm.ErrEOF && err != vm.ErrHalt {
			panic(err)
		}
	}
	dur = time.Since(start)
	logger.Printf("Runtime:\t %s\n", dur)
	logger.Printf("Exit condition:\t %s\n", err)
}

////////////////////////////////////////////////////////////////////////////////

const (
	HW_RNG int = iota
	//HW_DEBUGGER
)

type hwRNG struct {
}

func (h hwRNG) ID() int {
	return HW_RNG
}

//func (h hwRNG) Query() vm.HWInfo {
//return vm.HWInfo{
////ID:           HW_RNG,
//Name:         "Nuclear RNG",
//Manufacturer: "Nuclear Appliances Co.",
//Description:  "Random number generator powered by a radioactive isotope",
//}
//}

func (h hwRNG) Interrupt(v *vm.VirtualMachine) error {
	rand.Seed(time.Now().UnixNano())
	min := v.Registers[opcodes.RA]
	max := v.Registers[opcodes.RB]
	val := rand.Intn(max+1-min) + min
	v.SetRegister(opcodes.RX, val)
	return nil
}

//type hwDebugger struct {
//}

//func (h hwDebugger) ID() int {
//return HW_DEBUGGER
//}

//func (h hwDebugger) Query() vm.HWInfo {
//return vm.HWInfo{
////ID:           HW_DEBUGGER,
//Name:         "Telebugger 101",
//Manufacturer: "Telebuggers Ltd.",
//Description:  "Remotely prints out state of currently running machine",
//}
//}

//func (h hwDebugger) Interrupt(v *vm.VirtualMachine) error {
//fmt.Printf("HW:%d\t MEM:%d\t CYC:%d\t PC:%d\t SP:%d\t RA:%d\t RB:%d\t RC:%d\t RT:%d\t RX:%d\t RY:%d\t\n",
//len(v.Conf.Hardware),
//len(v.Memory),
//v.Cycles,
//v.PC,
//len(v.Stack),
//v.Registers[opcodes.RA],
//v.Registers[opcodes.RB],
//v.Registers[opcodes.RC],
//v.Registers[opcodes.RT],
//v.Registers[opcodes.RX],
//v.Registers[opcodes.RY],
//)
//return nil
//}
