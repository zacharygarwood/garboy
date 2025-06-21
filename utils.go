package main

func ExtractBits(b byte, positions []int) byte {
	var result byte
	for i, pos := range positions {
		bit := (b >> pos) & 1
		result |= bit << (len(positions) - 1 - i)
	}
	return result
}

// Ref: https://gist.github.com/meganesu/9e228b6b587decc783aa9be34ae27841
func IsHalfCarry8(a byte, b byte) bool {
	return (((a & 0xF) + (b & 0xF)) & 0x10) == 0x10
}

func IsHalfCarry16(a uint16, b uint16) bool {
	return (((a & 0xFFF) + (b & 0xFFF)) & 0x1000) == 0x1000
}

func IsHalfCarryWithCarry8(a byte, b byte, c byte) bool {
	return ((a & 0xF) + (b & 0xF) + c) > 0xF
}

func IsHalfCarryWithCarry16(a uint16, b uint16, c uint16) bool {
	return ((a & 0xFFF) + (b & 0xFFF) + c) > 0xFFF
}

func IsHalfBorrow8(a uint8, b uint8) bool {
	return (a & 0xF) < (b & 0xF)
}

func IsHalfBorrow16(a uint16, b uint16) bool {
	return (a & 0xFFF) < (b & 0xFFF)
}

func IsHalfBorrowWithCarry8(a uint8, b uint8, c uint8) bool {
	return (a & 0xF) < ((b & 0xF) + c)
}

func IsHalfBorrowWithCarry16(a uint16, b uint16, c uint16) bool {
	return (a & 0xFFF) < ((b & 0xFFF) + c)
}

// Flags need to be uint8 for some instructions
func AsUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func RotateRight(val uint8) (uint8, bool) {
	lsb := val & 0x01
	result := (val << 1) | lsb
	carry := (result & 0x01) != 0
	return result, carry
}

func RotateLeft(val uint8) (uint8, bool) {
	msb := (val & 0x80) >> 7
	result := (val << 1) | msb
	carry := (result & 0x01) != 0
	return result, carry
}

func RotateRightThroughCarry(val uint8, carry bool) (uint8, bool) {
	lsb := val & 0x01
	result := val >> 1
	if carry {
		result |= 0x80
	}
	newCarry := lsb != 0
	return result, newCarry
}

func RotateLeftThroughCarry(val uint8, carry bool) (uint8, bool) {
	msb := (val & 0x80) != 0
	result := val << 1
	if carry {
		result |= 0x01
	}
	newCarry := msb
	return result, newCarry
}

func ShiftLeftArithmetic(val uint8) (uint8, bool) {
	carry := (val & 0x80) != 0
	result := val << 1
	return result, carry
}

func ShiftRightArithmetic(val uint8) (uint8, bool) {
	carry := (val & 0x01) != 0
	result := (val >> 1) | (val & 0x80)
	return result, carry
}

func ShiftRightLogic(val uint8) (uint8, bool) {
	carry := (val & 0x01) != 0
	result := val >> 1
	return result, carry
}

// Swaps hi and lo bits. Ex: hi-lo becomes lo-hi
func Swap(val uint8) uint8 {
	return (val >> 4) + ((val & 0xF) << 4)
}

func IsBitSet(val uint8, bit uint8) bool {
	return (val & (1 << bit)) != 0
}

func ResetBit(val uint8, bit uint8) uint8 {
	return val & ^(1 << bit)
}

func SetBit(val uint8, bit uint8) uint8 {
	return val | (1 << bit)
}
