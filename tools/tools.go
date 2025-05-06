package tools

// Max returns the maximum of two integers.
//
// Parameters:
//   - a: The first integer.
//   - b: The second integer.
//
// Returns:
//   - int: The maximum of the two integers.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min returns the minimum of two integers.
//
// Parameters:
//   - a: The first integer.
//   - b: The second integer.
//
// Returns:
//   - int: The minimum of the two integers.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CalcBytesPerPixel calculates the number of bytes per pixel based on the bit count.
//
// Parameters:
//   - bitCount: The number of bits per pixel.
//
// Returns:
//   - int: The number of bytes per pixel.
func CalcBytesPerPixel(bitCount int) int {
	if bitCount >= 8 {
		return bitCount / 8
	} else {
		return 1 // For 1-bit and 4-bit BMPs, treat as 1 byte per pixel for row size calculation
	}
}
