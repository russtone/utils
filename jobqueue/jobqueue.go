package jobqueue

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrInvalidJob = errors.New("invalid job")
)

type Queue interface {
	Add(delta int)
	Schedule(job interface{})
	ScheduleDest(job interface{}, dest interface{})

	Start()
	Stop()

	Next(dist interface{}) bool
	Err(err *error) bool

	WaitWorkers()
	WaitJobs()

	Progress() float64
	Speed() float64
}

type Processor interface {
	Process(interface{}) (interface{}, bool, error)
}

type queue struct {
	processor Processor

	workersCount int
	workersWG    sync.WaitGroup

	jobsWG        sync.WaitGroup
	jobsCount     uint64
	jobsProcessed uint64

	todo chan job
	done chan interface{}
	errs chan error

	createdAt time.Time
}

type job struct {
	job  interface{}
	dest interface{}
}

func New(processor Processor, workersCount int, capacity int) Queue {
	return &queue{
		processor:    processor,
		workersCount: workersCount,
		todo:         make(chan job, capacity),
		done:         make(chan interface{}, capacity),
		errs:         make(chan error, capacity),
		createdAt:    time.Now(),
	}
}

func (jq *queue) Add(delta int) {
	atomic.AddUint64(&jq.jobsCount, uint64(delta))
	jq.jobsWG.Add(delta)
}

func (jq *queue) Schedule(j interface{}) {
	jq.todo <- job{job: j, dest: nil}
}

func (jq *queue) ScheduleDest(j interface{}, dest interface{}) {
	jq.todo <- job{job: j, dest: dest}
}

func (jq *queue) Next(dest interface{}) bool {
	res, ok := <-jq.done
	if !ok {
		return false
	}

	jq.setDest(dest, res)

	return true
}

func (jq *queue) Err(e *error) bool {
	err, ok := <-jq.errs
	if !ok {
		return false
	}

	*e = err
	return true
}

func (jq *queue) Start() {
	for i := 0; i < jq.workersCount; i++ {
		jq.workersWG.Add(1)
		go jq.worker(i)
	}

	go func() {
		jq.workersWG.Wait()
		close(jq.done)
		close(jq.errs)
	}()
}

func (jq *queue) Stop() {
	close(jq.todo)
}

func (jq *queue) WaitWorkers() {
	jq.workersWG.Wait()
}

func (jq *queue) WaitJobs() {
	jq.jobsWG.Wait()
}

func (jq *queue) Progress() float64 {
	processed := atomic.LoadUint64(&jq.jobsProcessed)
	return float64(processed) / float64(jq.jobsCount)
}

func (jq *queue) Speed() float64 {
	return float64(jq.jobsProcessed) / time.Since(jq.createdAt).Seconds()
}

func (jq *queue) worker(id int) {
	defer jq.workersWG.Done()

	for j := range jq.todo {
		jq.process(j)
	}
}

func (jq *queue) process(j job) {
	res, retry, err := jq.processor.Process(j.job)

	if err != nil {
		jq.errs <- err
	}

	if retry {
		jq.retry(j)
		return
	}

	jq.jobsWG.Done()
	atomic.AddUint64(&jq.jobsProcessed, 1)

	if err != nil {
		return
	}

	if j.dest != nil {
		jq.setDest(j.dest, res)
	}

	jq.done <- res
}

func (jq *queue) retry(j job) {
	go func() {
		jq.todo <- j
	}()
}

func (jq *queue) setDest(destination interface{}, result interface{}) {
	dst := reflect.ValueOf(destination)

	if dst.Kind() != reflect.Ptr {
		panic("destination must pass a pointer, not a value")
	}

	if result == nil {
		panic("result must not be nil")
	}

	res := reflect.ValueOf(result)

	if dst.Elem().Type() != res.Type() {
		panic(fmt.Sprintf("invalid destination type: expected *%v, got %v", res.Type(), dst.Type()))
	}

	dst.Elem().Set(res)
}
