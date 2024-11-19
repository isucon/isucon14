package world

import (
	"iter"
	"slices"
	"sync"
)

type RequestDB struct {
	counter int
	m       map[RequestID]*Request
	lock    sync.RWMutex
}

func NewRequestDB() *RequestDB {
	return &RequestDB{
		m: make(map[RequestID]*Request),
	}
}

func (db *RequestDB) Create(req *Request) *Request {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.counter++
	req.ID = RequestID(db.counter)
	db.m[req.ID] = req
	return req
}

func (db *RequestDB) Get(id RequestID) *Request {
	db.lock.RLock()
	defer db.lock.RUnlock()
	return db.m[id]
}

func (db *RequestDB) GetByServerID(serverID string) *Request {
	db.lock.RLock()
	defer db.lock.RUnlock()

	// TODO ハッシュマップで持って引くように
	for _, req := range db.m {
		if req.ServerID == serverID {
			return req
		}
	}
	return nil
}

func (db *RequestDB) Iter() iter.Seq2[RequestID, *Request] {
	return func(yield func(RequestID, *Request) bool) {
		db.lock.RLock()
		defer db.lock.RUnlock()
		for id, req := range db.m {
			if !yield(id, req) {
				return
			}
		}
	}
}

func (db *RequestDB) Size() int {
	db.lock.RLock()
	defer db.lock.RUnlock()
	return len(db.m)
}

func (db *RequestDB) Values() iter.Seq[*Request] {
	return func(yield func(*Request) bool) {
		db.lock.RLock()
		defer db.lock.RUnlock()
		for _, v := range db.m {
			if !yield(v) {
				return
			}
		}
	}
}

func (db *RequestDB) ToSlice() []*Request {
	return slices.Collect(db.Values())
}

type DBEntry[K ~int] interface {
	SetID(id K)
	GetServerID() string
}

type GenericDB[K ~int, V DBEntry[K]] struct {
	counter int
	m       map[K]V
	lock    sync.RWMutex
}

func NewGenericDB[K ~int, V DBEntry[K]]() *GenericDB[K, V] {
	return &GenericDB[K, V]{
		m: map[K]V{},
	}
}

func (db *GenericDB[K, V]) Create(v V) V {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.counter++
	v.SetID(K(db.counter))
	db.m[K(db.counter)] = v
	return v
}

func (db *GenericDB[K, V]) Get(id K) V {
	db.lock.RLock()
	defer db.lock.RUnlock()
	return db.m[id]
}

func (db *GenericDB[K, V]) GetByServerID(serverID string) V {
	var zero V
	db.lock.RLock()
	defer db.lock.RUnlock()

	// TODO ハッシュマップで持って引くように
	for _, req := range db.m {
		if req.GetServerID() == serverID {
			return req
		}
	}
	return zero
}

func (db *GenericDB[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		db.lock.RLock()
		defer db.lock.RUnlock()
		for id, v := range db.m {
			if !yield(id, v) {
				return
			}
		}
	}
}

func (db *GenericDB[K, V]) Size() int {
	db.lock.RLock()
	defer db.lock.RUnlock()
	return len(db.m)
}

func (db *GenericDB[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		db.lock.RLock()
		defer db.lock.RUnlock()
		for _, v := range db.m {
			if !yield(v) {
				return
			}
		}
	}
}

func (db *GenericDB[K, V]) ToSlice() []V {
	return slices.Collect(db.Values())
}
