package main

type Register8 struct {
	val uint8
}

type Register16 struct {
	val uint16
}

type CombinedRegister16 struct {
	hi *Register8
	lo *Register8
}

type Registers struct {
	a *Register8
	b *Register8
	c *Register8
	d *Register8
	e *Register8
	f *Register8 // Stores flags, not a real register
	h *Register8
	l *Register8

	af *CombinedRegister16
	bc *CombinedRegister16
	de *CombinedRegister16
	hl *CombinedRegister16

	sp *Register16
	pc *Register16
}

func NewRegisters() *Registers {
	a := &Register8{}
	b := &Register8{}
	c := &Register8{}
	d := &Register8{}
	e := &Register8{}
	f := &Register8{}
	h := &Register8{}
	l := &Register8{}

	af := &CombinedRegister16{hi: a, lo: f}
	bc := &CombinedRegister16{hi: b, lo: c}
	de := &CombinedRegister16{hi: d, lo: e}
	hl := &CombinedRegister16{hi: h, lo: l}

	sp := &Register16{}
	pc := &Register16{}

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

func (r Register8) read() uint8 {
	return r.val
}

func (r *Register8) write(val uint8) {
	r.val = val
}

func (r *Register16) read() uint16 {
	return r.val
}

func (r *Register16) write(val uint16) {
	r.val = val
}
