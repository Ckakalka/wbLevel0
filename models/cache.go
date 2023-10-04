package models

import (
	"sync"

	"github.com/Ckakalka/wbLevel0/db"
)

type OrderCash struct {
	mx sync.RWMutex
	m  map[string]db.Order
}

func NewOrderCash() *OrderCash {
	return &OrderCash{
		m: make(map[string]db.Order),
	}
}

func (ordCash *OrderCash) Load(key string) (db.Order, bool) {
	ordCash.mx.RLock()
	defer ordCash.mx.RUnlock()
	order, ok := ordCash.m[key]
	return order, ok
}

func (ordCash *OrderCash) Store(key string, order db.Order) {
	ordCash.mx.Lock()
	defer ordCash.mx.Unlock()
	ordCash.m[key] = order
}
