package access

import (
	"fmt"
	"time"

	"github.com/m1guelpf/chatgpt-telegram/src/expirymap"
)

type Manager struct {
	em expirymap.ExpiryMap
}

func New() *Manager {
	return &Manager{em: expirymap.New()}
}

// Grant даёт пользователю tgID «активен» на days дней
func (m *Manager) Grant(tgID int64, days int) {
	key := fmt.Sprint(tgID)
	m.em.Set(key, "active", time.Duration(days)*24*time.Hour)
}

// Has проверяет, не истёк ли доступ
func (m *Manager) Has(tgID int64) bool {
	key := fmt.Sprint(tgID)
	_, ok := m.em.Get(key)
	return ok
}
