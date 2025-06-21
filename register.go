package main

type Register8 interface {
	Read() uint8
	Write(uint8)
	Increment() uint8
	Decrement() uint8
}

type Register16 interface {
	Read() uint16
	Write(uint16)
	Increment() uint16
	Decrement() uint16
	PostIncrement() uint16
	PostDecrement() uint16
}

type Registers struct {
	a Register8
	b Register8
	c Register8
	d Register8
	e Register8
	h Register8
	l Register8

	f *FlagRegister // Stores flags, not a real register

	af Register16
	bc Register16
	de Register16
	hl Register16

	sp Register16
	pc Register16
}

type SingleRegister8 struct {
	val uint8
}

type SingleRegister16 struct {
	val uint16
}

type CombinedRegister16 struct {
	hi Register8
	lo Register8
}

type FlagRegister struct {
	val uint8 // Only top four bits matter so operations use masks
}

func NewRegisters() *Registers {
	a := &SingleRegister8{}
	b := &SingleRegister8{}
	c := &SingleRegister8{}
	d := &SingleRegister8{}
	e := &SingleRegister8{}
	f := &FlagRegister{}
	h := &SingleRegister8{}
	l := &SingleRegister8{}

	af := &CombinedRegister16{hi: a, lo: f}
	bc := &CombinedRegister16{hi: b, lo: c}
	de := &CombinedRegister16{hi: d, lo: e}
	hl := &CombinedRegister16{hi: h, lo: l}

	sp := &SingleRegister16{}
	pc := &SingleRegister16{}

	return &Registers{
		a: a,
		b: b,
		c: c,
		d: d,
		e: e,
		f: f,
		h: h,
		l: l,

		af: af,
		bc: bc,
		de: de,
		hl: hl,

		sp: sp,
		pc: pc,
	}
}

func (r *SingleRegister8) Read() uint8 {
	return r.val
}

func (r *SingleRegister8) Write(val uint8) {
	r.val = val
}

func (r *SingleRegister8) Increment() uint8 {
	r.val = r.val + 1
	return r.val
}

func (r *SingleRegister8) Decrement() uint8 {
	r.val = r.val - 1
	return r.val
}

func (r *SingleRegister16) Read() uint16 {
	return r.val
}

func (r *SingleRegister16) Write(val uint16) {
	r.val = val
}

func (r *SingleRegister16) Increment() uint16 {
	r.val = r.val + 1
	return r.val
}

func (r *SingleRegister16) Decrement() uint16 {
	r.val = r.val - 1
	return r.val
}

func (r *SingleRegister16) PostIncrement() uint16 {
	old := r.val
	r.Increment()
	return old
}

func (r *SingleRegister16) PostDecrement() uint16 {
	old := r.val
	r.Decrement()
	return old
}

func (r *CombinedRegister16) Read() uint16 {
	return (uint16(r.hi.Read()) << 8) | uint16(r.lo.Read())
}

func (r *CombinedRegister16) Write(val uint16) {
	r.hi.Write(uint8(val >> 8))
	r.lo.Write(uint8(val & 0xFF))
}

func (r *CombinedRegister16) Increment() uint16 {
	r.Write(r.Read() + 1)
	return r.Read()
}

func (r *CombinedRegister16) Decrement() uint16 {
	r.Write(r.Read() - 1)
	return r.Read()
}

func (r *CombinedRegister16) PostIncrement() uint16 {
	old := r.Read()
	r.Increment()
	return old
}

func (r *CombinedRegister16) PostDecrement() uint16 {
	old := r.Read()
	r.Decrement()
	return old
}

func (f *FlagRegister) Read() uint8 {
	return f.val & 0xF0
}

func (f *FlagRegister) Write(val uint8) {
	f.val = val & 0xF0
}

// These methods should not be used. They are only here to abide by the Register8 interface
func (f *FlagRegister) Increment() uint8 {
	panic("Should not be incrementing a flag register")
}

func (f *FlagRegister) Decrement() uint8 {
	panic("Should not be decrementing a flag register")
}

func (f *FlagRegister) Z() bool {
	return f.val&(1<<7) != 0
}

func (f *FlagRegister) N() bool {
	return f.val&(1<<6) != 0
}

func (f *FlagRegister) H() bool {
	return f.val&(1<<5) != 0
}

func (f *FlagRegister) C() bool {
	return f.val&(1<<4) != 0
}

func (f *FlagRegister) SetZ(val bool) {
	if val {
		f.val |= 1 << 7
	} else {
		f.val &^= 1 << 7
	}
}

func (f *FlagRegister) SetN(val bool) {
	if val {
		f.val |= 1 << 6
	} else {
		f.val &^= 1 << 6
	}
}

func (f *FlagRegister) SetH(val bool) {
	if val {
		f.val |= 1 << 5
	} else {
		f.val &^= 1 << 5
	}
}

func (f *FlagRegister) SetC(val bool) {
	if val {
		f.val |= 1 << 4
	} else {
		f.val &^= 1 << 4
	}
}
