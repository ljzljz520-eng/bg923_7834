package navigation

import "sync"

type memoryStore struct {
	mu    sync.RWMutex
	links []Link
}

func newMemoryStore(initial []Link) *memoryStore {
	links := make([]Link, len(initial))
	copy(links, initial)
	return &memoryStore{links: links}
}

func (s *memoryStore) add(link Link) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.links = append(s.links, link)
}

func (s *memoryStore) list(group string) []Link {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Link, 0, len(s.links))
	for _, link := range s.links {
		if group == "" || link.Group == group {
			result = append(result, link)
		}
	}
	return result
}

func (s *memoryStore) groups() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	seen := make(map[string]struct{})
	groups := make([]string, 0)
	for _, link := range s.links {
		if _, ok := seen[link.Group]; ok {
			continue
		}
		seen[link.Group] = struct{}{}
		groups = append(groups, link.Group)
	}
	return groups
}
