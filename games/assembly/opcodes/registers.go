package opcodes

const (
	RA int = iota
	RB
	RC
	RT
	RX
	RY
)

var RegisterMap = map[int]string{
	RA: "A",
	RB: "B",
	RC: "C",
	RT: "T",
	RX: "X",
	RY: "Y",
}

func RegisterToString(reg int) string {
	s, ok := RegisterMap[reg]
	if !ok {
		return ""
	}
	return s
}

func StringToRegister(s string) int {
	for k, v := range RegisterMap {
		if v == s {
			return k
		}
	}
	return -1
}

func ClampImmidiate(v int) int {
	if v < -9999 {
		return -9999
	} else if v > 9999 {
		return 9999
	}
	return v
}
