package cartridge

const (
	RomBankSize = 0x4000
	RamBankSize = 0x2000
)

type MBC interface {
	Read(addr uint16) uint8
	Write(addr uint16, value uint8)
}

func NewMBC(rom []uint8, header CartridgeHeader) MBC {
	switch header.CartType {
	case 0x00:
		return NewMBC0(rom, header)
	case 0x01, 0x02, 0x03:
		return NewMBC1(rom, header)
	case 0x0F, 0x10, 0x11, 0x12, 0x13:
		return NewMBC3(rom, header)
	default:
		panic("Unsupported MBC type")
	}
}

func getRomSize(romSizeCode uint8) int {
	switch romSizeCode {
	case 0x00: // 32KB
		return 32 * 1024
	case 0x01: // 64KB
		return 64 * 1024
	case 0x02: // 128KB
		return 128 * 1024
	case 0x03: // 256KB
		return 256 * 1024
	case 0x04: // 512KB
		return 512 * 1024
	case 0x05: // 1MB
		return 1024 * 1024
	case 0x06: // 2MB
		return 2048 * 1024
	case 0x07: // 4MB
		return 4096 * 1024
	case 0x08: // 8MB
		return 8192 * 1024
	default:
		return 32 * 1024
	}
}

func getRamSize(ramSizeCode uint8) int {
	switch ramSizeCode {
	case 0x00: // No RAM
		return 0
	case 0x01: // 2KB
		return 2 * 1024
	case 0x02: // 8KB
		return 8 * 1024
	case 0x03: // 32KB
		return 32 * 1024
	case 0x04: // 128KB
		return 128 * 1024
	case 0x05: // 64KB
		return 64 * 1024
	default:
		return 0
	}
}
