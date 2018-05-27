package ccittfaxdecode

type bitBuffer struct {
	Buffer    uint32
	EmptyBits uint8
	source    []byte
	sourcePos uint
}

func (b *bitBuffer) FlushBits(count uint8) {
	b.Buffer = b.Buffer << count
	b.EmptyBits += count
	b.tryfillBuffer()
}

func (b *bitBuffer) Peak8() (uint8, uint8) {
	return uint8(b.Buffer >> 24), 32 - b.EmptyBits
}

func (b *bitBuffer) Peak16() (uint16, uint8) {
	return uint16(b.Buffer >> 16), 32 - b.EmptyBits
}

func (b *bitBuffer) Peak32() (uint32, uint8) {
	return b.Buffer, 32 - b.EmptyBits
}

func (b *bitBuffer) HasData() bool {
	if b.EmptyBits == 32 && int(b.sourcePos) >= len(b.source) {
		return false
	}
	return true
}

func (b *bitBuffer) Clear() {
	b.Buffer = 0
	b.EmptyBits = 32
	b.sourcePos = 0
}

func (b *bitBuffer) tryfillBuffer() {
	for b.EmptyBits > 7 {
		if b.sourcePos >= uint(len(b.source)) {
			break
		}
		b.AddByte(b.source[b.sourcePos])
		b.sourcePos++
	}
}

func newBitBuffer(source []byte) *bitBuffer {
	buffer := &bitBuffer{
		EmptyBits: 32,
		Buffer:    0,
		source:    source,
		sourcePos: 0,
	}
	buffer.tryfillBuffer()
	return buffer
}

func (b *bitBuffer) AddByte(source byte) {
	padRight := b.EmptyBits - 8
	zeroed := b.Buffer >> (8 + padRight) << (8 + padRight) // switch to hex AND?
	b.Buffer = zeroed | (uint32(source) << padRight)
	b.EmptyBits -= 8
}
