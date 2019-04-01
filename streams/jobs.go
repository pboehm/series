package streams

import (
	"bytes"
	"github.com/gofrs/uuid"
	"io"
	"sync"
)

type Runnable func(io.Writer) error

type Job struct {
	Id       string `json:"id"`
	Running  bool   `json:"running"`
	Finished bool   `json:"finished"`
	Error    error  `json:"error"`
	runnable Runnable
	buffer   *bytes.Buffer
}

func NewJob(runnable Runnable) *Job {
	return &Job{
		Id:       uuid.Must(uuid.NewV4()).String(),
		Finished: false,
		Error:    nil,
		runnable: runnable,
		buffer:   &bytes.Buffer{},
	}
}

func (j *Job) Run() {
	if j.Running || j.Finished {
		return
	}

	j.Running = true
	j.Error = j.runnable(j.buffer)
	j.Running = false
	j.Finished = true
}

func (j *Job) Response() map[string]interface{} {
	var success, failure, err interface{} = nil, nil, nil
	if j.Finished {
		success = j.Error == nil
		failure = j.Error != nil
	}
	if j.Error != nil {
		err = j.Error.Error()
	}

	return map[string]interface{}{
		"id":       j.Id,
		"running":  j.Running,
		"finished": j.Finished,
		"success":  success,
		"failure":  failure,
		"error":    err,
		"output":   j.buffer.String(),
	}
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
