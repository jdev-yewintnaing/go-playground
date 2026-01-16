package scheduler

import (
	"fmt"
)

type Task interface {
	Run() error
}

type Manager struct {
	tasks []Task
}

func NewManager() *Manager {
	return &Manager{
		tasks: []Task{},
	}
}

func (m *Manager) Add(task Task) {
	m.tasks = append(m.tasks, task)
}

func (m *Manager) RunAll() []error {
	//var wg sync.WaitGroup
	errChan := make(chan error, len(m.tasks))
	//for _, task := range m.tasks {
	//	wg.Add(1)
	//	go func(t Task) {
	//		defer wg.Done()
	//		if err := t.Run(); err != nil {
	//			errChan <- err
	//		}
	//
	//	}(task)
	//
	//}
	//
	//wg.Wait()
	//close(errChan)

	for _, task := range m.tasks {
		go func(t Task) {
			errChan <- t.Run()
		}(task)
	}
	var results []error

	for range m.tasks {
		err := <-errChan
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
	//// 5. Collect results into a slice
	//for err := range errChan {
	//	results = append(results, err)
	//}

	return results
}

type EmailTask struct {
	Email string
}

func (e EmailTask) Run() error {
	fmt.Printf("Sending email to %s...\n", e.Email)
	return nil
}
