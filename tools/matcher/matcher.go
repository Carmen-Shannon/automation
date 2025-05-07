package matcher

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/Carmen-Shannon/automation/device/display"
	"github.com/Carmen-Shannon/automation/tools"
	"github.com/Carmen-Shannon/automation/tools/worker"
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

// NewMatcher creates a new matcher instance with the given BMP for scanning.
// It initializes a worker pool for processing matching tasks and returns the matcher instance.
//
// Parameters:
//   - bmp: The BMP to be used for scanning. This is the larger BMP image in which to search for the template.
//
// Returns:
//   - Matcher: A new matcher instance that can be used to find templates within the specified BMP.
func NewMatcher(bmp display.BMP) Matcher {
	return &matcher{
		pool: worker.NewDynamicWorkerPool(1, 3000, 500*time.Millisecond),
		scan: bmp,
	}
}

func (m *matcher) FindTemplate(template display.BMP, options ...FindBuilderOption) (int, int, error) {
	fbo := &findBuilderOption{}
	for _, opt := range options {
		opt(fbo)
	}
	if fbo.Threshold == 0 {
		fbo.Threshold = 100.0
	}
	if fbo.Timeout == 0 {
		fbo.Timeout = 500 * time.Millisecond
	}

	if err := validateBMPDimensions(m.scan, template); err != nil {
		return 0, 0, err
	}

	largeData, smallData := normalizeBMPData(m.scan), normalizeBMPData(template)

	largeBytesPerPixel := tools.CalcBytesPerPixel(int(m.scan.InfoHeader.BiBitCount))
	smallBytesPerPixel := tools.CalcBytesPerPixel(int(template.InfoHeader.BiBitCount))
	largeRowSize := ((m.scan.Width*largeBytesPerPixel + 3) / 4) * 4
	smallRowSize := ((template.Width*smallBytesPerPixel + 3) / 4) * 4

	chunks := chunkBMP(m.scan, template.Width, template.Height)

	numWorkers := tools.Max(runtime.NumCPU()-1, 1)
	chunkGroups := splitChunksForWorkers(chunks, numWorkers)
	if numWorkers > m.pool.GetMaxWorkers() {
		diff := numWorkers - m.pool.GetMaxWorkers()
		m.pool.IncreaseMaxWorkers(diff)
	}
	if !m.pool.IsWorking() {
		m.pool.Start()
	}

	resultChan := make(chan struct {
		X int
		Y int
	}, 1)
	matchFound := int32(0)
	var closeOnce sync.Once
	closeResultChan := func() {
		close(resultChan)
	}

	ctx, cancel := context.WithTimeout(context.Background(), fbo.Timeout)
	defer cancel()
	defer m.pool.Stop()
	defer closeOnce.Do(closeResultChan)

	// Submit tasks to the worker pool
	submitTasks(m.pool, chunkGroups, resultChan, &matchFound, largeData, smallData, largeRowSize, smallRowSize, largeBytesPerPixel, smallBytesPerPixel, template.Width, template.Height, fbo.Threshold, ctx)

	for {
		select {
		case <-ctx.Done():
			return 0, 0, fmt.Errorf("no match found - timeout")
		case res := <-resultChan:
			return res.X, res.Y, nil
		}
	}
}

func (m *matcher) SetScan(bmp display.BMP) {
	m.pool.ClearTaskQueue()
	m.pool.Stop()
	m.pool.Wait()

	m.scan = bmp
	m.pool.Start()
}
