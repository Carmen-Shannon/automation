package matcher

import (
	"automation/device/display"
	"automation/tools"
	"automation/tools/worker"
	"context"
	"fmt"
	"sync"
	"time"
)

type matcher struct {
	pool worker.DynamicWorkerPool
	scan display.BMP
}

type Matcher interface {
	// FindTemplate searches for a smaller BMP within another BMP using MSE for fuzzy matching.
	// It accepts a smaller template to search for as well as various options for the search, such as timeout and threshold.
	//
	// Parameters:
	//   - template: The smaller BMP image (template) to search for.
	//   - options: Optional parameters for the search, such as MSE threshold and timeout.
	//
	// Returns:
	//   - (x, y): The top-left coordinates of the match in the larger BMP.
	//     NOTE: The coordinates are relative to the larger BMP, not the screen.
	//   - error: An error if no match is found or if the search fails.
	FindTemplate(template display.BMP, options ...FindBuilderOption) (int, int, error)

	// SetScan sets the BMP to be used for scanning.
	// This is useful for updating the scan area without creating a new matcher instance.
	// It will stop the current worker pool and clear the task queue before setting the new BMP, as to stop any ongoing matching tasks.
	// Because of this, it has to wait for all workers to finish before setting the new BMP.
	//
	// Parameters:
	//   - bmp: The new BMP to set for scanning.
	SetScan(bmp display.BMP)
}

var _ Matcher = (*matcher)(nil)

func New(bmp display.BMP) Matcher {
	return &matcher{
		pool: worker.NewDynamicWorkerPool(1, 3000),
		scan: bmp,
	}
}

func (m *matcher) FindTemplate(template display.BMP, options ...FindBuilderOption) (int, int, error) {
	startTime := time.Now()
	fbo := &findBuilderOption{}
	for _, opt := range options {
		opt(fbo)
	}
	if fbo.MSEThreshold == 0 {
		fbo.MSEThreshold = 100.0
	}
	if fbo.Timeout == 0 {
		fbo.Timeout = 500 * time.Millisecond
	}

	// Step 1: Validate inputs
	if err := validateBMPDimensions(m.scan, template); err != nil {
		return 0, 0, err
	}

	// Step 2: Normalize BMPs
	largeData, smallData := normalizeBMPData(m.scan), normalizeBMPData(template)

	// Step 3: Calculate metadata
	largeBytesPerPixel := tools.CalcBytesPerPixel(int(m.scan.InfoHeader.BiBitCount))
	smallBytesPerPixel := tools.CalcBytesPerPixel(int(template.InfoHeader.BiBitCount))
	largeRowSize := ((m.scan.Width*largeBytesPerPixel + 3) / 4) * 4
	smallRowSize := ((template.Width*smallBytesPerPixel + 3) / 4) * 4

	// Step 4: Chunk the large BMP
	chunks := chunkBMP(m.scan, template.Width, template.Height)

	// Step 5: Initialize worker pool
	numWorkers := tools.Max((len(chunks)/3)/2, 1)
	chunkGroups := splitChunksForWorkers(chunks, numWorkers)
	if numWorkers > m.pool.GetMaxWorkers() {
		diff := numWorkers - m.pool.GetMaxWorkers()
		m.pool.IncreaseMaxWorkers(diff)
	}
	if !m.pool.IsWorking() {
		m.pool.Start()
	}

	// Step 6: Submit tasks and collect results
	resultChan := make(chan struct {
		X int
		Y int
	}, 1)
	matchFound := int32(0) // Atomic flag to signal match

	// Use sync.Once to ensure the resultChan is closed exactly once
	var closeOnce sync.Once
	closeResultChan := func() {
		close(resultChan)
	}

	// Submit tasks to the worker pool
	submitTasks(m.pool, chunkGroups, resultChan, &matchFound, largeData, smallData, largeRowSize, smallRowSize, largeBytesPerPixel, smallBytesPerPixel, template.Width, template.Height, fbo.MSEThreshold)

	// Step 7: Wait for results or pool completion
	var result struct {
		X int
		Y int
	}

	ctx, cancel := context.WithTimeout(context.Background(), fbo.Timeout)
	defer cancel()
	defer m.pool.Stop()

	select {
	case res := <-resultChan:
		result = res
		closeOnce.Do(closeResultChan)
		fmt.Println("time to match: ", time.Since(startTime))
		return result.X, result.Y, nil
	case <-ctx.Done():
		closeOnce.Do(closeResultChan)
		fmt.Println("No match found (timeout).")
		return 0, 0, fmt.Errorf("no match found")
	}
}

func (m *matcher) SetScan(bmp display.BMP) {
	m.pool.ClearTaskQueue()
	m.pool.Stop()
	m.pool.Wait()

	m.scan = bmp
	m.pool.Start()
}
