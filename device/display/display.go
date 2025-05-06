package display

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Display struct {
	X           int32
	Y           int32
	Width       int
	Height      int
	RefreshRate float32
	Primary     bool
}

type BMP struct {
	FileHeader bitmapHeader
	InfoHeader bitmapInfoHeader
	ColorTable [256][4]uint8
	Data       []byte
	Width      int
	Height     int
}

// ToBinary serializes the BMP struct into a byte slice in BMP format.
// It includes the file header, info header, and pixel data.
// The function returns the serialized byte slice.
//
// Returns:
//   - []byte: A byte slice containing the serialized BMP data.
func (b *BMP) ToBinary() []byte {
	// Create a buffer to hold the binary data
	var buffer bytes.Buffer

	// Serialize the file header
	binary.Write(&buffer, binary.LittleEndian, b.FileHeader.Type)      // 'BM'
	binary.Write(&buffer, binary.LittleEndian, b.FileHeader.Size)      // File size
	binary.Write(&buffer, binary.LittleEndian, b.FileHeader.Reserved1) // Reserved1
	binary.Write(&buffer, binary.LittleEndian, b.FileHeader.Reserved2) // Reserved2
	binary.Write(&buffer, binary.LittleEndian, b.FileHeader.OffBits)   // Offset to pixel data

	// Serialize the info header
	binary.Write(&buffer, binary.LittleEndian, b.InfoHeader.BiSize)
	binary.Write(&buffer, binary.LittleEndian, b.InfoHeader.BiWidth)
	binary.Write(&buffer, binary.LittleEndian, b.InfoHeader.BiHeight)
	binary.Write(&buffer, binary.LittleEndian, b.InfoHeader.BiPlanes)
	binary.Write(&buffer, binary.LittleEndian, b.InfoHeader.BiBitCount)
	binary.Write(&buffer, binary.LittleEndian, b.InfoHeader.BiCompression)
	binary.Write(&buffer, binary.LittleEndian, b.InfoHeader.BiSizeImage)
	binary.Write(&buffer, binary.LittleEndian, b.InfoHeader.BiXPelsPerMeter)
	binary.Write(&buffer, binary.LittleEndian, b.InfoHeader.BiYPelsPerMeter)
	binary.Write(&buffer, binary.LittleEndian, b.InfoHeader.BiClrUsed)
	binary.Write(&buffer, binary.LittleEndian, b.InfoHeader.BiClrImportant)

	// Serialize the color table if BiBitCount is 8
	if b.InfoHeader.BiBitCount == 8 {
		for _, entry := range b.ColorTable {
			binary.Write(&buffer, binary.LittleEndian, entry)
		}
	}

	// Append the pixel data
	buffer.Write(b.Data)

	return buffer.Bytes()
}

type bitmapInfoHeader struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

type bitmapInfo struct {
	BmiHeader bitmapInfoHeader
	BmiColors [1]uint32
}

type bitmapHeader struct {
	Type      uint16
	Size      uint32
	Reserved1 uint16
	Reserved2 uint16
	OffBits   uint32
}

type virtualScreen struct {
	Left     int32
	Right    int32
	Top      int32
	Bottom   int32
	Displays []Display
}

type VirtualScreen interface {
	// CaptureBmp captures the current screen and saves the bitmap as a byte slice.
	// It accepts options to specify which display(s) to capture, if none are provided then the primary display is captured.
	//
	// Parameters:
	//   - options: Optional parameters for the display capture, such as the display to capture.
	//
	// Returns:
	//   - [][]byte: A byte slice containing the bitmap data of the captured screen.
	//   - error: An error if the capture fails.
	CaptureBmp(options ...DisplayCaptureOption) ([]BMP, error)

	// DetectDisplays detects all displays connected to the system and returns a slice of display structs.
	// It also modifies the virtual screen Displays field to include the detected displays.
	// If no displays are found, it returns an error.
	//
	// Returns:
	//   - []Display: A slice of Display structs representing the detected displays.
	//   - error: An error if the detection fails or no displays are found.
	DetectDisplays() ([]Display, error)

	// GetPrimaryDisplay retrieves the primary display from the virtual screen.
	// If no primary display is found, it returns an error.
	//
	// Returns:
	//   - Display: The primary display struct.
	//   - error: An error if no primary display is found.
	GetPrimaryDisplay() (Display, error)

	// Displays returns a slice of all displays connected to the system.
	// Returns:
	//   - []Display: A slice of Display structs representing all connected displays.
	GetDisplays() []Display

	// Left returns the left bound of the virtual screen.
	// Returns:
	//   - int32: The left bound of the virtual screen.
	GetLeft() int32

	// Right returns the right bound of the virtual screen.
	// Returns:
	//   - int32: The right bound of the virtual screen.
	GetRight() int32

	// Top returns the top bound of the virtual screen.
	// Returns:
	//   - int32: The top bound of the virtual screen.
	GetTop() int32

	// Bottom returns the bottom bound of the virtual screen.
	// Returns:
	//   - int32: The bottom bound of the virtual screen.
	GetBottom() int32
}

var _ VirtualScreen = (*virtualScreen)(nil) // compile-time check to ensure that virtualScreen implements VirtualScreen

func (vs *virtualScreen) GetPrimaryDisplay() (Display, error) {
	displays := vs.Displays

	if displays == nil {
		displays, err := vs.DetectDisplays()
		if err != nil || len(displays) == 0 {
			return Display{}, err
		}
	}
	for _, display := range displays {
		if display.Primary {
			return display, nil
		}
	}
	return Display{}, errors.New("no primary display found")
}

func (vs *virtualScreen) GetDisplays() []Display {
	return vs.Displays
}

func (vs *virtualScreen) GetLeft() int32 {
	return vs.Left
}

func (vs *virtualScreen) GetRight() int32 {
	return vs.Right
}

func (vs *virtualScreen) GetTop() int32 {
	return vs.Top
}

func (vs *virtualScreen) GetBottom() int32 {
	return vs.Bottom
}
