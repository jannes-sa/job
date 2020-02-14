package job

import "sync"

type storage struct {
	routines map[string]map[int]interface{}
	mu       *sync.RWMutex
}

func (s *storage) numOfTask(routineName string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.routines[routineName])
}

func (s *storage) setTasks(routineName string, tasks map[int]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.routines[routineName] = tasks
}

func (s *storage) getTasks(routineName string) map[int]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.routines[routineName]
}

func (s *storage) deleteTask(routineName string, key int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.routines[routineName], key)
}
