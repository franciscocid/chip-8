package main

import (
	"fmt"
	"strings"
)

func Disassemble(program []uint16) string {
	var res strings.Builder

	for _, code := range program {
		decoded := decode(code)
		res.WriteString(fmt.Sprintf("0x%04x\t", decoded))
		res.WriteString(generate(decoded))
		res.WriteString("\n")
	}

	return res.String()
}

func decode(code uint16) uint16 {
	return (code & 0xFF00 >> 8) | (code & 0x00FF << 8)
}

func generate(code uint16) string {
	switch code {
	case 0x00E0:
		return "CLS"

	case 0x00EE:
		return "RET"
	}

	x := xRegister(code)
	y := yRegister(code)
	addr := address(code)
	value := value(code)

	switch code & 0xF000 {
	case 0x0000:
		return fmt.Sprintf("SYS\t0x%03x", addr)

	case 0x1000:
		return fmt.Sprintf("JMP\t0x%03x", addr)

	case 0x2000:
		return fmt.Sprintf("CALL\t0x%03x", addr)

	case 0x3000:
		return fmt.Sprintf("SE\tV%d, 0x%03x", x, addr)

	case 0x4000:
		return fmt.Sprintf("SNE\tV%d, 0x%02x", x, value)

	case 0x5000:
		return fmt.Sprintf("SNE\tV%d, V%d", x, y)

	case 0x6000:
		return fmt.Sprintf("LD\tV%d, 0x%02x", x, value)

	case 0x7000:
		return fmt.Sprintf("ADD\tV%d, 0x%02x", x, value)

	case 0x8000:
		suffix := code & 0x000F

		switch suffix {
		case 0x0:
			return fmt.Sprintf("LD\tV%d, V%d", x, y)

		case 0x1:
			return fmt.Sprintf("OR\tV%d, V%d", x, y)

		case 0x2:
			return fmt.Sprintf("AND\tV%d, V%d", x, y)

		case 0x3:
			return fmt.Sprintf("XOR\tV%d, V%d", x, y)

		case 0x4:
			return fmt.Sprintf("ADD\tV%d, V%d", x, y)

		case 0x5:
			return fmt.Sprintf("SUB\tV%d, V%d", x, y)

		case 0x6:
			return fmt.Sprintf("SHR\tV%d, V%d", x, y)

		case 0x7:
			return fmt.Sprintf("SUBN\tV%d, V%d", x, y)

		case 0xE:
			return fmt.Sprintf("SHL\tV%d, V%d", x, y)
		}

	case 0x9000:
		return fmt.Sprintf("SNE\tV%d, V%d", x, y)

	case 0xA000:
		return fmt.Sprintf("LD\tI, 0x%03x", addr)

	case 0xB000:
		return fmt.Sprintf("JMP\tV0, 0x%03x", addr)

	case 0xC000:
		return fmt.Sprintf("RND\tV%d, 0x%02x", x, value)

	case 0xD000:
		nib := nibble(code)
		return fmt.Sprintf("DRW\tV%d, V%d, 0x%x", x, y, nib)

	case 0xE000:
		suffix := code & 0x00FF
		switch suffix {
		case 0x9E:
			return fmt.Sprintf("SKP\tV%d", x)
		case 0xA1:
			return fmt.Sprintf("SKNP\tV%d", x)
		}

	case 0xF000:
		suffix := code & 0x00FF
		switch suffix {
		case 0x07:
			return fmt.Sprintf("LD\tV%d, DT", x)

		case 0x0A:
			return fmt.Sprintf("LD\tV%d, KEY", x)

		case 0x15:
			return fmt.Sprintf("LD\tDT, V%d", x)

		case 0x18:
			return fmt.Sprintf("LD\tST, V%d", x)

		case 0x1E:
			return fmt.Sprintf("ADD\tI, V%d", x)
			
		case 0x29:
			return fmt.Sprintf("LD\tF, V%d", x)

		case 0x33:
			return fmt.Sprintf("LD\tB, V%d", x)
			
		case 0x55:
			return fmt.Sprintf("LD\t[I], V%d", x)
			
		case 0x65:
			return fmt.Sprintf("LD\tV%d, [I]", x)
		}
		
	}

	return fmt.Sprintf("DAT\t0x%04x", code)
}

func xRegister(code uint16) uint16 {
	return code & 0x0F00 >> 8
}

func yRegister(code uint16) uint16 {
	return code & 0x00F0 >> 8
}

func value(code uint16) uint16 {
	return code & 0x00FF
}

func address(code uint16) uint16 {
	return code & 0x0FFF
}

func nibble(code uint16) uint16 {
	return code & 0x000F
}

func main() {

	program := []uint16{
		0x1a12, 0x2e32, 0x3030, 0x4320, 0x202e, 0x6745, 0x6265, 0x7265,
		0x2067, 0x3831, 0x382f, 0x272d, 0x3139, 0x0380, 0x1381, 0xc8a8,
		0x55f1, 0x0560, 0xcca8, 0x55f0, 0x7387, 0x6386, 0x7227, 0xe000,
		0x9427, 0x406e, 0xe287, 0x276e, 0xe187, 0x1a68, 0x0c69, 0x386a,
		0x006b, 0x026c, 0x1a6d, 0x5027, 0xeda8, 0xb4da, 0xd4dc, 0xd023,
		0x003e, 0x7c12, 0xcca8, 0x65f0, 0x0085, 0xffc4, 0x5284, 0xf624,
		0xffc4, 0x5284, 0x1e26, 0x0160, 0xa1e0, 0xd627, 0xf736, 0x4e12,
		0x608e, 0x7a28, 0x646e, 0x7a28, 0xd627, 0x2a12, 0x07f0, 0x0040,
		0x1013, 0x8080, 0x0680, 0xa081, 0x0681, 0x1580, 0x0040, 0x9a12,
		0x0140, 0x9a12, 0xff40, 0x9a12, 0xc812, 0x9080, 0x0680, 0xb081,
		0x0681, 0x1580, 0x0040, 0xb212, 0x0140, 0xb212, 0xff40, 0xb212,
		0xc812, 0xeda8, 0xb4da, 0x386a, 0x006b, 0xb4da, 0xf36e, 0xe287,
		0x046e, 0xe187, 0x326e, 0x7a28, 0x8080, 0x0680, 0xc081, 0x0681,
		0x1580, 0x0040, 0xe012, 0x0140, 0xe012, 0xff40, 0xe012, 0x5412,
		0x9080, 0x0680, 0xd081, 0x0681, 0x1580, 0x0040, 0xf812, 0x0140,
		0xf812, 0xff40, 0xf812, 0x5412, 0xeda8, 0xd4dc, 0x026c, 0x1a6d,
		0xd4dc, 0xcf6e, 0xe287, 0x206e, 0xe187, 0x196e, 0x7a28, 0x5412,
		0x3f60, 0xa828, 0x5027, 0xeda8, 0xb4da, 0xd4dc, 0x406e, 0xe387,
		0x7080, 0xe280, 0x0030, 0x3212, 0x608e, 0x7a28, 0x8a28, 0xe000,
		0x1166, 0x0a67, 0xcaa8, 0xe627, 0x1166, 0x1067, 0xc8a8, 0xe627,
		0x0064, 0x0865, 0x0066, 0x0f67, 0x19ab, 0x69d4, 0x22ab, 0x69d5,
		0x0360, 0xa828, 0x003e, 0xc613, 0x19ab, 0x69d4, 0x22ab, 0x69d5,
		0x0274, 0x0275, 0x3034, 0x4813, 0x19ab, 0x69d4, 0x22ab, 0x69d5,
		0x0360, 0xa828, 0x003e, 0xc613, 0x19ab, 0x69d4, 0x22ab, 0x69d5,
		0x0276, 0x1636, 0x6813, 0x19ab, 0x69d4, 0x22ab, 0x69d5, 0x0360,
		0xa828, 0x003e, 0xc613, 0x19ab, 0x69d4, 0x22ab, 0x69d5, 0xfe74,
		0xfe75, 0x0034, 0x8613, 0x19ab, 0x69d4, 0x22ab, 0x69d5, 0x0360,
		0xa828, 0x003e, 0xc613, 0x19ab, 0x69d4, 0x22ab, 0x69d5, 0xfe76,
		0x0036, 0xa613, 0x4813, 0x22ab, 0x69d5, 0x2bab, 0x69d5, 0x1a12,
		0x7083, 0x036e, 0xe283, 0x8084, 0x9085, 0x066e, 0xa1ee, 0x3214,
		0x036e, 0xa1ee, 0x4a14, 0x086e, 0xa1ee, 0x6214, 0x076e, 0xa1ee,
		0x7a14, 0x0343, 0x0275, 0x0043, 0xfe75, 0x0243, 0x0274, 0x0143,
		0xfe74, 0x4080, 0x5081, 0xba27, 0x0082, 0x086e, 0xe280, 0x0030,
		0x9214, 0x076e, 0x2080, 0xe282, 0x0542, 0x9a14, 0x0642, 0xb214,
		0x0742, 0xec14, 0x5027, 0xfc6e, 0xe287, 0x3187, 0x4088, 0x5089,
		0x5017, 0x4080, 0x5081, 0x0271, 0xba27, 0x0082, 0x086e, 0xe280,
		0x0030, 0xf213, 0x0363, 0x0275, 0x0e14, 0x4080, 0x5081, 0xfe71,
		0xba27, 0x0082, 0x086e, 0xe280, 0x0030, 0xf213, 0x0063, 0xfe75,
		0x0e14, 0x4080, 0x5081, 0x0270, 0xba27, 0x0082, 0x086e, 0xe280,
		0x0030, 0xf213, 0x0263, 0x0274, 0x0e14, 0x4080, 0x5081, 0xfe70,
		0xba27, 0x0082, 0x086e, 0xe280, 0x0030, 0xf213, 0x0163, 0xfe74,
		0x0e14, 0x5027, 0x94d8, 0xf08e, 0xee00, 0xf06e, 0xe280, 0x3180,
		0x55f0, 0xf1a8, 0x54d4, 0x0176, 0x0561, 0x07f0, 0x0040, 0x18f1,
		0x2414, 0xf06e, 0xe280, 0x3180, 0x55f0, 0xf5a8, 0x54d4, 0x0476,
		0xa080, 0xb081, 0xba27, 0xf06e, 0xe280, 0x0030, 0xd214, 0x0c6e,
		0xe387, 0xc080, 0xd081, 0xba27, 0xf06e, 0xe280, 0x0030, 0xe414,
		0x306e, 0xe387, 0xff60, 0x18f0, 0x15f0, 0x2414, 0x0143, 0x3a64,
		0x0243, 0x0064, 0x2414, 0x7082, 0x7083, 0x0c6e, 0xe282, 0xa080,
		0xb081, 0xba27, 0xeda8, 0xf06e, 0xe280, 0x0030, 0x2415, 0xb4da,
		0x0c42, 0x027b, 0x0042, 0xfe7b, 0x0842, 0x027a, 0x0442, 0xfe7a,
		0xb4da, 0xee00, 0x806e, 0x07f1, 0x0031, 0xd415, 0x0034, 0xd415,
		0x0081, 0x0e83, 0x003f, 0x5615, 0x9083, 0xb583, 0x004f, 0x8c15,
		0x0033, 0x7415, 0xe387, 0x8083, 0xa583, 0x004f, 0xbc15, 0x0033,
		0xa415, 0xe387, 0xd415, 0x8083, 0xa583, 0x004f, 0xbc15, 0x0033,
		0xa415, 0xe387, 0x9083, 0xb583, 0x004f, 0x8c15, 0x0033, 0x7415,
		0xe387, 0xd415, 0x4063, 0x3281, 0x0041, 0xd415, 0xb4da, 0x027b,
		0xb4da, 0xf36e, 0xe287, 0x0c62, 0x2187, 0xee00, 0x1063, 0x3281,
		0x0041, 0xd415, 0xb4da, 0xfe7b, 0xb4da, 0xf36e, 0xe287, 0x0062,
		0x2187, 0xee00, 0x2063, 0x3281, 0x0041, 0xd415, 0xb4da, 0x027a,
		0xb4da, 0xf36e, 0xe287, 0x0862, 0x2187, 0xee00, 0x8063, 0x3281,
		0x0041, 0xd415, 0xb4da, 0xfe7a, 0xb4da, 0xf36e, 0xe287, 0x0462,
		0x2187, 0xee00, 0xf0c1, 0x1280, 0x0030, 0xe415, 0x0c6e, 0xe387,
		0xe382, 0x0e15, 0xb4da, 0x0e80, 0x004f, 0xf215, 0x0462, 0xfe7a,
		0x1416, 0x0e80, 0x004f, 0xfe15, 0x0c62, 0x027b, 0x1416, 0x0e80,
		0x004f, 0x0a16, 0x0862, 0x027a, 0x1416, 0x0e80, 0x004f, 0xdc15,
		0x0062, 0xfe7b, 0xb4da, 0xf36e, 0xe287, 0x2187, 0xee00, 0x7082,
		0x7083, 0x306e, 0xe282, 0xc080, 0xd081, 0xba27, 0xeda8, 0xf06e,
		0xe280, 0x0030, 0x4c16, 0xd4dc, 0x3042, 0x027d, 0x0042, 0xfe7d,
		0x2042, 0x027c, 0x1042, 0xfe7c, 0xd4dc, 0xee00, 0x806e, 0x07f1,
		0x0031, 0x0417, 0x0034, 0x0417, 0x0081, 0x0e83, 0x004f, 0x7e16,
		0x9083, 0xd583, 0x004f, 0xb616, 0x0033, 0x9c16, 0xe387, 0x8083,
		0xc583, 0x004f, 0xea16, 0x0033, 0xd016, 0xe387, 0x0417, 0x8083,
		0xc583, 0x004f, 0xea16, 0x0033, 0xd016, 0xe387, 0x9083, 0xd583,
		0x004f, 0xb616, 0x0033, 0x9c16, 0xe387, 0x0417, 0x4063, 0x3281,
		0x0041, 0x0417, 0xd4dc, 0x027d, 0xd4dc, 0xe387, 0xcf6e, 0xe287,
		0x3062, 0x2187, 0xee00, 0x1063, 0x3281, 0x0041, 0x0417, 0xd4dc,
		0xfe7d, 0xd4dc, 0xe387, 0xcf6e, 0xe287, 0x0062, 0x2187, 0xee00,
		0x2063, 0x3281, 0x0041, 0x0417, 0xd4dc, 0x027c, 0xd4dc, 0xe387,
		0xcf6e, 0xe287, 0x2062, 0x2187, 0xee00, 0x8063, 0x3281, 0x0041,
		0x0417, 0xd4dc, 0xfe7c, 0xd4dc, 0xe387, 0xcf6e, 0xe287, 0x1062,
		0x2187, 0xee00, 0xf0c1, 0x1280, 0x0030, 0x1617, 0xe387, 0x306e,
		0xe387, 0xe382, 0x3616, 0xd4dc, 0x0e80, 0x004f, 0x2417, 0x9062,
		0xfe7c, 0x4617, 0x0e80, 0x004f, 0x3017, 0x3062, 0x027d, 0x4617,
		0x0e80, 0x004f, 0x3c17, 0xa062, 0x027c, 0x4617, 0x0e80, 0x004f,
		0x0c17, 0x0062, 0xfe7d, 0xd4dc, 0x4f6e, 0xe287, 0x2187, 0xee00,
		0x7080, 0x036e, 0xe280, 0x0e80, 0x8081, 0x9481, 0x026e, 0xe281,
		0x0041, 0x0170, 0x0e80, 0x0e80, 0xcda8, 0x1ef0, 0x94d8, 0xf08e,
		0xee00, 0x006e, 0x19a9, 0x1efe, 0x1efe, 0x1efe, 0x1efe, 0x65f3,
		0x34ab, 0x1efe, 0x1efe, 0x1efe, 0x1efe, 0x55f3, 0x017e, 0x803e,
		0x7417, 0xee00, 0x2382, 0x3383, 0x0f6e, 0x2080, 0x3081, 0xbe27,
		0xe280, 0x0e80, 0xf9a8, 0x1ef0, 0x32d2, 0x0272, 0x4032, 0x9a17,
		0x2382, 0x0273, 0x2043, 0xee00, 0x9a17, 0x0270, 0x0271, 0x0680,
		0x0681, 0x0e81, 0x0e81, 0x0e81, 0x0e81, 0x34ab, 0x1ef1, 0x1ef1,
		0x1ef0, 0x65f0, 0xee00, 0xcca8, 0x65f0, 0x0680, 0x55f0, 0x0160,
		0xa1e0, 0xe017, 0xee00, 0x65f1, 0x016e, 0x4384, 0x0082, 0x1083,
		0x1065, 0x5583, 0x004f, 0xe582, 0x004f, 0x0c18, 0x2765, 0x5582,
		0x004f, 0x0c18, 0x2080, 0x3081, 0xe484, 0xf017, 0x29f4, 0x75d6,
		0x0676, 0x4384, 0x0082, 0x1083, 0xe865, 0x5583, 0x004f, 0xe582,
		0x004f, 0x3418, 0x0365, 0x5582, 0x004f, 0x3418, 0x2080, 0x3081,
		0xe484, 0x1818, 0x29f4, 0x75d6, 0x0676, 0x4384, 0x0082, 0x1083,
		0x6465, 0x5583, 0x004f, 0xe582, 0x004f, 0x5418, 0x2080, 0x3081,
		0xe484, 0x4018, 0x29f4, 0x75d6, 0x0676, 0x4384, 0x0082, 0x1083,
		0x0a65, 0x5583, 0x004f, 0x6e18, 0x3081, 0xe484, 0x6018, 0x29f4,
		0x75d6, 0x0676, 0x29f1, 0x75d6, 0xee00, 0xc8a8, 0x65f1, 0xe481,
		0x003f, 0x0170, 0xc8a8, 0x55f1, 0xee00, 0xc8a8, 0x65f3, 0x008e,
		0x258e, 0x004f, 0xee00, 0x003e, 0xa218, 0x108e, 0x358e, 0x004f,
		0xee00, 0xcaa8, 0x55f1, 0xee00, 0xe38e, 0x0f62, 0xff63, 0x1061,
		0xa1e2, 0xc418, 0x3481, 0x0031, 0xb018, 0x1061, 0x3480, 0x0030,
		0xb018, 0xee00, 0x016e, 0xee00, 0x0000, 0x0000, 0x0005, 0x7050,
		0x0020, 0x7050, 0x0020, 0x3060, 0x0060, 0x3060, 0x0060, 0x6030,
		0x0030, 0x6030, 0x0030, 0x7020, 0x0050, 0x7020, 0x0050, 0x7020,
		0x0070, 0x2000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000,
		0x0000, 0x8000, 0x0000, 0x0000, 0xc000, 0x0000, 0x8000, 0x0080,
		0xc000, 0x8080, 0xc080, 0x8000, 0x0c00, 0x0808, 0x0808, 0x0808,
		0x0808, 0x0808, 0x0808, 0x0808, 0x0c0d, 0x0808, 0x0808, 0x0808,
		0x0808, 0x0808, 0x0808, 0x0808, 0x0a0d, 0x0565, 0x0505, 0xe505,
		0x0505, 0x05e5, 0x0505, 0xc505, 0x0a0a, 0x0565, 0x0505, 0xe505,
		0x0505, 0x05e5, 0x0505, 0xc505, 0x0a0a, 0x0c05, 0x0808, 0x050f,
		0x0d0c, 0x0805, 0x0808, 0x050d, 0x0f0e, 0x0c05, 0x0808, 0x050f,
		0x0d0c, 0x0805, 0x0808, 0x050d, 0x0a0a, 0x0a05, 0x0665, 0x9505,
		0x0a0a, 0x0535, 0xc505, 0x350a, 0x0505, 0x0a95, 0x0565, 0x9505,
		0x0a0a, 0x0535, 0xc506, 0x050a, 0x0a0a, 0x0f05, 0x0805, 0x0808,
		0x0808, 0x080c, 0x050f, 0x0808, 0x0808, 0x0f08, 0x0805, 0x0c08,
		0x0808, 0x0808, 0x050f, 0x050f, 0x0a0a, 0x0575, 0x05b5, 0x0505,
		0xc505, 0x650a, 0xb505, 0xe505, 0x0505, 0x05e5, 0x05b5, 0x0ac5,
		0x0565, 0x0505, 0xb505, 0xd505, 0x0a0a, 0x0c05, 0x0808, 0x0808,
		0x050d, 0x050f, 0x080c, 0x050f, 0x0f08, 0x0805, 0x0d08, 0x0f05,
		0x0c05, 0x0808, 0x0808, 0x050d, 0x0f0a, 0x0f05, 0x0565, 0xc505,
		0x350a, 0x95e5, 0x650a, 0xb005, 0x0505, 0x05b5, 0x0ac5, 0xe535,
		0x0a95, 0x0565, 0xc505, 0x050f, 0x070f, 0x0574, 0x08d5, 0x050f,
		0x0f0e, 0x0805, 0x050f, 0x080c, 0x0808, 0x0d08, 0x0805, 0x050f,
		0x0f08, 0x0805, 0x750f, 0xd405, 0x0a07, 0x0a05, 0x0535, 0xf505,
		0x0505, 0x05b5, 0xd505, 0x0808, 0x0c0d, 0x0f08, 0x0575, 0xb505,
		0x0505, 0x05f5, 0x9505, 0x050a, 0x0a0a, 0x0805, 0x0808, 0x050d,
		0x080c, 0x0808, 0x350d, 0xc505, 0x0a0a, 0x0565, 0x0c95, 0x0808,
		0x0d08, 0x0c05, 0x0808, 0x050f, 0x0a0a, 0x0575, 0xc506, 0x050a,
		0x0808, 0x0808, 0x0808, 0x050f, 0x0f08, 0x0805, 0x0808, 0x0808,
		0x0f08, 0x0a05, 0x0665, 0xd505, 0x0a0a, 0x0c05, 0x050d, 0x350a,
		0x0505, 0x0505, 0x05e5, 0xf505, 0x0505, 0x05f5, 0xe505, 0x0505,
		0x0505, 0x0a95, 0x0c05, 0x050d, 0x0a0a, 0x0805, 0x050f, 0x0808,
		0x0808, 0x0f08, 0x0c05, 0x050d, 0x0f08, 0x0c05, 0x050d, 0x0808,
		0x0808, 0x0f08, 0x0805, 0x050f, 0x0a0a, 0x0535, 0xb505, 0x0505,
		0x0505, 0x0505, 0x0a95, 0x350a, 0x0505, 0x0a95, 0x350a, 0x0505,
		0x0505, 0x0505, 0x05b5, 0x9505, 0x080a, 0x0808, 0x0808, 0x0808,
		0x0808, 0x0808, 0x0f08, 0x0808, 0x0808, 0x0f08, 0x0808, 0x0808,
		0x0808, 0x0808, 0x0808, 0x0808, 0x3c0f, 0x9942, 0x4299, 0x013c,
		0x0f10, 0x8478, 0x3232, 0x7884, 0x1000, 0x78e0, 0xfefc, 0x84fe,
		0x0078, 0xe010,
	}

	fmt.Println("OPCODE\tSYNTAX\t")
	fmt.Println(Disassemble(program))
}
