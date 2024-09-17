package gentt

import (
	"fmt"
	"reflect"

	"github.com/nassorc/go-codebase/lib/bitset"
	"github.com/nassorc/go-codebase/lib/ringbuffer"
)

const WORLD_SIZE = 10000
const BITSET_SIZE = 8

type EntityId = int

type IComponentType interface {
	RType() reflect.Type
	Zero() interface{}
}

type World struct {
	stores               []*Store
  // Maps a component to its actual index in the stores array 
	componentTypeToStore map[reflect.Type]int
	entityPool           *ringbuffer.Ringbuffer[EntityId]

	// reserve last bit to check if entity is valid
	// 1 = valid, set when creating entity
	// 0 = not valid, reset when destroying entity
	entitySignatures []bitset.Bitset
  systems   []func(*World) 
  renderers   []any
}

func NewWorld() *World {
	entityPool := ringbuffer.NewRingBuffer[EntityId](WORLD_SIZE)
	entitySignatures := make([]bitset.Bitset, WORLD_SIZE)

	for idx := 0; idx < WORLD_SIZE; idx++ {
		entityPool.Enqueue(idx)
		entitySignatures[idx] = bitset.NewBitset(BITSET_SIZE)
	}

	return &World{
		componentTypeToStore: make(map[reflect.Type]int),
		entityPool:           entityPool,
		entitySignatures:     entitySignatures,
	}
}

func (w World) IsValid(entity EntityId) bool {
	return w.entitySignatures[entity].Test(BITSET_SIZE - 1)
}

func (w *World) RegisterSystem(system func(*World)) {
  w.systems = append(w.systems, system)
}

func (w *World) RegisterRenderer(renderer any) {
  w.renderers = append(w.renderers, renderer)
}

func (w *World) Tick() {
  for _, system := range w.systems {
    system(w)
  }
}

func (w *World) Draw(arg any) {
  for _, renderer := range w.renderers {
    val := reflect.ValueOf(renderer)
    val.Call([]reflect.Value{ reflect.ValueOf(w), reflect.ValueOf(arg) })
  }
}

func (w *World) Get(store reflect.Type, entity EntityId) (any, bool) {
	data, ok := w.stores[w.componentTypeToStore[store]].Get(entity)
	return data, ok
}

func (w *World) GetStore(store reflect.Type) *Store {
	return w.stores[w.componentTypeToStore[store]]
}

func (w *World) recycleId(id EntityId) {
	w.entityPool.Enqueue(id)
}

func (w *World) issueId() EntityId {
	if w.entityPool.Empty() {
		panic("Entity capacity reached.")
	}
	id, ok := w.entityPool.Deque()

	if !ok {
		panic("Failed to create entity.")
	}

	return id
}

func (w *World) SetData(entity EntityId, component interface{}) {
	t := reflect.TypeOf(component)

	if !w.HasStore(t) {
		// create store
		idx := len(w.stores)
		w.stores = append(w.stores, NewStore(t))
		w.componentTypeToStore[t] = idx
	}

	storeIdx := w.componentTypeToStore[t]
	w.stores[w.componentTypeToStore[t]].Insert(entity, reflect.ValueOf(component))
	w.entitySignatures[entity].Set(storeIdx)
}

func (w *World) Create(components ...IComponentType) EntityId {
	id := w.issueId()

	for _, component := range components {
		t := component.RType()

    // create store if store of the given type doesn't exist
		if !w.HasStore(t) {
			id := len(w.stores)
			w.stores = append(w.stores, NewStore(component.RType()))
			w.componentTypeToStore[t] = id
		}

    // insert data to store
		storeIdx := w.componentTypeToStore[t]
    store := w.stores[w.componentTypeToStore[t]]
		store.Insert(id, reflect.ValueOf(component.Zero()))
		w.entitySignatures[id].Set(storeIdx)
	}

	// make entity valid
	w.entitySignatures[id].Set(BITSET_SIZE - 1)
	return id
}

func (w *World) Remove(entity EntityId, component reflect.Type) {
	idx, ok := w.componentTypeToStore[component]

	if !ok {
		panic(fmt.Sprintf("Failed to remove component %v from entity %v. Store does not exist", component, entity))
	}

	ok = w.stores[idx].Remove(entity)

	if !ok {
		panic(fmt.Sprintf("Failed to remove component %v from entity %v", component, entity))
	}

	w.entitySignatures[entity].Set(idx)
}

func (w *World) Destroy(entity EntityId) {
	for _, store := range w.stores {
		if store.Has(entity) {
			store.Remove(entity)
		}
	}

	w.recycleId(entity)
	w.entitySignatures[entity].ResetAll()
}

func (w *World) HasStore(t reflect.Type) bool {
	idx, ok := w.componentTypeToStore[t]
	return ok && idx < len(w.stores)
}

func (w *World) Query(components ...IComponentType) []EntityId {
	result := make([]EntityId, 0)
	qSig := bitset.NewBitset(BITSET_SIZE)

	// points to store with the least amount of entities
	var minStore *Store
	// var minSize = WORLD_SIZE + 1

	for _, component := range components {
		t := component.RType()
		idx := w.componentTypeToStore[t]
		store := w.stores[idx]

		if minStore == nil {
			minStore = store
		} else if store.Size() < minStore.Size() {
			minStore = store
		}

		qSig.Set(idx)
	}

	for _, entity := range minStore.ReverseLookup[:minStore.Size()] {
		eSig := w.entitySignatures[entity]

		if (eSig.Int() & qSig.Int()) == qSig.Int() {
			result = append(result, entity)
		}
	}

	return result
}
