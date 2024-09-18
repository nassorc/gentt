package gentt

import (
	"reflect"
	"slices"
)

const (
	PAGE_SIZE = 3
)

type Page [PAGE_SIZE]int
type Pages []*Page

type Store struct {
	size          int
	Dense         reflect.Value
	Sparse        Pages
	ReverseLookup []EntityId
}

func NewStore(t reflect.Type) *Store {
	return &Store{
		size:          0,
		Dense:         reflect.MakeSlice(reflect.SliceOf(t), 0, 100),
		ReverseLookup: make([]EntityId, 0, 100),
		Sparse:        make(Pages, 0, 100),
	}
}

func (s *Store) Size() EntityId {
	return s.size
}

func (s *Store) SetSparseIdx(id EntityId, idx int) {
	page := id / PAGE_SIZE
	dataIdx := id % PAGE_SIZE

	if page >= len(s.Sparse) {
		s.Sparse = slices.Grow(s.Sparse, page+1)

    for idx := len(s.Sparse); idx <= page; idx++ {
      s.Sparse = append(s.Sparse, nil)
    }
	}

	if s.Sparse[page] == nil {
    s.Sparse[page] = new(Page)
  }

	s.Sparse[page][dataIdx] = idx
}

func (s Store) GetDataIdx(id EntityId) (int, bool) {
	page := id / PAGE_SIZE
	dataIdx := id % PAGE_SIZE

	if page >= len(s.Sparse) {
		return 0, false
	}

	return s.Sparse[page][dataIdx], true
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
	return s.Dense.Index(idx), true
}

func (s *Store) Get(id EntityId) (reflect.Value, bool) {
	if !s.Has(id) {
		return reflect.Value{}, false
	}

	idx, _ := s.GetDataIdx(id)

	return s.Dense.Index(idx), true
}

func (s *Store) Insert(id EntityId, value reflect.Value) {
	// if (s.Size() + 1) >= s.capacity {
	// 	panic("Full component store.")
	// }

	if s.Has(id) {
		idx, _ := s.GetDataIdx(id)
		s.Dense.Index(idx).Set(value)

	} else {
		idx := s.Size()

		s.SetSparseIdx(id, idx)

		s.Dense = reflect.Append(s.Dense, value)
		s.ReverseLookup = append(s.ReverseLookup, id)

		s.size += 1
	}
}

func (s *Store) Remove(id EntityId) bool {
	if !s.Has(id) || s.Size() == 0 {
		return false
	}

	idx, _ := s.GetDataIdx(id)
	lastIdx := s.Size() - 1
	lastOwnerId := s.ReverseLookup[lastIdx]

	// replace target data with the last element and create new slice excluding the value
	s.Dense.Index(idx).Set(s.Dense.Index(lastIdx))

	// update data
	s.SetSparseIdx(lastOwnerId, idx)
	s.ReverseLookup[idx] = lastOwnerId
	s.size -= 1

	return true
}
