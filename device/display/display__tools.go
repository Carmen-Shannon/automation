package display

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

// FindSubBMP searches for a smaller BMP within a larger BMP using MSE for fuzzy matching.
// Parameters:
//   - largeBMP: The larger BMP image.
//   - smallBMP: The smaller BMP image (template).
//   - mseThreshold: The maximum allowable MSE for a match.
//
// Returns:
//   - (x, y): The top-left coordinates of the match in the larger BMP.
//   - error: An error if no match is found.
func FindSubBMP(largeBMP, smallBMP BMP, mseThreshold float64) (int, int, error) {
	// Validate dimensions
	if smallBMP.Width > largeBMP.Width || smallBMP.Height > largeBMP.Height {
		return 0, 0, fmt.Errorf("small BMP dimensions exceed large BMP dimensions")
	}

	// Normalize BMPs to top-down format
	largeData := normalizeBMPData(largeBMP)
	smallData := normalizeBMPData(smallBMP)

	// Calculate bytes per pixel for both BMPs
	largeBytesPerPixel := calcBytesPerPixel(int(largeBMP.InfoHeader.BiBitCount))
	smallBytesPerPixel := calcBytesPerPixel(int(smallBMP.InfoHeader.BiBitCount))

	// Calculate row sizes (BMP rows are padded to 4-byte boundaries)
	largeRowSize := ((largeBMP.Width*largeBytesPerPixel + 3) / 4) * 4
	smallRowSize := ((smallBMP.Width*smallBytesPerPixel + 3) / 4) * 4

	// Sliding window search
	for y := 0; y <= largeBMP.Height-smallBMP.Height; y++ {
		for x := 0; x <= largeBMP.Width-smallBMP.Width; x++ {
			// Calculate MSE for the current window
			mse := calculateMSE(largeData, smallData, x, y, largeRowSize, smallRowSize, largeBytesPerPixel, smallBytesPerPixel, smallBMP.Width, smallBMP.Height)
			if mse <= mseThreshold {
				return x, y, nil // Match found
			}
		}
	}

	// No match found
	return 0, 0, fmt.Errorf("no match found")
}

// calculateMSE calculates the Mean Squared Error (MSE) between the current window in the larger BMP and the smaller BMP.
// Parameters:
//   - largeData: The pixel data of the larger BMP.
//   - smallData: The pixel data of the smaller BMP.
//   - startX, startY: The top-left coordinates of the current window in the larger BMP.
//   - largeRowSize, smallRowSize: The row sizes of the larger and smaller BMPs.
//   - largeBytesPerPixel, smallBytesPerPixel: The bytes per pixel for the larger and smaller BMPs.
//   - smallWidth, smallHeight: The dimensions of the smaller BMP.
//
// Returns:
//   - mse: The calculated Mean Squared Error.
func calculateMSE(largeData, smallData []byte, startX, startY, largeRowSize, smallRowSize, largeBytesPerPixel, smallBytesPerPixel, smallWidth, smallHeight int) float64 {
	var totalError float64
	var pixelCount int

	for row := 0; row < smallHeight; row++ {
		// Calculate the starting index for the current row in both BMPs
		largeRowStart := (startY+row)*largeRowSize + startX*largeBytesPerPixel
		smallRowStart := row * smallRowSize

		for col := 0; col < smallWidth; col++ {
			// Calculate the starting index for the current pixel in both BMPs
			largePixelStart := largeRowStart + col*largeBytesPerPixel
			smallPixelStart := smallRowStart + col*smallBytesPerPixel

			// Compare pixel values (assume RGB format for simplicity)
			for i := 0; i < 3; i++ { // Compare R, G, B channels
				largeValue := float64(largeData[largePixelStart+i])
				smallValue := float64(smallData[smallPixelStart+i])
				totalError += (largeValue - smallValue) * (largeValue - smallValue)
			}

			pixelCount++
		}
	}

	// Calculate the mean squared error
	return totalError / float64(pixelCount*3) // Multiply by 3 for RGB channels
}

// normalizeBMPData ensures that the BMP data is in top-down format.
// If the BMP is bottom-up (BiHeight > 0), it flips the rows.
func normalizeBMPData(bmp BMP) []byte {
	// If the BMP is already top-down (negative height), return the data as-is
	if bmp.InfoHeader.BiHeight < 0 {
		return bmp.Data
	}

	// Otherwise, flip the rows to make it top-down
	bytesPerPixel := calcBytesPerPixel(int(bmp.InfoHeader.BiBitCount))
	rowSize := ((bmp.Width*bytesPerPixel + 3) / 4) * 4
	height := int(bmp.InfoHeader.BiHeight)

	normalizedData := make([]byte, len(bmp.Data))
	for row := 0; row < height; row++ {
		srcOffset := (height - 1 - row) * rowSize
		dstOffset := row * rowSize
		copy(normalizedData[dstOffset:dstOffset+rowSize], bmp.Data[srcOffset:srcOffset+rowSize])
	}

	return normalizedData
}

// LoadBmp parses BMP data from a byte slice and extracts the raw pixel data, width, and height.
//
// Parameters:
//   - data: A byte slice containing the BMP file data.
//
// Returns:
//   - *BMP: A pointer to a BMP struct containing the raw pixel data, width, and height.
//   - error: An error if the BMP data is invalid or unsupported.
func LoadBmp(data []byte) (*BMP, error) {
	// Ensure the data is large enough to contain the BMP headers
	if len(data) < 54 { // Minimum size for BMP headers (14 bytes for file header + 40 bytes for info header)
		return nil, fmt.Errorf("invalid BMP data: too small")
	}

	// Read the BMP file header
	fileHeader := bitmapHeader{
		Type:      binary.LittleEndian.Uint16(data[0:2]),
		Size:      binary.LittleEndian.Uint32(data[2:6]),
		Reserved1: binary.LittleEndian.Uint16(data[6:8]),
		Reserved2: binary.LittleEndian.Uint16(data[8:10]),
		OffBits:   binary.LittleEndian.Uint32(data[10:14]),
	}
	if fileHeader.Type != 0x4D42 { // 'BM'
		return nil, fmt.Errorf("invalid BMP file type: 0x%x", fileHeader.Type)
	}

	// Read the BMP info header
	infoHeader := bitmapInfoHeader{
		BiSize:          binary.LittleEndian.Uint32(data[14:18]),
		BiWidth:         int32(binary.LittleEndian.Uint32(data[18:22])),
		BiHeight:        int32(binary.LittleEndian.Uint32(data[22:26])),
		BiPlanes:        binary.LittleEndian.Uint16(data[26:28]),
		BiBitCount:      binary.LittleEndian.Uint16(data[28:30]),
		BiCompression:   binary.LittleEndian.Uint32(data[30:34]),
		BiSizeImage:     binary.LittleEndian.Uint32(data[34:38]),
		BiXPelsPerMeter: int32(binary.LittleEndian.Uint32(data[38:42])),
		BiYPelsPerMeter: int32(binary.LittleEndian.Uint32(data[42:46])),
		BiClrUsed:       binary.LittleEndian.Uint32(data[46:50]),
		BiClrImportant:  binary.LittleEndian.Uint32(data[50:54]),
	}

	// Debugging: Print out the info header details, will delete later
	fmt.Println("BMP Info Header Details:")
	fmt.Printf("  BiSize: %d\n", infoHeader.BiSize)
	fmt.Printf("  BiWidth: %d\n", infoHeader.BiWidth)
	fmt.Printf("  BiHeight: %d\n", infoHeader.BiHeight)
	fmt.Printf("  BiPlanes: %d\n", infoHeader.BiPlanes)
	fmt.Printf("  BiBitCount: %d\n", infoHeader.BiBitCount)
	fmt.Printf("  BiCompression: %d\n", infoHeader.BiCompression)
	fmt.Printf("  BiSizeImage: %d\n", infoHeader.BiSizeImage)
	fmt.Printf("  BiXPelsPerMeter: %d\n", infoHeader.BiXPelsPerMeter)
	fmt.Printf("  BiYPelsPerMeter: %d\n", infoHeader.BiYPelsPerMeter)
	fmt.Printf("  BiClrUsed: %d\n", infoHeader.BiClrUsed)
	fmt.Printf("  BiClrImportant: %d\n", infoHeader.BiClrImportant)
	fmt.Println("BMP File Header Details:")
	fmt.Printf("  Type: %x\n", fileHeader.Type)
	fmt.Printf("  Size: %d\n", fileHeader.Size)
	fmt.Printf("  Reserved1: %d\n", fileHeader.Reserved1)
	fmt.Printf("  Reserved2: %d\n", fileHeader.Reserved2)
	fmt.Printf("  OffBits: %d\n", fileHeader.OffBits)

	// Validate the BMP format
	if infoHeader.BiCompression != 0 {
		return nil, fmt.Errorf("unsupported BMP format (must be uncompressed)")
	}

	switch infoHeader.BiBitCount {
	case 32:
		return processBmp32bit(data, fileHeader, infoHeader)
	case 24:
		return processBmp24bit(data, fileHeader, infoHeader)
	case 16:
		return processBmp16bit(data, fileHeader, infoHeader)
	case 8:
		return processBmp8bit(data, fileHeader, infoHeader)
	case 4:
		return processBmp4bit(data, fileHeader, infoHeader)
	case 1:
		return processBmp1bit(data, fileHeader, infoHeader)
	default:
		return nil, fmt.Errorf("unsupported BMP bit count: %d", infoHeader.BiBitCount)
	}
}

func buildBitMapInfoHeader(width, height, ppmX, ppmY int32, bitCount uint16, compressionMode uint32) *bitmapInfoHeader {
	return &bitmapInfoHeader{
		BiSize:          uint32(unsafe.Sizeof(bitmapInfoHeader{})),
		BiWidth:         width,
		BiHeight:        -height,
		BiPlanes:        1,
		BiBitCount:      bitCount,
		BiCompression:   compressionMode,
		BiXPelsPerMeter: ppmX,
		BiYPelsPerMeter: ppmY,
	}
}

func buildBitMapHeader(headerSize, dataSize uint32) *bitmapHeader {
	return &bitmapHeader{
		Type:    0x4D42, // 'BM'
		Size:    uint32(14 + headerSize + dataSize),
		OffBits: 14 + headerSize,
	}
}

func calcPixelsPerMeter(dpi float64) int32 {
	return int32(dpi * 39.3701)
}

func calcBytesPerPixel(bitCount int) int {
	if bitCount >= 8 {
		return bitCount / 8
	} else {
		return 1 // For 1-bit and 4-bit BMPs, treat as 1 byte per pixel for row size calculation
	}
}

func calcBmpSize(width, height, bytesPerPixel, bitCount int) int {
	var rowSize int
	switch bitCount {
	case 1:
		rowSize = ((width+7)/8 + 3) & ^3 // 1 bit per pixel, 8 pixels per byte
	case 4:
		rowSize = ((width+1)/2 + 3) & ^3 // 4 bits per pixel, 2 pixels per byte
	default:
		rowSize = (width*bytesPerPixel + 3) & ^3 // For 8-bit, 24-bit, and 32-bit BMPs
	}

	return rowSize * height
}

func processBmp32bit(data []byte, fileHeader bitmapHeader, infoHeader bitmapInfoHeader) (*BMP, error) {
	// Extract dimensions
	width := int(infoHeader.BiWidth)
	height := int(infoHeader.BiHeight)
	if height < 0 {
		height = -height // Convert to positive for consistent processing
	}

	// Calculate the pixel data offset and size
	pixelDataOffset := int(fileHeader.OffBits)
	rowSize := (width*4 + 3) & ^3 // Row size with padding
	dataSize := rowSize * height

	// Ensure the pixel data is within bounds
	if pixelDataOffset+dataSize > len(data) {
		return nil, fmt.Errorf("invalid BMP data: pixel data out of bounds")
	}

	// Extract the raw pixel data
	pixelData := data[pixelDataOffset : pixelDataOffset+dataSize]

	return &BMP{FileHeader: fileHeader, InfoHeader: infoHeader, Data: pixelData, Width: width, Height: height}, nil
}

func processBmp24bit(data []byte, fileHeader bitmapHeader, infoHeader bitmapInfoHeader) (*BMP, error) {
	// Extract dimensions
	width := int(infoHeader.BiWidth)
	height := int(infoHeader.BiHeight)
	if height < 0 {
		height = -height // Convert to positive for consistent processing
	}

	// Calculate the pixel data offset and size
	pixelDataOffset := int(fileHeader.OffBits)
	rowSize := ((width*3 + 3) / 4) * 4 // Row size with padding (3 bytes per pixel)
	dataSize := rowSize * height

	// Ensure the pixel data is within bounds
	if pixelDataOffset+dataSize > len(data) {
		return nil, fmt.Errorf("invalid BMP data: pixel data out of bounds")
	}

	// Extract the raw pixel data, including padding bytes
	pixelData := data[pixelDataOffset : pixelDataOffset+dataSize]

	// Debugging: Print calculated values
	fmt.Printf("processBmp24bit Debugging:\n")
	fmt.Printf("  Width: %d, Height: %d\n", width, height)
	fmt.Printf("  RowSize: %d, DataSize: %d\n", rowSize, dataSize)
	fmt.Printf("  PixelDataOffset: %d, TotalDataLength: %d\n", pixelDataOffset, len(data))

	return &BMP{FileHeader: fileHeader, InfoHeader: infoHeader, Data: pixelData, Width: width, Height: height}, nil
}

func processBmp16bit(data []byte, fileHeader bitmapHeader, infoHeader bitmapInfoHeader) (*BMP, error) {
	// Extract dimensions
	width := int(infoHeader.BiWidth)
	height := int(infoHeader.BiHeight)
	if height < 0 {
		height = -height // Convert to positive for consistent processing
	}

	// Calculate the pixel data offset and size
	pixelDataOffset := int(fileHeader.OffBits)
	rowSize := (width*2 + 3) & ^3 // Row size with padding (2 bytes per pixel)
	dataSize := rowSize * height

	// Ensure the pixel data is within bounds
	if pixelDataOffset+dataSize > len(data) {
		return nil, fmt.Errorf("invalid BMP data: pixel data out of bounds")
	}

	// Extract the raw pixel data
	rawPixelData := data[pixelDataOffset : pixelDataOffset+dataSize]

	// Convert the padded rows into a contiguous pixel array
	pixelData := make([]byte, width*height*3) // 3 bytes per pixel (RGB format)
	for y := 0; y < height; y++ {
		srcOffset := y * rowSize
		dstOffset := y * width * 3
		for x := 0; x < width; x++ {
			// Read 2 bytes per pixel
			pixelOffset := srcOffset + x*2
			pixel := binary.LittleEndian.Uint16(rawPixelData[pixelOffset : pixelOffset+2])

			// Extract RGB values (assuming 5-6-5 format)
			red := uint8((pixel>>11)&0x1F) << 3  // 5 bits for Red
			green := uint8((pixel>>5)&0x3F) << 2 // 6 bits for Green
			blue := uint8(pixel&0x1F) << 3       // 5 bits for Blue

			// Store the RGB values in the pixel data array
			pixelData[dstOffset+x*3+0] = blue
			pixelData[dstOffset+x*3+1] = green
			pixelData[dstOffset+x*3+2] = red
		}
	}

	return &BMP{FileHeader: fileHeader, InfoHeader: infoHeader, Data: pixelData, Width: width, Height: height}, nil
}

func processBmp8bit(data []byte, fileHeader bitmapHeader, infoHeader bitmapInfoHeader) (*BMP, error) {
	// Extract dimensions
	width := int(infoHeader.BiWidth)
	height := int(infoHeader.BiHeight)
	if height < 0 {
		height = -height // Convert to positive for consistent processing
	}

	// Calculate the pixel data offset and size
	pixelDataOffset := int(fileHeader.OffBits)
	rowSize := (width + 3) & ^3 // Row size with padding (1 byte per pixel)
	dataSize := rowSize * height

	// Ensure the pixel data is within bounds
	if pixelDataOffset+dataSize > len(data) {
		return nil, fmt.Errorf("invalid BMP data: pixel data out of bounds")
	}

	// Extract the color table
	colorTableSize := int(infoHeader.BiClrUsed)
	if colorTableSize == 0 {
		colorTableSize = 256 // Default to 256 colors for 8-bit BMPs
	}
	colorTableOffset := 14 + int(infoHeader.BiSize) // File header (14 bytes) + Info header size
	colorTable := data[colorTableOffset : colorTableOffset+colorTableSize*4]

	// Extract the raw pixel data
	rawPixelData := data[pixelDataOffset : pixelDataOffset+dataSize]

	// Convert the indexed pixel data into RGB format
	pixelData := make([]byte, width*height*3) // 3 bytes per pixel (RGB format)
	for y := 0; y < height; y++ {
		srcOffset := y * rowSize
		dstOffset := y * width * 3
		for x := 0; x < width; x++ {
			// Get the color index
			colorIndex := rawPixelData[srcOffset+x]

			// Look up the RGB values in the color table
			blue := colorTable[colorIndex*4+0]
			green := colorTable[colorIndex*4+1]
			red := colorTable[colorIndex*4+2]

			// Store the RGB values in the pixel data array
			pixelData[dstOffset+x*3+0] = blue
			pixelData[dstOffset+x*3+1] = green
			pixelData[dstOffset+x*3+2] = red
		}
	}

	return &BMP{FileHeader: fileHeader, InfoHeader: infoHeader, Data: pixelData, Width: width, Height: height}, nil
}

func processBmp4bit(data []byte, fileHeader bitmapHeader, infoHeader bitmapInfoHeader) (*BMP, error) {
	// Extract dimensions
	width := int(infoHeader.BiWidth)
	height := int(infoHeader.BiHeight)
	if height < 0 {
		height = -height // Convert to positive for consistent processing
	}

	// Calculate the pixel data offset and size
	pixelDataOffset := int(fileHeader.OffBits)
	rowSize := ((width+1)/2 + 3) & ^3 // Row size with padding (4 bits per pixel, 2 pixels per byte)
	dataSize := rowSize * height

	// Ensure the pixel data is within bounds
	if pixelDataOffset+dataSize > len(data) {
		return nil, fmt.Errorf("invalid BMP data: pixel data out of bounds")
	}

	// Extract the color table
	colorTableSize := int(infoHeader.BiClrUsed)
	if colorTableSize == 0 {
		colorTableSize = 16 // Default to 16 colors for 4-bit BMPs
	}
	colorTableOffset := 14 + int(infoHeader.BiSize) // File header (14 bytes) + Info header size
	colorTable := data[colorTableOffset : colorTableOffset+colorTableSize*4]

	// Extract the raw pixel data
	rawPixelData := data[pixelDataOffset : pixelDataOffset+dataSize]

	// Convert the indexed pixel data into RGB format
	pixelData := make([]byte, width*height*3) // 3 bytes per pixel (RGB format)
	for y := 0; y < height; y++ {
		srcOffset := y * rowSize
		dstOffset := y * width * 3
		for x := 0; x < width; x++ {
			// Get the color index (4 bits per pixel)
			byteIndex := srcOffset + x/2
			colorIndex := uint8(0)
			if x%2 == 0 {
				// High nibble
				colorIndex = rawPixelData[byteIndex] >> 4
			} else {
				// Low nibble
				colorIndex = rawPixelData[byteIndex] & 0x0F
			}

			// Look up the RGB values in the color table
			blue := colorTable[colorIndex*4+0]
			green := colorTable[colorIndex*4+1]
			red := colorTable[colorIndex*4+2]

			// Store the RGB values in the pixel data array
			pixelData[dstOffset+x*3+0] = blue
			pixelData[dstOffset+x*3+1] = green
			pixelData[dstOffset+x*3+2] = red
		}
	}

	return &BMP{FileHeader: fileHeader, InfoHeader: infoHeader, Data: pixelData, Width: width, Height: height}, nil
}

func processBmp1bit(data []byte, fileHeader bitmapHeader, infoHeader bitmapInfoHeader) (*BMP, error) {
	// Extract dimensions
	width := int(infoHeader.BiWidth)
	height := int(infoHeader.BiHeight)
	if height < 0 {
		height = -height // Convert to positive for consistent processing
	}

	// Calculate the pixel data offset and size
	pixelDataOffset := int(fileHeader.OffBits)
	rowSize := ((width+7)/8 + 3) & ^3 // Row size with padding (1 bit per pixel, 8 pixels per byte)
	dataSize := rowSize * height

	// Ensure the pixel data is within bounds
	if pixelDataOffset+dataSize > len(data) {
		return nil, fmt.Errorf("invalid BMP data: pixel data out of bounds")
	}

	// Extract the color table
	colorTableSize := int(infoHeader.BiClrUsed)
	if colorTableSize == 0 {
		colorTableSize = 2 // Default to 2 colors for 1-bit BMPs
	}
	colorTableOffset := 14 + int(infoHeader.BiSize) // File header (14 bytes) + Info header size
	colorTable := data[colorTableOffset : colorTableOffset+colorTableSize*4]

	// Extract the raw pixel data
	rawPixelData := data[pixelDataOffset : pixelDataOffset+dataSize]

	// Convert the indexed pixel data into RGB format
	pixelData := make([]byte, width*height*3) // 3 bytes per pixel (RGB format)
	for y := 0; y < height; y++ {
		srcOffset := y * rowSize
		dstOffset := y * width * 3
		for x := 0; x < width; x++ {
			// Get the color index (1 bit per pixel)
			byteIndex := srcOffset + x/8
			bitIndex := 7 - (x % 8) // Bits are stored from MSB to LSB
			colorIndex := (rawPixelData[byteIndex] >> bitIndex) & 0x01

			// Look up the RGB values in the color table
			blue := colorTable[colorIndex*4+0]
			green := colorTable[colorIndex*4+1]
			red := colorTable[colorIndex*4+2]

			// Store the RGB values in the pixel data array
			pixelData[dstOffset+x*3+0] = blue
			pixelData[dstOffset+x*3+1] = green
			pixelData[dstOffset+x*3+2] = red
		}
	}

	return &BMP{FileHeader: fileHeader, InfoHeader: infoHeader, Data: pixelData, Width: width, Height: height}, nil
}
