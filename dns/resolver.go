package dns

import (
	"net"

	"github.com/russtone/utils/jobqueue"
)

type Resolver struct {
	jobqueue.Queue

	pool *pool
}

func NewResolver(servers []net.IP, workersCount int, rateLimit float64, capacity int) *Resolver {

	pool := &pool{
		servers:   make(chan server, len(servers)),
		rateLimit: rateLimit,
	}

	for _, ip := range servers {
		pool.add(ip)
	}

	r := &Resolver{
		pool: pool,
	}

	r.Queue = jobqueue.New(r, workersCount, capacity)

	return r
}

func (r *Resolver) Process(j interface{}) (interface{}, bool, error) {
	job, ok := j.(*Job)
	if !ok || job == nil {
		return nil, false, jobqueue.ErrInvalidJob
	}

	ns := r.pool.take()

	defer func() {
		r.pool.release(ns)
	}()

	qtype := job.qtype()

	answer, err := ns.query(job.Name, qtype)

	if err != nil {
		return nil, true, err
	}

	job.setAnswer(qtype, answer)

	res := Result{
		Name:    job.Name,
		Answers: job.Answers,
		Meta:    job.Meta,
	}

	return res, !job.done(), nil
}

func (r *Resolver) Schedule(name string, qtypes []string, meta map[string]interface{}) {
	r.Queue.Schedule(&Job{
		Name:    name,
		Qtypes:  qtypes,
		Answers: make(map[string][]string),
		Meta:    meta,
	})
}

//
// Job
//

type Job struct {
	Name    string
	Qtypes  []string
	Answers map[string][]string
	Meta    map[string]interface{}

	qtypeIdx int
}

func (j *Job) done() bool {
	return len(j.Answers) == len(j.Qtypes)
}

func (j *Job) qtype() string {
	if j.qtypeIdx > len(j.Qtypes) {
		panic("Invalid qtype intex")
	}

	return j.Qtypes[j.qtypeIdx]
}

func (j *Job) setAnswer(qtype string, answer []string) {
	j.Answers[qtype] = answer
	j.qtypeIdx++
}

//
// Result
//

type Result struct {
	Name    string
	Answers map[string][]string
	Meta    map[string]interface{}
}

func (r *Result) IsEmpty() bool {
	count := 0

	for _, ans := range r.Answers {
		count += len(ans)
	}

	return count == 0
}
