package payment

import "sync"

type orderInfo struct {
	TgID int64
}

type Store struct {
	mu     sync.RWMutex
	orders map[string]orderInfo
}

func NewStore() *Store {
	return &Store{orders: make(map[string]orderInfo)}
}

func (s *Store) Put(orderID string, tgID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[orderID] = orderInfo{TgID: tgID}
}

func (s *Store) Get(orderID string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	info, ok := s.orders[orderID]
	return info.TgID, ok
}

func (s *Store) Delete(orderID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.orders, orderID)
}