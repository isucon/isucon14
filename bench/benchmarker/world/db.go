package world

import (
	"iter"
	"slices"
	"sync"
)

type RideDB struct {
	counter int
	m       map[RideID]*Ride
	lock    sync.RWMutex
}

func NewRideDB() *RideDB {
	return &RideDB{
		m: make(map[RideID]*Ride),
	}
}

func (db *RideDB) Create(req *Ride) *Ride {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.counter++
	req.ID = RideID(db.counter)
	db.m[req.ID] = req
	return req
}

func (db *RideDB) Get(id RideID) *Ride {
	db.lock.RLock()
	defer db.lock.RUnlock()
	return db.m[id]
}

func (db *RideDB) GetByServerID(serverID string) *Ride {
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

func (db *RideDB) Iter() iter.Seq2[RideID, *Ride] {
	return func(yield func(RideID, *Ride) bool) {
		db.lock.RLock()
		defer db.lock.RUnlock()
		for id, req := range db.m {
			if !yield(id, req) {
				return
			}
		}
	}
}

func (db *RideDB) Size() int {
	db.lock.RLock()
	defer db.lock.RUnlock()
	return len(db.m)
}

func (db *RideDB) Values() iter.Seq[*Ride] {
	return func(yield func(*Ride) bool) {
		db.lock.RLock()
		defer db.lock.RUnlock()
		for _, v := range db.m {
			if !yield(v) {
				return
			}
		}
	}
}

func (db *RideDB) ToSlice() []*Ride {
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
