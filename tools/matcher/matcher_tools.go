package matcher

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"

	"github.com/Carmen-Shannon/automation/device/display"
	"github.com/Carmen-Shannon/automation/tools"
	"github.com/Carmen-Shannon/automation/tools/worker"
)

type chunk struct {
	Data          []byte // pixel data for the chunk
	X, Y          int    // top-left coordinates of the chunk in the original BMP
	Width, Height int    // dimensions of the chunk
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
		largeRowStart := (startY+row)*largeRowSize + startX*largeBytesPerPixel
		smallRowStart := row * smallRowSize

		for col := 0; col < smallWidth; col++ {
			largePixelStart := largeRowStart + col*largeBytesPerPixel
			smallPixelStart := smallRowStart + col*smallBytesPerPixel

			// ker-chow
			largeR := float64(largeData[largePixelStart])
			largeG := float64(largeData[largePixelStart+1])
			largeB := float64(largeData[largePixelStart+2])
			smallR := float64(smallData[smallPixelStart])
			smallG := float64(smallData[smallPixelStart+1])
			smallB := float64(smallData[smallPixelStart+2])

			dr := largeR - smallR
			dg := largeG - smallG
			db := largeB - smallB

			totalError += dr*dr + dg*dg + db*db
			pixelCount++
		}
	}

	return totalError / float64(pixelCount*3)
}

// chunkBMP divides a larger BMP into dynamically sized chunks based on the size of the smaller BMP.
// Parameters:
//   - largeBMP: The larger BMP to be divided.
//   - smallWidth: The width of the smaller BMP.
//   - smallHeight: The height of the smaller BMP.
//
// Returns:
//   - []chunk: A list of chunks with their relative positions.
func chunkBMP(largeBMP display.BMP, smallWidth, smallHeight int) []chunk {
	bytesPerPixel := tools.CalcBytesPerPixel(int(largeBMP.InfoHeader.BiBitCount))
	rowSize := ((largeBMP.Width*bytesPerPixel + 3) / 4) * 4

	// Define chunk dimensions and overlap
	widthRatio := float64(largeBMP.Width) / float64(smallWidth)
	heightRatio := float64(largeBMP.Height) / float64(smallHeight)
	// chunkWidth := tools.Max(tools.Max(largeBMP.Width/25, smallWidth*2), smallWidth)
	// chunkHeight := tools.Max(tools.Max(largeBMP.Height/25, smallHeight*2), smallHeight)
	chunkWidth := int(float64(smallWidth) * math.Min(6, math.Max(2, widthRatio/4)))
	chunkWidth = tools.Min(chunkWidth, largeBMP.Width/3)
	chunkHeight := int(float64(smallHeight) * math.Min(6, math.Max(2, heightRatio/4)))
	chunkHeight = tools.Min(chunkHeight, largeBMP.Height/3)

	// For very small images, just use the image size
	if largeBMP.Width < smallWidth*6 {
		chunkWidth = largeBMP.Width
	}
	if largeBMP.Height < smallHeight*6 {
		chunkHeight = largeBMP.Height
	}
	// overlapX := smallWidth + 1
	// overlapY := smallHeight + 1
	overlapX := tools.Max(smallWidth-1, int(float64(smallWidth)/math.Max(1.5, widthRatio/8)))
	overlapY := tools.Max(smallHeight-1, int(float64(smallHeight)/math.Max(1.5, heightRatio/8)))
	if chunkWidth == largeBMP.Width {
		overlapX = smallWidth
	}
	if chunkHeight == largeBMP.Height {
		overlapY = smallHeight
	}

	// Estimate the number of chunks for preallocation
	estimatedChunks := ((largeBMP.Height + chunkHeight - overlapY - 1) / (chunkHeight - overlapY)) *
		((largeBMP.Width + chunkWidth - overlapX - 1) / (chunkWidth - overlapX))
	chunks := make([]chunk, 0, estimatedChunks)

	// Preallocate a buffer for extracting chunk data
	buffer := make([]byte, chunkWidth*chunkHeight*bytesPerPixel)

	// Synchronization for parallel processing
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Parallelize the outer loop (y-axis)
	for y := 0; y < largeBMP.Height; y += chunkHeight - overlapY {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()

			// Local slice to collect chunks for this row
			rowChunks := []chunk{}
			// Iterate over the x-axis
			for x := 0; x < largeBMP.Width; x += chunkWidth - overlapX {
				actualChunkWidth := chunkWidth
				if x+chunkWidth > largeBMP.Width {
					actualChunkWidth = largeBMP.Width - x
				}
				if actualChunkWidth < smallWidth {
					continue // skip this chunk, too small for template
				}

				actualChunkHeight := chunkHeight
				if y+chunkHeight > largeBMP.Height {
					actualChunkHeight = largeBMP.Height - y
				}
				if actualChunkHeight < smallHeight {
					continue // skip this chunk, too small for template
				}

				// Extract chunk data using the preallocated buffer
				chunkData := extractChunk(largeBMP.Data, x, y, actualChunkWidth, actualChunkHeight, rowSize, bytesPerPixel, buffer)

				// Append the chunk to the local slice
				rowChunks = append(rowChunks, chunk{
					Data:   chunkData,
					X:      x,
					Y:      y,
					Width:  actualChunkWidth,
					Height: actualChunkHeight,
				})
			}

			// Append the row's chunks to the global slice
			mu.Lock()
			chunks = append(chunks, rowChunks...)
			mu.Unlock()
		}(y)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	return chunks
}

// extractChunk extracts the pixel data for a specific chunk from the larger BMP.
// Parameters:
//   - data: The pixel data of the larger BMP.
//   - startX, startY: The top-left position of the chunk in the larger BMP.
//   - chunkWidth, chunkHeight: The dimensions of the chunk.
//   - rowSize: The row size of the larger BMP (including padding).
//   - bytesPerPixel: The number of bytes per pixel in the larger BMP.
//
// Returns:
//   - []byte: The pixel data for the chunk.
func extractChunk(data []byte, startX, startY, chunkWidth, chunkHeight, rowSize, bytesPerPixel int, buffer []byte) []byte {
	chunkSize := chunkWidth * chunkHeight * bytesPerPixel
	if len(buffer) < chunkSize {
		buffer = make([]byte, chunkSize)
	}

	if startX*bytesPerPixel+chunkWidth*bytesPerPixel <= rowSize {
		srcOffset := startY*rowSize + startX*bytesPerPixel
		copy(buffer[:chunkSize], data[srcOffset:srcOffset+chunkSize])
	} else {
		for row := 0; row < chunkHeight; row++ {
			srcOffset := (startY+row)*rowSize + startX*bytesPerPixel
			dstOffset := row * chunkWidth * bytesPerPixel
			copy(buffer[dstOffset:dstOffset+chunkWidth*bytesPerPixel], data[srcOffset:srcOffset+chunkWidth*bytesPerPixel])
		}
	}
	return buffer[:chunkSize]
}

// normalizeBMPData ensures that the BMP data is in top-down format.
// If the BMP is bottom-up (BiHeight > 0), it flips the rows.
//
// Parameters:
//   - bmp: The BMP struct containing the pixel data.
//
// Returns:
//   - []byte: The normalized pixel data.
func normalizeBMPData(bmp display.BMP) []byte {
	// If the BMP is already top-down (negative height), return the data as-is
	if bmp.InfoHeader.BiHeight < 0 {
		return bmp.Data
	}

	// Otherwise, flip the rows to make it top-down
	bytesPerPixel := tools.CalcBytesPerPixel(int(bmp.InfoHeader.BiBitCount))
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

// sendResult sends the result to the result channel and recovers from panic if the channel is closed.
//
// Parameters:
//   - resultChan: The channel to send the result to.
//   - result: The result to be sent.
//
// Returns:
//   - bool: True if the result was sent successfully, false if the channel was closed.
func sendResult(resultChan chan struct {
	X int
	Y int
}, result struct {
	X int
	Y int
}) bool {
	defer func() {
		// Recover from panic if the channel is closed
		if r := recover(); r != nil {
			// no-op
		}
	}()

	select {
	case resultChan <- result:
		return true
	default:
		return false
	}
}

// splitChunksForWorkers divides the chunks into groups for parallel processing.
// It ensures that the chunks are distributed evenly among the workers and reverses the order of chunks for alternate groups.
//
// Parameters:
//   - chunks: The list of chunks to be divided.
//   - numWorkers: The number of workers to distribute the chunks among.
//
// Returns:
//   - [][]chunk: A slice of slices, where each inner slice contains the chunks for a specific worker.
func splitChunksForWorkers(chunks []chunk, numWorkers int) [][]chunk {
	groups := make([][]chunk, numWorkers)
	n := len(chunks)
	left, right := 0, n-1
	assignIdx := 0

	for left <= right {
		// Assign from the left
		if left <= right {
			groups[assignIdx%numWorkers] = append(groups[assignIdx%numWorkers], chunks[left])
			assignIdx++
			left++
		}
		// Assign from the right
		if left <= right {
			groups[assignIdx%numWorkers] = append(groups[assignIdx%numWorkers], chunks[right])
			assignIdx++
			right--
		}
	}
	return groups
}

// submitTasks submits tasks to the worker pool for processing the chunks of the large BMP.
// Each task processes a chunk and checks for matches with the small BMP.
//
// Parameters:
//   - worker: The worker pool to submit tasks to.
//   - chunkGroups: The groups of chunks to be processed.
//   - resultChan: The channel to send results back to the main thread.
//   - matchFound: A pointer to an atomic integer to signal when a match is found.
//   - largeData: The pixel data of the larger BMP.
//   - smallData: The pixel data of the smaller BMP.
//   - largeRowSize: The row size of the larger BMP (including padding).
//   - smallRowSize: The row size of the smaller BMP (including padding).
//   - largeBytesPerPixel: The number of bytes per pixel in the larger BMP.
//   - smallBytesPerPixel: The number of bytes per pixel in the smaller BMP.
//   - smallWidth: The width of the smaller BMP.
//   - smallHeight: The height of the smaller BMP.
//   - mseThreshold: The maximum allowable MSE for a match.
func submitTasks(pool worker.DynamicWorkerPool, chunkGroups [][]chunk, resultChan chan struct {
	X int
	Y int
}, matchFound *int32, largeData, smallData []byte, largeRowSize, smallRowSize, largeBytesPerPixel, smallBytesPerPixel, smallWidth, smallHeight int, mseThreshold float64, ctx context.Context) {
	for _, chunkGroup := range chunkGroups {
		chunkGroup := chunkGroup // Capture chunkGroup in the loop

		task := worker.Task{
			ID: len(chunkGroup),
			Do: func() (any, error) {
				for _, chunk := range chunkGroup {
					if ctx.Err() != nil {
						return nil, nil
					}
					for y := 0; y <= chunk.Height-smallHeight; y++ {
						if atomic.LoadInt32(matchFound) == 1 {
							return nil, nil
						} else if ctx.Err() != nil {
							return nil, nil
						}

						for x := 0; x <= chunk.Width-smallWidth; x++ {
							if ctx.Err() != nil {
								return nil, nil
							}
							absoluteX := chunk.X + x
							absoluteY := chunk.Y + y

							// Calculate MSE for the current window
							mse := calculateMSE(
								largeData, smallData,
								absoluteX, absoluteY,
								largeRowSize, smallRowSize,
								largeBytesPerPixel, smallBytesPerPixel,
								smallWidth, smallHeight,
							)

							// Early exit if the MSE is significantly below the threshold
							if mse <= mseThreshold/5 {
								if atomic.CompareAndSwapInt32(matchFound, 0, 1) {
									sendResult(resultChan, struct {
										X int
										Y int
									}{X: absoluteX, Y: absoluteY})
									return nil, nil
								}
							}

							// If the MSE is below the threshold, validate the match
							if mse <= mseThreshold {
								validationMSE := calculateMSE(
									largeData, smallData,
									absoluteX, absoluteY,
									largeRowSize, smallRowSize,
									largeBytesPerPixel, smallBytesPerPixel,
									smallWidth, smallHeight,
								)

								if validationMSE <= mseThreshold {
									if atomic.CompareAndSwapInt32(matchFound, 0, 1) {
										sendResult(resultChan, struct {
											X int
											Y int
										}{X: absoluteX, Y: absoluteY})
										return nil, nil
									}
								}
							}
						}
					}
				}
				return nil, nil
			},
		}
		if ctx.Err() != nil {
			return
		}
		pool.SubmitTask(task)
	}
}

// validateBMPDimensions checks if the dimensions of the small BMP are within the bounds of the large BMP.
//
// Parameters:
//   - largeBMP: The larger BMP image.
//   - smallBMP: The smaller BMP image.
//
// Returns:
//   - error: An error if the small BMP dimensions exceed the large BMP dimensions.
func validateBMPDimensions(largeBMP, smallBMP display.BMP) error {
	if smallBMP.Width > largeBMP.Width || smallBMP.Height > largeBMP.Height {
		return fmt.Errorf("small BMP dimensions exceed large BMP dimensions")
	}
	return nil
}
