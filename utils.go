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

func IsHalfBorrow8(a byte, b byte) bool {
	return ((a & 0xF) - (b & 0xF)) < 0
}

func IsHalfBorrow16(a uint16, b uint16) bool {
	return ((a & 0xFFF) - (b & 0xFFF)) < 0
}
