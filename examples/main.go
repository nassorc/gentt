package main

import (
	"fmt"
	"gentt"
)

// Components
type PositionData struct { X, Y int }

type VelocityData struct { X, Y int }

type ConfigData struct { msg string }

var (
  Position = gentt.CreateComponent[PositionData]()
  Velocity = gentt.CreateComponent[VelocityData]()
  Config = gentt.CreateComponent[ConfigData]()
)

func DisplayConfig(world *gentt.World) {
  if entity, ok := Config.First(world); ok {
    config := Config.Get(world, entity)
    fmt.Println(config.msg)
  }
}

func main() {
  world := gentt.NewWorld()

  // systems and renderers
  world.RegisterSystem(DisplayConfig)

  e1 := world.Create(Position, Velocity) // creates a new entity and returns its id
  config := world.Create(Config)

  Config.SetData(world, config, ConfigData{ "hello world" })
  Position.SetData(world, e1, PositionData{10, 20})
  Velocity.SetData(world, e1, VelocityData{1, 2})

  world.Tick() // advances the game world and calls each system and renderer

  // Query queries all entities that match function arguments 
  // For example, this code queries all entities that is a subset of {Position, Velocity}
  query := world.Query(Position, Velocity)

  query.Each(func(e gentt.EntityId) {
    pos := Position.Get(world, e) 
    vel := Velocity.Get(world, e) 

    pos.X += vel.X
    pos.Y += vel.Y
  })

  for _, entity := range Position.All(world) {
    world.Destroy(entity)
  }

  moveableQuery := world.Query(Position, Velocity)

  moveableQuery.Each(func(entity gentt.EntityId) {
    world.Destroy(entity)
  })
}
