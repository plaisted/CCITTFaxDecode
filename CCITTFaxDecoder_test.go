package ccittfaxdecode

import (
	"image/png"

	"io/ioutil"
	"os"
	"testing"

	modeEnums "github.com/plaisted/CCITTFaxDecode/ModeCodes"
)

var flagtests = []struct {
	in  byte
	out modeEnums.Mode
}{
	{0x10, modeEnums.Pass},
	{0x20, modeEnums.Horizontal},
	{0x80, modeEnums.VerticalZero},
	{0x60, modeEnums.VerticalR1},
	{0x0c, modeEnums.VerticalR2},
	{0x06, modeEnums.VerticalR3},
	{0x40, modeEnums.VerticalL1},
	{0x08, modeEnums.VerticalL2},
	{0x04, modeEnums.VerticalL3},
	{0x02, modeEnums.Extension},
}

func TestFlagParser(t *testing.T) {
	for _, tt := range flagtests {
		t.Run(string(tt.out), func(t *testing.T) {
			bytes := []byte{tt.in, 0x00}
			decoder := NewCCITTFaxDecoder(80, bytes)
			code, _ := decoder.getMode()
			if code.Type != tt.out {
				t.Errorf("got %d, want %d", code.Type, tt.out)
			} else {
				t.Logf("ModeCode %d correct", code.Type)
			}
		})
	}
}

type imageTest struct {
	sourceBinary string
	baseLineBmp  string
	reverse      bool
}

var imageTests = []imageTest{
	imageTest{".\\testfiles\\18x18.bin", ".\\testfiles\\18x18.png", false},
	imageTest{".\\testfiles\\80x80reversed.bin", ".\\testfiles\\80x80reversed.png", true},
	imageTest{".\\testfiles\\b0.bin", ".\\testfiles\\b0_baseline.png", false},
	imageTest{".\\testfiles\\b1.bin", ".\\testfiles\\b1_baseline.png", false},
	imageTest{".\\testfiles\\b2.bin", ".\\testfiles\\b2_baseline.png", false},
	imageTest{".\\testfiles\\b4.bin", ".\\testfiles\\b4_baseline.png", false},
	imageTest{".\\testfiles\\b5.bin", ".\\testfiles\\b5_baseline.png", false},
	imageTest{".\\testfiles\\b6.bin", ".\\testfiles\\b6_baseline.png", false},
	imageTest{".\\testfiles\\b8.bin", ".\\testfiles\\b8_baseline.png", false},
	imageTest{".\\testfiles\\b9.bin", ".\\testfiles\\b9_baseline.png", false},
	imageTest{".\\testfiles\\CCITT.498.bin", ".\\testfiles\\CCITT.498.png", true},
	// {".\\testfiles\\80x80reversed.bin", ".\\testfiles\\80x80reversed_bad.png", 80, true},
}

func BenchmarkDecode(b *testing.B) {
	bytes, err := ioutil.ReadFile(".\\testfiles\\b0.bin")
	if err != nil {
		b.Errorf("Error: %v\n", err)
		return
	}

	for n := 0; n < b.N; n++ {
		decoder := NewCCITTFaxDecoder(uint(18), bytes)
		_, err := decoder.DecodeToImg()
		if err != nil {
			b.Errorf("Error: %v\n", err)
			return
		}
	}
}

func TestImages(t *testing.T) {
	for _, tt := range imageTests {
		t.Run(string(tt.sourceBinary), func(t *testing.T) {
			runImage(t, tt)
		})
	}
}

func TestBad(t *testing.T) {
	runImage(t, imageTests[10])
}

func runImage(t *testing.T, tt imageTest) {
	baseline, err := os.Open(tt.baseLineBmp)
	if err != nil {
		t.Errorf("Error: %v\n", err)
		return
	}
	defer baseline.Close()

	bImg, err := png.Decode(baseline)
	if err != nil {
		t.Errorf("Error: %v\n", err)
		return
	}

	bytes, err := ioutil.ReadFile(tt.sourceBinary)
	if err != nil {
		t.Errorf("Error: %v\n", err)
		return
	}

	decoder := NewCCITTFaxDecoder(uint(bImg.Bounds().Dx()), bytes)
	decoder.ReverseColor = tt.reverse
	img, err := decoder.DecodeToImg()
	if err != nil {
		f, err2 := os.Create(tt.sourceBinary + ".err.png")
		if err2 == nil {
			png.Encode(f, img)
		}
		f.Close()
		t.Errorf("Error: %v\n", err)
		return
	}

	if img.Bounds().Dx() != bImg.Bounds().Dx() || bImg.Bounds().Dy() != img.Bounds().Dy() {
		t.Error("Bad image dimensions")
		return
	}

	errored := false
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			r, b, g, a := img.At(x, y).RGBA()
			rb, bb, gb, ab := bImg.At(x, y).RGBA()
			if r != rb || b != bb || g != gb || a != ab {
				t.Errorf("Bad pixel: %d %d", x, y)
				errored = true
			}
		}
	}
	if !errored {
		t.Logf("E2E success: %s", tt.baseLineBmp)
	}

}
