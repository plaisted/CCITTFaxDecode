package ccittfaxdecode

import (
	"errors"
	"sort"
)

type horizontalCode struct {
	BitsUsed    uint8
	Mask        uint16
	Value       uint16
	CColor      ccolor
	Pixels      uint16
	Terminating bool
}

func (code *horizontalCode) Matches(data uint16) bool {
	if data&code.Mask == code.Value {
		return true
	}
	return false
}

type horizontalCodes struct {
	whiteCodes []horizontalCode
	blackCodes []horizontalCode
}

func newHorizontalCodes() *horizontalCodes {
	return &horizontalCodes{
		whiteCodes: loadWhiteCodes(),
		blackCodes: loadBlackCodes(),
	}

}

func (codes *horizontalCodes) FindMatch32(data uint32, white bool) (horizontalCode, error) {
	return codes.FindMatch(uint16(data>>16), white)
}

func (codes *horizontalCodes) FindMatch(data uint16, white bool) (horizontalCode, error) {
	// this should be modified to use some form of tree search algorithm
	var match horizontalCode
	var lookup []horizontalCode
	if white {
		lookup = codes.whiteCodes
	} else {
		lookup = codes.blackCodes
	}
	for i := 0; i < len(lookup); i++ {
		if lookup[i].Matches(data) {
			return lookup[i], nil
		}
	}

	return match, errors.New("bad horizontal")
}

type ccolor uint8

const (
	black ccolor = 0
	white ccolor = 255
	both  ccolor = 127
)

var whiteTermCodes = [...]uint16{
	0x35, 8, 0x07, 6, 0x07, 4, 0x08, 4, 0x0b, 4, 0x0c, 4, 0x0e, 4, 0x0f, 4,
	0x13, 5, 0x14, 5, 0x07, 5, 0x08, 5, 0x08, 6, 0x03, 6, 0x34, 6, 0x35, 6,
	0x2a, 6, 0x2b, 6, 0x27, 7, 0x0c, 7, 0x08, 7, 0x17, 7, 0x03, 7, 0x04, 7,
	0x28, 7, 0x2b, 7, 0x13, 7, 0x24, 7, 0x18, 7, 0x02, 8, 0x03, 8, 0x1a, 8,
	0x1b, 8, 0x12, 8, 0x13, 8, 0x14, 8, 0x15, 8, 0x16, 8, 0x17, 8, 0x28, 8,
	0x29, 8, 0x2a, 8, 0x2b, 8, 0x2c, 8, 0x2d, 8, 0x04, 8, 0x05, 8, 0x0a, 8,
	0x0b, 8, 0x52, 8, 0x53, 8, 0x54, 8, 0x55, 8, 0x24, 8, 0x25, 8, 0x58, 8,
	0x59, 8, 0x5a, 8, 0x5b, 8, 0x4a, 8, 0x4b, 8, 0x32, 8, 0x33, 8, 0x34, 8,
}
var whiteMakeUpCodes = [...]uint16{
	0x1b, 5, 0x12, 5, 0x17, 6, 0x37, 7, 0x36, 8, 0x37, 8, 0x64, 8, 0x65, 8,
	0x68, 8, 0x67, 8, 0xcc, 9, 0xcd, 9, 0xd2, 9, 0xd3, 9, 0xd4, 9, 0xd5, 9,
	0xd6, 9, 0xd7, 9, 0xd8, 9, 0xd9, 9, 0xda, 9, 0xdb, 9, 0x98, 9, 0x99, 9,
	0x9a, 9, 0x18, 6, 0x9b, 9,
}

var commonMakeUpCodes = [...]uint16{
	0x08, 11, 0x0c, 11, 0x0d, 11, 0x12, 12, 0x13, 12, 0x14, 12, 0x15, 12, 0x16, 12,
	0x17, 12, 0x1c, 12, 0x1d, 12, 0x1e, 12, 0x1f, 12,
}

var blackTermCodes = [...]uint16{
	0x37, 10, 0x02, 3, 0x03, 2, 0x02, 2, 0x03, 3, 0x03, 4, 0x02, 4, 0x03, 5,
	0x05, 6, 0x04, 6, 0x04, 7, 0x05, 7, 0x07, 7, 0x04, 8, 0x07, 8, 0x18, 9,
	0x17, 10, 0x18, 10, 0x08, 10, 0x67, 11, 0x68, 11, 0x6c, 11, 0x37, 11, 0x28, 11,
	0x17, 11, 0x18, 11, 0xca, 12, 0xcb, 12, 0xcc, 12, 0xcd, 12, 0x68, 12, 0x69, 12,
	0x6a, 12, 0x6b, 12, 0xd2, 12, 0xd3, 12, 0xd4, 12, 0xd5, 12, 0xd6, 12, 0xd7, 12,
	0x6c, 12, 0x6d, 12, 0xda, 12, 0xdb, 12, 0x54, 12, 0x55, 12, 0x56, 12, 0x57, 12,
	0x64, 12, 0x65, 12, 0x52, 12, 0x53, 12, 0x24, 12, 0x37, 12, 0x38, 12, 0x27, 12,
	0x28, 12, 0x58, 12, 0x59, 12, 0x2b, 12, 0x2c, 12, 0x5a, 12, 0x66, 12, 0x67, 12,
}

var blackMakeUpCodes = [...]uint16{
	0x0f, 10, 0xc8, 12, 0xc9, 12, 0x5b, 12, 0x33, 12, 0x34, 12, 0x35, 12, 0x6c, 13,
	0x6d, 13, 0x4a, 13, 0x4b, 13, 0x4c, 13, 0x4d, 13, 0x72, 13, 0x73, 13, 0x74, 13,
	0x75, 13, 0x76, 13, 0x77, 13, 0x52, 13, 0x53, 13, 0x54, 13, 0x55, 13, 0x5a, 13,
	0x5b, 13, 0x64, 13, 0x65, 13,
}

//LoadCodes loads and sorts horizontal code objects
func loadWhiteCodes() []horizontalCode {
	totalCodes := (len(blackTermCodes) + len(whiteTermCodes) + len(whiteMakeUpCodes) +
		len(blackMakeUpCodes) + len(commonMakeUpCodes)) / 2
	codes := make([]horizontalCode, totalCodes)
	c := 0

	// white
	for i := 0; i < len(whiteTermCodes)/2; i++ {
		code := &horizontalCode{}
		code.BitsUsed = uint8(whiteTermCodes[i*2+1])
		code.Value = whiteTermCodes[i*2] << (16 - code.BitsUsed)
		code.CColor = white
		code.Mask = 0xffff
		code.Mask = code.Mask << (16 - code.BitsUsed)
		code.Pixels = uint16(i)
		code.Terminating = true
		codes[c] = *code
		c++
	}

	// white make up
	for i := 0; i < len(whiteMakeUpCodes)/2; i++ {
		code := &horizontalCode{}
		code.BitsUsed = uint8(whiteMakeUpCodes[i*2+1])
		code.Value = whiteMakeUpCodes[i*2] << (16 - code.BitsUsed)
		code.CColor = white
		code.Mask = 0xffff
		code.Mask = code.Mask << (16 - code.BitsUsed)
		code.Pixels = uint16((i + 1) * 64)
		codes[c] = *code
		c++
	}

	// common make up
	for i := 0; i < len(commonMakeUpCodes)/2; i++ {
		code := &horizontalCode{}
		code.BitsUsed = uint8(commonMakeUpCodes[i*2+1])
		code.Value = commonMakeUpCodes[i*2] << (16 - code.BitsUsed)
		code.CColor = both
		code.Mask = 0xffff
		code.Mask = code.Mask << (16 - code.BitsUsed)
		code.Pixels = uint16((i+1)*64 + 1728)
		codes[c] = *code
		c++
	}

	sort.Slice(codes, func(i, j int) bool {
		return codes[i].BitsUsed > codes[j].BitsUsed
	})
	return codes
}

func loadBlackCodes() []horizontalCode {
	totalCodes := (len(blackTermCodes) + len(whiteTermCodes) + len(whiteMakeUpCodes) +
		len(blackMakeUpCodes) + len(commonMakeUpCodes)) / 2
	codes := make([]horizontalCode, totalCodes)
	c := 0

	// black
	for i := 0; i < len(blackTermCodes)/2; i++ {
		code := &horizontalCode{}
		code.BitsUsed = uint8(blackTermCodes[i*2+1])
		code.Value = blackTermCodes[i*2] << (16 - code.BitsUsed)
		code.CColor = black
		code.Mask = 0xffff
		code.Mask = code.Mask << (16 - code.BitsUsed)
		code.Pixels = uint16(i)
		code.Terminating = true
		codes[c] = *code
		c++
	}

	// black make up
	for i := 0; i < len(blackMakeUpCodes)/2; i++ {
		code := &horizontalCode{}
		code.BitsUsed = uint8(blackMakeUpCodes[i*2+1])
		code.Value = blackMakeUpCodes[i*2] << (16 - code.BitsUsed)
		code.CColor = black
		code.Mask = 0xffff
		code.Mask = code.Mask << (16 - code.BitsUsed)
		code.Pixels = uint16((i + 1) * 64)
		codes[c] = *code
		c++
	}

	// common make up
	for i := 0; i < len(commonMakeUpCodes)/2; i++ {
		code := &horizontalCode{}
		code.BitsUsed = uint8(commonMakeUpCodes[i*2+1])
		code.Value = commonMakeUpCodes[i*2] << (16 - code.BitsUsed)
		code.CColor = both
		code.Mask = 0xffff
		code.Mask = code.Mask << (16 - code.BitsUsed)
		code.Pixels = uint16((i+1)*64 + 1728)
		codes[c] = *code
		c++
	}

	sort.Slice(codes, func(i, j int) bool {
		return codes[i].BitsUsed > codes[j].BitsUsed
	})
	return codes
}
