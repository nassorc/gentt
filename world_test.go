package gentt

import (
	"testing"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type PositionData struct { X, Y int }
type VelocityData struct { X, Y int }
type ConfigData struct { }

var Position = CreateComponent[PositionData]()
var Velocity = CreateComponent[VelocityData]()
var Config = CreateComponent[ConfigData]()

func Test_EcsSystems(t *testing.T) {
  world := NewWorld()
  tickCount := 0
   
  world.RegisterSystem(func(world *World) {
    tickCount += 1
  })

  world.Tick()
  if tickCount != 1 {
    t.Errorf("Expected %v, Got %v", 1, tickCount)
  }

  world.Tick()
  world.Tick()
  if tickCount != 3 {
    t.Errorf("Expected %v, Got %v", 3, tickCount)
  }
}

func Test_EcsRenderers(t *testing.T) {
  world := NewWorld()
  drawCount := 0

  renderer :=  func(world *World, adder int) {
    drawCount += adder
  }
   
  world.RegisterRenderer(renderer)

  world.Draw(2)

  if drawCount != 2 {
    t.Errorf("Expected %v, Got %v", 2, drawCount)
  }

  world.Draw(3)
  world.Draw(5)
  if drawCount != 10 {
    t.Errorf("Expected %v, Got %v", 3, drawCount)
  }
}

func Test_ComponentType(t *testing.T) {
  // The component store does NOT preserve the order of the data.
  // You'll likely use First to query a single global component.
  world := NewWorld()

  world.Create(Position)
  world.Create(Position)
  expected := world.Create(Config)

  if first, ok := Config.First(world); ok {
    if first != expected {
      t.Errorf("Expected %v, Got %v", expected, first)
    }
  } else {
    t.Errorf("Expected %v, Got %v", true, ok)
  }
}

func Test_CreateEntity(t *testing.T) {
  world := NewWorld()

  for idx := 0; idx < 9000; idx++ {
    world.Create(Position)
  }
  for idx := 0; idx < 1000; idx++ {
    world.Create(Position, Velocity)
  }
}

func Benchmark_CreateEntity(b *testing.B) {
  for idx := 0; idx < b.N; idx++ {
    world := NewWorld()

    for idx := 0; idx < 9000; idx++ {
      world.Create(Position)
    }

    for idx := 0; idx < 1000; idx++ {
      world.Create(Position, Velocity)
    }
  }
}

func Benchmark_SystemIteration(b *testing.B) {
  world := NewWorld()

  for idx := 0; idx < 9000; idx++ {
    world.Create(Position)
  }
  for idx := 0; idx < 1000; idx++ {
    world.Create(Position, Velocity)
  }

  query := world.Query(Position, Velocity)

  b.ResetTimer()


  for i := 0; i < b.N; i++ {
    for idx := 0; idx < len(query); idx++ {
      pos := Position.Get(world, query[idx])
      vel := Velocity.Get(world, query[idx])

      pos.X += vel.X
      pos.Y += vel.Y
    }
  }
}

func Benchmark_SystemIterationDonburi(b *testing.B) {
  world := donburi.NewWorld()

  var (
    Pos = donburi.NewComponentType[PositionData]()
    Vel = donburi.NewComponentType[VelocityData]()
  )

  for idx := 0; idx < 9000; idx++ {
    world.Create(Pos)
  }
  for idx := 0; idx < 1000; idx++ {
    world.Create(Pos, Vel)
  }

  // query := world.Query(Position, Velocity)
  query := donburi.NewQuery(filter.Contains(Pos, Vel))

  b.ResetTimer()

  for i := 0; i < b.N; i++ {
    query.Each(world, func(e *donburi.Entry) {
      pos := Pos.Get(e)
      vel := Vel.Get(e)

      pos.X += vel.X
      pos.Y += vel.Y

    })
  }
}

// func Benchmark_Query(b *testing.B) {
//   world := NewWorld()
//
//   for idx := 0; idx < 9000; idx++ {
//     world.Create(Position)
//   }
//   for idx := 0; idx < 5; idx++ {
//     world.Create(Position, Velocity)
//   }
//
//   query := world.Query(Position, Velocity)
//
//   b.ResetTimer()
//
//
//   for idx := 0; idx < b.N; idx++ {
//     for _, entity := range query {
//       pos := Position.Get(world, entity)
//       vel := Velocity.Get(world, entity)
//
//       pos.X += vel.X
//       pos.Y += vel.Y
//     }
//   }
//
//
// }

// func Benchmark_DebugCreating(b *testing.B) {
//   for idx := 0; idx < b.N; idx++ {
//     world := NewWorld()
//     world.Create(Position, Velocity)
//   }
// }
