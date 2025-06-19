package main

func ExtractBits(b byte, positions []int) byte {
	var result byte
	for i, pos := range positions {
		bit := (b >> pos) & 1
		result |= bit << (len(positions) - 1 - i)
	}
	return result
}
