package memory

type Memory interface {
	Read(offset uint16) byte
	Write(offset uint16, val byte)
}

type RAM struct {
	data []byte
}

func NewRAM(size int) *RAM {
	return &RAM{
		data: make([]byte, size),
	}
}

func (r *RAM) Read(offset uint16) byte {
	return r.data[offset]
}

func (r *RAM) Write(offset uint16, val byte) {
	r.data[offset] = val
}

type ROM struct {
	data []byte
}

func NewROM(data []byte) *ROM {
	return &ROM{
		data: data,
	}
}

func (r *ROM) Read(offset uint16) byte {
	return r.data[offset]
}

func (r *ROM) Write(offset uint16, val byte) {
	panic("Should not be writing to ROM")
}

type IORegisters struct {
	data [0x80]byte // FF00-FF7F
}

func NewIORegisters() *IORegisters {
	return &IORegisters{
		data: [0x80]byte{},
	}
}

func (io *IORegisters) Read(offset uint16) byte {
	return io.data[offset]
}

func (io *IORegisters) Write(offset uint16, val byte) {
	io.data[offset] = val
	// TODO: Trigger hardware effects (timer, joypad, etc)
}

// Not the greatest, but I had to add this so that the byte in memory returned by cpu.byteAt() would be treated as a Register8
// Ex: INC (HL) calls inc_r8 which takes Register8. Called like inc_r8(c.byteAt(c.reg.hl.Read()))
type MemoryReference8 struct {
	Mmu  Memory
	Addr uint16
}

func (m *MemoryReference8) Read() uint8 {
	return m.Mmu.Read(m.Addr)
}

func (m *MemoryReference8) Write(val uint8) {
	m.Mmu.Write(m.Addr, val)
}

func (m *MemoryReference8) Increment() uint8 {
	res := m.Mmu.Read(m.Addr) + 1
	m.Mmu.Write(m.Addr, res)
	return res
}

func (m *MemoryReference8) Decrement() uint8 {
	res := m.Mmu.Read(m.Addr) - 1
	m.Mmu.Write(m.Addr, res)
	return res
}
