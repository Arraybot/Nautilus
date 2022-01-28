package commands

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

// Convert bytes to megabytes.
func bytesToMegabytes(b uint64) uint64 {
	return b / 1024 / 1024
}

// Paginate a certain length of items into pages of 10 and return slice from (incl) and to (excl) indices and total pages.
func paginate(length, page int) (int, int, int) {
	if length == 0 {
		return 0, 0, 0
	}
	pageFloor := length / 10
	pageModulo := length % 10
	var pageTotal int
	if pageModulo == 0 {
		pageTotal = pageFloor
	} else {
		pageTotal = pageFloor + 1
	}
	if page < 1 {
		page = 1
	}
	if page > pageTotal {
		page = pageTotal
	}
	lower := (page - 1) * 10
	var upper int
	if page != pageTotal || pageModulo == 0 {
		upper = lower + 10
	} else {
		upper = lower + pageModulo
	}
	return lower, upper, pageTotal
}
