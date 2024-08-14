package main

import (
	"reflect"
)

type ComponentType[T any] struct {
  rType reflect.Type
}

func CreateComponent[T any]() ComponentType[T] {
  var val T
  t := reflect.TypeOf(val)
  return ComponentType[T]{
    t,
  }
}

func (c ComponentType[T]) Add(world *World, entity EntityId, val T) {
  world.SetData(entity, val)
}

func (c ComponentType[T]) All(world *World) []EntityId {
  store := world.GetStore(c.rType)

  return store.Sparse[:store.Size()]
}

func (c ComponentType[T]) Create(world *World, entity EntityId) {
  var zero T
  world.SetData(entity, &zero)
}

func (c ComponentType[T]) Each(world *World, action func(entity EntityId, data *T))  {
  store := world.GetStore(c.rType)

  // We are purposely reversing the returned array so the user won't invalidate
  // the array when the user removes the entity's component or detroy the entity during iteration.

  // Removing a component does a swap and pop with the last element of the store, so when you
  // are iterating over the data, and in an iteration you remove the entity's component or
  // destroy an entity, the entity that is currently being processed will be swapped with
  // the unprocessed last element in the store.

  // Reverising the data guarantees that swapping only happens with processed elements and
  // it also lessens the burden to the users.

  // E.g. During iteration, if index is 0, we remove the first entity in the list.
  // Otherwise, do nothing.

  //                      Data [d1, d2, d3, d4]
  // current iteration ---------^             ^----size
  //
  // We process d1 and see that index is 0, so we remove d1.
  // swap d1 with d4      Data [d4, d2, d3, d1]
  //                            *            *
  //
  //                      Data [d4, d2, d3, d1]
  // next iteration ----------------^    ^---------size
  // In this example, we never got the chance to process d4

  for idx := store.Size() - 1; idx >= 0; idx-- {
    entity := store.ReverseLookup[idx]
    elm, _ := store.Index(idx)
    data := elm.Addr().Interface().(*T)
    action(entity, data)
  }
}

func (c ComponentType[T]) First(world *World) EntityId {
  // return world.GetStore(c.rType).ReverseLookup[0]
  return world.stores[world.componentTypeToStore[c.rType]].ReverseLookup[0]
}

func (c ComponentType[T]) Get(world *World, entity EntityId) *T {
  store := world.stores[world.componentTypeToStore[c.rType]]
  // data := store.Data.Index(store.ReverseLookup[entity]).Addr().Interface().(*T)
  data, _ := store.Get(entity)

  // store.Data.Index(store.idToDataLookup[entity])
  return data.Addr().Interface().(*T)
  // data := world.stores[world.componentTypeToStore[c.rType]].Index(entity)
  // return data.(*T)
  // var data *T
  // data := &world.stores[world.componentTypeToStore[c.rType]].Dense[entity]
  //
  // return (*data).(T), false

  // val, ok := world.Get(c.rType, entity)
  // // fmt.Println("RECEIVED", *val)
  // return (val).(T), ok
}

func (c ComponentType[T]) Remove(world *World, entity EntityId) {
  world.Remove(entity, Velocity.RType())
}

func (s ComponentType[T]) RType() reflect.Type {
  return s.rType
}

func (s ComponentType[T]) SetData(world *World, entity EntityId, val T)  {
  // data := &world.stores[world.componentTypeToStore[s.rType]].Dense[0]
  // *data = any(val)
  world.SetData(entity, val)

  // data := v.()
  // *data = val
  // world.SetData(entity, val)
}

func (s ComponentType[T]) Zero() interface{} {
  var zero T
  return zero
}

