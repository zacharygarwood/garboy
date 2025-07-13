package interrupts

type InterruptRegister struct {
	val byte
}

func (i *InterruptRegister) Read() uint8 {
	return i.val
}

func (i *InterruptRegister) Write(val uint8) {
	i.val = val
}
