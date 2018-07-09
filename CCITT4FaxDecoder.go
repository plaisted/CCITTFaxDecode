package ccittfaxdecode

import (
	"errors"
	"image"
	"image/color"

	modes "github.com/plaisted/CCITTFaxDecode/ModeCodes"
)

//
type CCITT4FaxDecoder struct {
	ReverseColor    bool
	width           uint
	buffer          *bitBuffer
	modeCodes       []modeCode
	horizontalCodes *horizontalCodes
}

func NewCCITTFaxDecoder(width uint, bytes []byte) *CCITT4FaxDecoder {
	return &CCITT4FaxDecoder{
		width:           width,
		horizontalCodes: newHorizontalCodes(),
		modeCodes:       getModes(),
		buffer:          newBitBuffer(bytes),
	}
}

func (r *CCITT4FaxDecoder) Decode() ([][]uint8, error) {
	lines := make([][]uint8, 0)
	line := make([]uint8, r.width)
	linePos := 0
	curLine := 0
	var a0Color uint8
	a0Color = 255 // start white

	for r.buffer.HasData() {
		if linePos > int(r.width)-1 {
			lines = append(lines, line)
			line = make([]uint8, r.width)
			linePos = 0
			a0Color = 255 // start white
			curLine++
			if endOfBlock(r.buffer.Buffer) {
				break
			}
		}

		// end on trailing zeros padding
		if v, _ := r.buffer.Peak32(); v == 0x00000000 {
			break
		}

		// mode lookup
		mode, err := r.getMode()
		if err != nil {
			return append(lines, line), err
		}
		r.buffer.FlushBits(mode.BitsUsed)

		// act on mode
		switch mode.Type {
		case modes.Pass:
			_, b2 := findBValues(getPreviousLine(lines, curLine, r.width), linePos, a0Color, false)
			for p := linePos; p < b2; p++ {
				line[linePos] = a0Color
				linePos++
			}
			// a0 color should stay the same
		case modes.Extension:
			return lines, errors.New("CCITTFax extensions not supported")
		case modes.Horizontal:
			isWhite := a0Color == 255

			length := []uint16{0, 0}
			color := []uint8{127, 127}
			for i := 0; i < 2; i++ {
				scan := true
				for scan {
					h, err := r.horizontalCodes.FindMatch32(r.buffer.Buffer, isWhite)
					if err != nil {
						return nil, err
					}
					r.buffer.FlushBits(h.BitsUsed)
					length[i] += h.Pixels
					color[i] = uint8(h.CColor)

					if h.Terminating {
						isWhite = !isWhite
						scan = false
					}
				}
			}

			for i := 0; i < 2; i++ {
				for p := 0; p < int(length[i]); p++ {
					if linePos < len(line) {
						line[linePos] = color[i]
					}
					linePos++
				}
			}
			// a0 color should stay the same
		case modes.VerticalZero:
			fallthrough
		case modes.VerticalL1:
			fallthrough
		case modes.VerticalR1:
			fallthrough
		case modes.VerticalL2:
			fallthrough
		case modes.VerticalR2:
			fallthrough
		case modes.VerticalL3:
			fallthrough
		case modes.VerticalR3:
			offset := mode.GetVerticalOffset()
			b1, _ := findBValues(getPreviousLine(lines, curLine, r.width), linePos, a0Color, true)

			for i := linePos; i < b1+offset; i++ {
				if linePos < len(line) {
					line[linePos] = a0Color
				}
				linePos++
			}

			// a0 color changes
			a0Color = reverseColor(a0Color)
		default:
			return lines, errors.New("unknown mode type")
		}

	}

	if r.ReverseColor {
		for i := 0; i < len(lines); i++ {
			for x := 0; x < len(lines[i]); x++ {
				lines[i][x] = reverseColor(lines[i][x])
			}
		}
	}

	return lines, nil
}

func (r *CCITT4FaxDecoder) DecodeToImg() (image.Image, error) {
	var img *image.Gray
	lines, err := r.Decode()
	img = image.NewGray(image.Rect(0, 0, int(r.width), len(lines)))
	for y := 0; y < len(lines); y++ {
		for x := 0; x < int(r.width); x++ {
			if len(lines[y]) > x {
				img.SetGray(x, y, color.Gray{Y: lines[y][x]})
			}
		}
	}
	return img, err
}

func reverseColor(current uint8) uint8 {
	if current == 0 {
		return 255
	} else {
		return 0
	}
}

func endOfBlock(buffer uint32) bool {
	return (buffer & 0xffffff00) == 0x00100100
}

func getPreviousLine(lines [][]byte, currentLine int, width uint) []byte {
	if currentLine == 0 {
		whiteOut := make([]byte, width)
		for i := 0; i < len(whiteOut); i++ {
			whiteOut[i] = 255
		}
		return whiteOut
	} else {
		return lines[currentLine-1]
	}
}

func findBValues(refLine []byte, a0pos int, a0Color uint8, justb1 bool) (b1, b2 int) {
	other := reverseColor(a0Color)
	startPos := a0pos
	if startPos != 0 {
		startPos += 1
	}

	for i := startPos; i < len(refLine); i++ {
		var curColor byte
		var lastColor byte
		if i == 0 {
			curColor = refLine[0]
			lastColor = 255
		} else {
			curColor = refLine[i]
			lastColor = refLine[i-1]
		}

		if b1 != 0 {
			if curColor == a0Color && lastColor == other {
				b2 = i
				return
			}
		}

		if curColor == other && lastColor == a0Color {
			b1 = i
			if b2 != 0 || justb1 {
				return
			}
		}
	}
	if b1 == 0 {
		b1 = len(refLine)
	} else {
		b2 = len(refLine)
	}
	return
}

// TODO: move to faxmodes file and add tree lookup
func (r *CCITT4FaxDecoder) getMode() (modeCode, error) {
	var match modeCode

	// find mode
	b8, _ := r.buffer.Peak8()
	for i := 0; i < len(r.modeCodes); i++ {
		if r.modeCodes[i].Matches(b8) {
			return r.modeCodes[i], nil
		}
	}

	return match, errors.New("bad start")
}
