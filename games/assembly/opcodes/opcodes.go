package opcodes

const (
	// Misc.
	NOOP int = iota
	HALT
	MOV
	MOVI
	HWQ // TODO: add register version?
	HWI

	// Math
	ADD
	ADDI
	SUB
	SUBI
	MUL
	MULI
	DIV
	DIVI
	MOD
	MODI

	// Jumps
	JUMP
	JPT
	JPF
	CALL
	RET

	// Tests
	TESTEQ // TODO: add register version?
	TESTNE
	TESTGT
	TESTGE
	TESTLT
	TESTLE
)

var OpMap = map[int]string{
	NOOP: "NOOP",
	HALT: "HALT",
	MOV:  "MOV",
	MOVI: "MOVI",
	HWQ:  "HWQ",
	HWI:  "HWI",

	ADD:  "ADD",
	ADDI: "ADDI",
	SUB:  "SUB",
	SUBI: "SUBI",
	MUL:  "MUL",
	MULI: "MULI",
	DIV:  "DIV",
	DIVI: "DIVI",
	MOD:  "MOD",
	MODI: "MODI",

	JUMP: "JUMP",
	JPT:  "JPT",
	JPF:  "JPF",
	CALL: "CALL",
	RET:  "RET",

	TESTEQ: "TEQ",
	TESTNE: "TNE",
	TESTGT: "TGT",
	TESTGE: "TGE",
	TESTLT: "TLT",
	TESTLE: "TLE",
}

func OpToString(op int) string {
	s, ok := OpMap[op]
	if !ok {
		return ""
	}
	return s
}

func StringToOp(s string) int {
	for k, v := range OpMap {
		if v == s {
			return k
		}
	}
	return -1
}
