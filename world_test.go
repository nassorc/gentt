package gentt

import "testing"

type PositionData struct { X, Y int }
type VelocityData struct { X, Y int }

var Position = CreateComponent[PositionData]()
var Velocity = CreateComponent[VelocityData]()

func Benchmark_CreateEntity(b *testing.B) {
  world := NewWorld()

  for idx := 0; idx < b.N; idx++ {
    world.Create(Position, Velocity)
  }
}

func Benchmark_SystemIteration(b *testing.B) {
  world := NewWorld()

  for idx := 0; idx < 10000; idx++ {
    world.Create(Position, Velocity)
  }

  // for idx := 0; idx < 1000; idx++ {
  //   world.Create(Position, Velocity)
  // }

  query := world.Query(Position, Velocity)

  b.ResetTimer()

  for entity := 0; entity < len(query); entity++ {
    pos := Position.Get(world, entity)
    vel := Velocity.Get(world, entity)

    pos.X += vel.X
    pos.Y += vel.Y
  }
}

func Benchmark_DebugCreating(b *testing.B) {

  for idx := 0; idx < b.N; idx++ {
    world := NewWorld()
    world.Create(Position, Velocity)
  }
}
