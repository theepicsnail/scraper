package main

import "sync"

type Task func(chan<- Task)

type Runner struct {
	tasks chan Task
	quit  chan bool
	sync.WaitGroup
}

func NewRunner() *Runner {
	return &Runner{
		tasks: make(chan Task, task_queue_count),
		quit:  make(chan bool),
	}
}

func (r *Runner) Run() {
	// Start the workers
	for i := 0; i < worker_count; i++ {
		r.Add(1)
		go r.spawnWorker()
	}

	// Wait for them to all finish
	r.Wait()
}

func (r *Runner) spawnWorker() {
	defer r.WaitGroup.Done()
	for task := range r.tasks {
		task(r.tasks)
	}
}
