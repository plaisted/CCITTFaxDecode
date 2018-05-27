package ccittfaxdecode

import (
	modeEnums "github.com/plaisted/CCITTFaxDecode/ModeCodes"
)

type modeCode struct {
	BitsUsed uint8
	Mask     uint8
	Value    uint8
	Type     modeEnums.Mode
}

var modeCodes = [...]uint8{
	0x1, 4, 1,
	0x1, 3, 2,
	0x1, 1, 3, // 1
	0x03, 3, 4, // 011
	0x03, 6, 5, // 0000 11
	0x03, 7, 6, // 0000 011
	0x2, 3, 7, // 010
	0x02, 6, 8, // 0000 10
	0x02, 7, 9, // 0000 010
	0x01, 7, 10, // 0000 010
}

func (code *modeCode) GetVerticalOffset() int {
	switch code.Type {
	case modeEnums.VerticalZero:
		return 0
	case modeEnums.VerticalL1:
		return -1
	case modeEnums.VerticalR1:
		return 1
	case modeEnums.VerticalL2:
		return -2
	case modeEnums.VerticalR2:
		return 2
	case modeEnums.VerticalL3:
		return -3
	case modeEnums.VerticalR3:
		return 3
	default:
		return 0
	}
}

func (code *modeCode) Matches(data uint8) bool {
	if data&code.Mask == code.Value {
		return true
	}
	return false
}

func getModes() []modeCode {
	modes := make([]modeCode, len(modeCodes)/3)

	for i := 0; i < len(modeCodes)/3; i++ {
		code := &modeCode{}
		code.BitsUsed = uint8(modeCodes[i*3+1])
		code.Value = modeCodes[i*3] << (8 - code.BitsUsed)
		code.Mask = 0xff
		code.Mask = code.Mask << (8 - code.BitsUsed)
		code.Type = modeEnums.Mode(modeCodes[i*3+2])
		modes[i] = *code
	}
	return modes
}
