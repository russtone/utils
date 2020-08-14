package jobqueue_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/russtone/utils/jobqueue"
)

type ProcessorMock struct {
	mock.Mock
	errID    int
	errMet   bool
	progress int
	wg       sync.WaitGroup
}

var errTest = errors.New("test")

func (p *ProcessorMock) Process(job interface{}) (interface{}, bool, error) {

	p.Called(job)

	j, ok := job.(*Job)
	if !ok {
		return nil, false, jobqueue.ErrInvalidJob
	}

	if j.ID == p.errID && !p.errMet {
		p.errMet = true
		// Retry once on error.
		return nil, true, errTest
	}

	// Must be after error check to avoid calling twice for id == errID.
	if j.ID < p.progress {
		defer p.wg.Done()
	}

	return Result{j}, !j.done(), nil
}

type Job struct {
	ID    int
	retry bool
}

type Result struct {
	*Job
}

func (j *Job) done() bool {
	if j.retry {
		j.retry = false
		return false
	}

	return true
}

func TestJobqueue(t *testing.T) {

	tests := []struct {
		jobs     int
		workers  int
		errID    int
		retryID  int
		destID   int
		progress int
	}{
		{jobs: 10, workers: 3, errID: 5, retryID: 6, destID: 7, progress: 5},
		{jobs: 10, workers: 3, errID: 6, retryID: 6, destID: 7, progress: 5},
		{jobs: 100, workers: 5, errID: 5, retryID: 50, destID: 69, progress: 10},
		{jobs: 100, workers: 50, errID: 50, retryID: 80, destID: 33, progress: 80},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {

			if tt.progress > tt.jobs {
				panic("invald test case: progress must be less then jobs count")
			}

			if tt.errID > tt.jobs || tt.retryID > tt.jobs || tt.destID > tt.jobs {
				panic("invalid test case: id must be less then jobs count")
			}

			proc := &ProcessorMock{
				errID:    tt.errID,
				progress: tt.progress,
			}

			proc.wg.Add(tt.progress)

			for i := 0; i < tt.jobs; i++ {
				if i == tt.retryID && i == tt.errID {
					proc.On("Process", &Job{i, true}).Return().Twice()
					proc.On("Process", &Job{i, false}).Return().Once()
				} else if i == tt.retryID {
					proc.On("Process", &Job{i, true}).Return().Once()
					proc.On("Process", &Job{i, false}).Return().Once()
				} else if i == tt.errID {
					proc.On("Process", &Job{i, false}).Return().Twice()
				} else {
					proc.On("Process", &Job{i, false}).Return().Once()
				}
			}

			queue := jobqueue.New(proc, tt.workers, tt.jobs)

			queue.Start()

			wg := sync.WaitGroup{}
			wg.Add(2)

			// Process errors
			go func() {
				defer wg.Done()

				var err error

				count := 0

				for queue.Err(&err) {
					if err != nil {
						// The only error that could happen is testErr.
						assert.Error(t, errTest, err)
						count++
					}
				}

				// There must be only 1 error.
				assert.Equal(t, 1, count)
			}()

			// Process results
			go func() {
				defer wg.Done()

				var res Result

				count := 0

				for queue.Next(&res) {
					count += 1
				}

				// All jobs must be proccessed.
				assert.Equal(t, count, tt.jobs)
			}()

			queue.Add(tt.jobs)

			var dest Result

			for i := 0; i < tt.progress; i++ {
				job := &Job{i, i == tt.retryID}

				if i != tt.destID {
					queue.Schedule(job)
				} else {
					queue.ScheduleDest(job, &dest)
				}
			}

			// Wait for jobs to finish to compare progress.
			proc.wg.Wait()
			assert.Equal(t, float64(tt.progress)/float64(tt.jobs), queue.Progress())

			for i := tt.progress; i < tt.jobs; i++ {
				job := &Job{i, i == tt.retryID}

				if i != tt.destID {
					queue.Schedule(job)
				} else {
					queue.ScheduleDest(job, &dest)
				}
			}

			queue.WaitJobs()

			// Check job with dest.
			assert.Equal(t, Result{&Job{tt.destID, false}}, dest)

			queue.Stop()

			queue.WaitWorkers()

			// Wait erors and output processors.
			wg.Wait()
		})
	}

}
