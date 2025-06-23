package tasks

import (
	"errors"
	"fmt"
)

var (
	NotFoundErr = errors.New("not found")
)

type MemStore struct {
	list map[string]Task
}

func NewMemStore() *MemStore {
	list := make(map[string]Task)
	return &MemStore{
		list,
	}
}

func (m MemStore) Add(id string, task Task) error {
	m.list[id] = task
	fmt.Println(task.String())
	return nil
}

func (m MemStore) Get(id string) (Task, error) {

	if val, ok := m.list[id]; ok {
		return val, nil
	}

	return Task{}, NotFoundErr
}

func (m MemStore) List() (map[string]Task, error) {
	return m.list, nil
}

func (m MemStore) Update(id string, task Task) error {
	fmt.Println(task.String())
	if _, ok := m.list[id]; ok {
		m.list[id] = task
		return nil
	}

	return NotFoundErr
}

func (m MemStore) Remove(id string) error {
	delete(m.list, id)
	return nil
}
