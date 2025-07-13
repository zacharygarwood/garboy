package memory

type Register8 interface {
	Read() uint8
	Write(uint8)
}

type Register16 interface {
	Read() uint16
	Write(uint16)
	Increment() uint16
	Decrement() uint16
	PostIncrement() uint16
	PostDecrement() uint16
}
