package vm

//// Fun and fluff stuff
//type HWInfo struct {
//Name         string
//Manufacturer string
//Description  string
//}

type HWModule interface {
	ID() int
	//Query() HWInfo
	Interrupt(*VirtualMachine) error
}
