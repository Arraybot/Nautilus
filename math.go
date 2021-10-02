package main

// Convert RGB values to integer.
func rgbToHex(r, g, b uint32) uint32 {
	// Red is MS 8 bits.
	val := (r << 16)
	// Green is middle 8 bits.
	val |= (g << 8)
	// Blue is LS 8 bits.
	val |= b
	return val
}

// Convert integer to RGB values.
func hexToRgb(hex uint64) (uint8, uint8, uint8) {
	// Red is MS 8 bits.
	r := uint8(hex >> 16)
	// Green is middle 8 bits
	g := uint8((hex >> 8) & 0xFF)
	// Blue is LS 8 bits.
	b := uint8(hex & 0xFF)
	return r, g, b
}
