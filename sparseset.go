package main

import (
	"reflect"
	"slices"
)

const (
  PAGE_SIZE = 10
  MAX_PAGES = 1000
)

type Page [PAGE_SIZE]int
type Pages []*Page

type Store struct {
	capacity       int
	size           int
	Data           reflect.Value
	Sparse         []EntityId
	ReverseLookup  []EntityId
  SparsePages    Pages // sparse set
}

func NewStore(t reflect.Type, capacity int) *Store {
	return &Store{
		capacity:       capacity,
		size:           0,
		Data:           reflect.MakeSlice(reflect.SliceOf(t), 0, 10),
		ReverseLookup:  make([]EntityId, 0, 10),

		// Sparse:         make([]EntityId, capacity),

    SparsePages: make(Pages, 0, 0),
	}
}

func (s *Store) Size() EntityId {
	return s.size
}

func (s *Store) SetSparseIdx(id EntityId, idx int) {
  page := id / PAGE_SIZE
  dataIdx := id % PAGE_SIZE

  if page >= len(s.SparsePages) {
    s.SparsePages = slices.Grow(s.SparsePages, page)
    s.SparsePages = append(s.SparsePages, new(Page))
  }

  s.SparsePages[page][dataIdx] = idx
}

func (s Store) GetDataIdx(id EntityId) (int, bool) {
  page := id / PAGE_SIZE
  dataIdx := id % PAGE_SIZE

  if page >= len(s.SparsePages) {
    return 0, false
  }

  return s.SparsePages[page][dataIdx], true
}

func (s *Store) Has(id EntityId) bool {
  idx, ok := s.GetDataIdx(id)

  if !ok {
    return false
  }

  if s.ReverseLookup[idx] == id && idx < s.Size() {
    return true
  }

	return false
}

func (s *Store) Index(idx int) (reflect.Value, bool) {
	return s.Data.Index(idx), true // ! or s.Data.Index(idx).Addr().Elem()
}

func (s *Store) Get(id EntityId) (reflect.Value, bool) {
	if !s.Has(id) {
		return reflect.Value{}, false
	}

  idx, _ := s.GetDataIdx(id)

	return s.Data.Index(idx), true // ! or s.Data.Index(idx).Addr().Elem()
}

func (s *Store) Insert(id EntityId, value reflect.Value) {
	if (s.Size() + 1) >= s.capacity {
		panic("Full component store.")
	}

	if s.Has(id) {
		// idx := s.ReverseLookup[id]
		// s.Data.Index(idx).Set(value)

    idx, _ := s.GetDataIdx(id)
		s.Data.Index(idx).Set(value)

	} else {
    idx := s.Size()

    s.SetSparseIdx(id, idx)

    // s.Data.Index(idx).Set(value)
    s.Data = reflect.Append(s.Data, value)
    s.ReverseLookup = append(s.ReverseLookup, id)
    // s.ReverseLookup[idx] = id
    // s.ReverseLookup[idx] = id
    // s.ReverseLookup[id] = idx
    // s.Sparse[idx] = id

    s.size += 1
	}

}

// This function removes the data of the given id by performing
// a move and pop with the last element.
func (s *Store) Remove(id EntityId) bool {
	if !s.Has(id) || s.Size() == 0 {
		return false
	}

	// idx := s.ReverseLookup[id]
  idx, _ := s.GetDataIdx(id)
	lastIdx := s.Size() - 1
	lastOwnerId := s.ReverseLookup[lastIdx]

	// replace target data with the last element and create new slice excluding the value
	s.Data.Index(idx).Set(s.Data.Index(lastIdx))
	// s.Data = s.Data.Slice(0, lastIdx)  // !

	// update data
	// s.ReverseLookup[lastOwnerId] = idx
	// s.Sparse[idx] = lastOwnerId
  s.SetSparseIdx(lastOwnerId, idx)
  s.ReverseLookup[idx] = lastOwnerId
	s.size -= 1

	return true
}

