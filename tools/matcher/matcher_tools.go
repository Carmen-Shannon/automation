package matcher

import (
	"automation/device/display"
	"automation/tools"
	"automation/tools/worker"
	"fmt"
	"sync"
	"sync/atomic"
)

type chunk struct {
	Data          []byte // pixel data for the chunk
	X, Y          int    // top-left coordinates of the chunk in the original BMP
	Width, Height int    // dimensions of the chunk
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
	chunkWidth := tools.Max(largeBMP.Width/10, smallWidth*2)
	chunkHeight := tools.Max(largeBMP.Height/10, smallHeight*2)
	overlapX := smallWidth
	overlapY := smallHeight

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
			actualChunkHeight := chunkHeight
			if y+chunkHeight > largeBMP.Height {
				actualChunkHeight = largeBMP.Height - y
			}

			// Iterate over the x-axis
			for x := 0; x < largeBMP.Width; x += chunkWidth - overlapX {
				actualChunkWidth := chunkWidth
				if x+chunkWidth > largeBMP.Width {
					actualChunkWidth = largeBMP.Width - x
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
	// Use the provided buffer if it’s large enough
	chunkSize := chunkWidth * chunkHeight * bytesPerPixel
	if len(buffer) < chunkSize {
		buffer = make([]byte, chunkSize)
	}

	// Check if the chunk data is contiguous in memory
	if startX*bytesPerPixel+chunkWidth*bytesPerPixel <= rowSize {
		// Contiguous case: Use a single copy operation
		srcOffset := startY*rowSize + startX*bytesPerPixel
		copy(buffer[:chunkSize], data[srcOffset:srcOffset+chunkSize])
	} else {
		// Non-contiguous case: Copy row by row
		for row := 0; row < chunkHeight; row++ {
			// Calculate the source offset in the larger BMP
			srcOffset := (startY+row)*rowSize + startX*bytesPerPixel

			// Calculate the destination offset in the chunk
			dstOffset := row * chunkWidth * bytesPerPixel

			// Copy the row data
			copy(buffer[dstOffset:dstOffset+chunkWidth*bytesPerPixel], data[srcOffset:srcOffset+chunkWidth*bytesPerPixel])
		}
	}

	return buffer[:chunkSize]
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
	for i, chunk := range chunks {
		groupIndex := i % numWorkers
		groups[groupIndex] = append(groups[groupIndex], chunk)
	}

	// Reverse the order of chunks for alternate groups
	for i := 1; i < numWorkers; i += 2 {
		for j, k := 0, len(groups[i])-1; j < k; j, k = j+1, k-1 {
			groups[i][j], groups[i][k] = groups[i][k], groups[i][j]
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
}, matchFound *int32, largeData, smallData []byte, largeRowSize, smallRowSize, largeBytesPerPixel, smallBytesPerPixel, smallWidth, smallHeight int, mseThreshold float64) {
	for _, chunkGroup := range chunkGroups {
		chunkGroup := chunkGroup // Capture chunkGroup in the loop

		task := worker.Task{
			ID: len(chunkGroup),
			Do: func() (any, error) {
				for _, chunk := range chunkGroup {
					for y := 0; y <= chunk.Height-smallHeight; y++ {
						if atomic.LoadInt32(matchFound) == 1 {
							return nil, nil
						}

						for x := 0; x <= chunk.Width-smallWidth; x++ {
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

		pool.SubmitTask(task)
	}
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
