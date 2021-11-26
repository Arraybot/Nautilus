package commands

import "testing"

func TestWhiteToHex(t *testing.T) {
	res := rgbToHex(255, 255, 255)
	if res != 0xFFFFFF {
		t.Error("White RGB->Hex failed")
	}
}

func TestBlackToHex(t *testing.T) {
	res := rgbToHex(0, 0, 0)
	if res != 0x000000 {
		t.Error("Black RGB->Hex failed")
	}
}

func TestBlueToHex(t *testing.T) {
	res := rgbToHex(0, 0, 255)
	if res != 0x0000FF {
		t.Error("Blue RGB->Hex failed")
	}
}

func TestRandomToHex(t *testing.T) {
	res := rgbToHex(18, 52, 86)
	if res != 0x123456 {
		t.Error("Color RGB->Hex failed")
	}
}

func TestWhiteToRGB(t *testing.T) {
	r, g, b := hexToRgb(0xFFFFFF)
	if r != 255 && g != 255 && b != 255 {
		t.Error("White Hex->RGB failed")
	}
}

func TestBlackToRGB(t *testing.T) {
	r, g, b := hexToRgb(0x000000)
	if r != 0 && g != 0 && b != 0 {
		t.Error("Black Hex->RGB failed")
	}
}

func TestBlueToRGB(t *testing.T) {
	r, g, b := hexToRgb(0x0000FF)
	if r != 0 && g != 0 && b != 255 {
		t.Error("Blue Hex->RGB failed")
	}
}

func TestRandomToRGB(t *testing.T) {
	r, g, b := hexToRgb(0x123456)
	if r != 18 && g != 52 && b != 86 {
		t.Error("Color Hex->RGB failed")
	}
}

func Test0BytesToMB(t *testing.T) {
	res := bytesToMegabytes(0)
	if res != 0 {
		t.Error("0 bytes should be 0 megabytes")
	}
}

func Test1000BytesToMB(t *testing.T) {
	res := bytesToMegabytes(1000000)
	if res >= 1 {
		t.Error("1000000 bytes should be < 1 megabytes")
	}
}

func Test1024BytesToMB(t *testing.T) {
	res := bytesToMegabytes(1048576)
	if res != 1 {
		t.Error("1048576 bytes should be 1 megabyte")
	}
}

func TestManyBytesToMB(t *testing.T) {
	res := bytesToMegabytes(4194304)
	if res != 4 {
		t.Error("4194304 bytes should be 4 megabytes")
	}
}
