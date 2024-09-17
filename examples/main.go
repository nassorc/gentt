package main

import (
	"fmt"
  "gentt"
)

type PositionData struct { X, Y int }

type VelocityData struct { X, Y int }

var (
  Position = gentt.CreateComponent[PositionData]()
  Velocity = gentt.CreateComponent[VelocityData]()
)

type QueryResult struct {
  Result []int
}

func (res QueryResult) Each() {
}

func main() {
  world := gentt.NewWorld()

  e1 := world.Create(Position, Velocity)
  Position.SetData(world, e1, PositionData{10, 20})
  Velocity.SetData(world, e1, VelocityData{1, 2})

  query := world.Query(Position, Velocity)

  for i := 0; i < len(query); i++ {
    pos := Position.Get(world, query[i]) 
    vel := Velocity.Get(world, query[i]) 

    pos.X += vel.X
    pos.Y += vel.Y
  }

  if entity, ok := Position.First(world); ok {
    fmt.Println("Has first", entity)
  } else {
    fmt.Println("No first", entity)
  }

  fmt.Println(Position.All(world))

  Position.Each(world, func(entity int, data *PositionData) {
    Position.Remove(world, entity)
  })

  fmt.Println(Position.All(world))
}
