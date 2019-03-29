package streams

import (
	"github.com/satori/go.uuid"
	"sync"
)

type Runnable func() error

type Job struct {
	Id       string `json:"id"`
	Running  bool   `json:"running"`
	Finished bool   `json:"finished"`
	Error    error  `json:"error"`
	runnable Runnable
}

func NewJob(runnable Runnable) *Job {
	return &Job{
		Id:       uuid.NewV4().String(),
		Finished: false,
		Error:    nil,
		runnable: runnable,
	}
}

func (j *Job) Run() {
	if j.Running || j.Finished {
		return
	}

	j.Running = true
	j.Error = j.runnable()
	j.Running = false
	j.Finished = true
}

type JobPool struct {
	data   map[string]*Job
	rwLock sync.RWMutex
}

func NewJobPool() *JobPool {
	return &JobPool{
		data:   map[string]*Job{},
		rwLock: sync.RWMutex{},
	}
}

func (jp *JobPool) Get(id string) (*Job, bool) {
	jp.rwLock.RLock()
	defer jp.rwLock.RUnlock()
	job, found := jp.data[id]
	return job, found
}

func (jp *JobPool) Set(id string, job *Job) {
	jp.rwLock.Lock()
	defer jp.rwLock.Unlock()
	jp.data[id] = job
}
